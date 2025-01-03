package getthis

// compiler du go         : go build getthis.go
// execution du programme : ./getthis.exe
import (
	"aquarium/modules/aquabase"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	_ "modernc.org/sqlite"
)

var colnameFileName string = "extracteur"
var columnSelection []string = []string{colnameFileName, "FullName", "MD5", "CreationDate", "LastModificationDate", "LastAccessDate"}

var pourcentageChargement float32 = -1

var annulationDemandee bool = false
var annulationReussie bool = false

type Getthis struct{}

/* ******************************************************************** */
/* ********************** Getthis Methods ***************************** */
/* ******************************************************************** */

func (gt Getthis) PrerequisOK(cheminORC string) bool {
	return true
}

func (gt Getthis) Description() string {
	return "Fichier Getthis"
}

func (gt Getthis) CreationTable(cheminProjet string) error {
	base := aquabase.InitDB_Extraction(cheminProjet)
	base.CreateTableIfNotExist("getthis", columnSelection)
	return nil
}

func (gt Getthis) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if pourcentageChargement == -1 {
		base := aquabase.InitDB_Extraction(cheminProjet)
		if !base.EstTableVide("getthis") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (gt Getthis) Annuler() bool {
	if annulationReussie {
		annulationReussie = false
		return true
	}
	if !annulationDemandee {
		annulationDemandee = true
		annulationReussie = false
	}
	return annulationReussie
}

func (gt Getthis) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}

func (gt Getthis) Extraction(cheminProjet string) error {
	pourcentageChargement = 0
	//log.Println("Bonjour, je suis censé faire des extractions {Getthis}")
	//log.Println("dbPath:" + filepath.Join(cheminProjet, "analyse", "extractions.db"))
	fileToSearch := "GetThis.csv"
	zipExtension := ".7z"
	zipPassword := "avproof"
	// For each unzipped CSV files found
	list_GetThis, err := searchFilesInFolder(fileToSearch, cheminProjet)
	if err != nil {
		return err
	}
	for numFile, getThis := range list_GetThis {
		// Extract data
		if annulationDemandee {
			err := viderTableGetThis(cheminProjet)
			if err == nil {
				annulationDemandee = false
				annulationReussie = true
				return nil
			}
		}
		df, err := readCsv(getThis)
		if err != nil {
			return err
		}
		// Export data
		err = exportDfToDb(df, cheminProjet, getThis, "getthis")
		if err != nil {
			return err
		}
		pourcentageChargement = (float32(numFile) * 10) / float32(len(list_GetThis))
	}
	// For each zipped .7z files found
	list_7zFile, err := searchFilesInFolder(zipExtension, cheminProjet)
	if err != nil {
		return err
	}
	for numFile, archivePath := range list_7zFile {
		// Search CSV files in zip file
		if annulationDemandee {
			err := viderTableGetThis(cheminProjet)
			if err == nil {
				annulationDemandee = false
				annulationReussie = true
				return nil
			}
		}
		list7z_GetThis, err := searchFilesIn7z(fileToSearch, archivePath, zipPassword)
		if err != nil {
			fmt.Println("Error skipped searchFilesIn7z: "+fileToSearch+" --- ", err)
			continue //return err
		}
		for _, getThis7z := range list7z_GetThis {
			zipPath, csvName := splitEndPath(getThis7z, "::")
			// Extract data for all CSV found
			df, err := readCsvIn7zFile(zipPath, csvName, zipPassword)
			if err != nil {
				fmt.Println("Error skipped readCsvIn7zFile: "+getThis7z+" --- ", err)
				continue //return err
			}
			// Add a filePath column
			err = exportDfToDb(df, cheminProjet, getThis7z, "getthis")
			if err != nil {
				fmt.Println("Error skipped exportDfToDb: "+getThis7z+" --- ", err)
				continue //return err
			}
			//return nil // Used for save just one file
		}
		// Search zip file in zip file
		z, err := searchFilesIn7z(zipExtension, archivePath, zipPassword)
		if err != nil {
			return err
		}
		if len(z) > 0 {
			fmt.Println("WARNING: files under two archive layers : ", z)
		}
		pourcentageChargement = 10 + (float32(numFile)*90)/float32(len(list_7zFile))
	}
	// No problem in the function
	pourcentageChargement = 101
	if annulationReussie {
		pourcentageChargement = -1
	}
	return nil
}

/* **************************************************************************** */
/* *********************** GetThis Utils Methods ****************************** */
/* **************************************************************************** */

func searchFilesInFolder(fileName string, folderPath string) ([]string, error) {
	// For all files in folderPath and his subFolders
	var paths []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, fileName) {
			// Save the path if it ends with fileName
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ERROR: searchFilesInFolder() : %w", err)
	}
	return paths, nil
}

func readCsv(filePath string) (dataframe.DataFrame, error) {
	var df dataframe.DataFrame
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return df, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	// Load the CSV file into a dataframe
	df = dataframe.ReadCSV(file)
	return df, nil
}

func searchFilesIn7z(fileName string, zipPath string, password string) ([]string, error) {
	var paths []string
	// Extract and open archivePath
	r, err := sevenzip.OpenReaderWithPassword(zipPath, password)
	if err != nil {
		return paths, fmt.Errorf("ERROR: searchFilesIn7z() : %w", err)
	}
	defer func(r *sevenzip.ReadCloser) {
		_ = r.Close()
	}(r)
	// Search fileName in zipPath
	for _, file := range r.File {
		if strings.HasSuffix(file.Name, fileName) {
			// Save the path if it ends with fileName
			paths = append(paths, zipPath+"::"+file.Name)
		}
	}
	return paths, nil
}

func readCsvIn7zFile(zipPath string, localisationIn7zFile string, password string) (dataframe.DataFrame, error) {
	var df dataframe.DataFrame

	// Extract and open archivePath
	r, err := sevenzip.OpenReaderWithPassword(zipPath, password)
	if err != nil {
		return df, fmt.Errorf("ERROR: readCsvIn7zFile(E01) : %w", err)
	}
	defer func(r *sevenzip.ReadCloser) {
		_ = r.Close()
	}(r)
	// Search fileName in zipPath
	for _, file := range r.File {
		if file.Name == localisationIn7zFile {
			// Open file inside the archive
			rc, err := file.Open()
			if err != nil {
				return df, fmt.Errorf("ERROR: readCsvIn7zFile(E02) : %w", err)
			}
			defer func(rc io.ReadCloser) {
				_ = rc.Close()
			}(rc)
			// Read the file contents into a buffer
			var buf bytes.Buffer
			_, err = io.Copy(&buf, rc)
			if err != nil {
				return df, fmt.Errorf("ERROR: readCsvIn7zFile(E03) : %w", err)
			}
			// Load the CSV data into a dataframe
			df := dataframe.ReadCSV(&buf)
			return df, nil
		}
	}
	return df, fmt.Errorf("ERROR: readCsvIn7zFile(E04) : %w", err)

}

func exportDfToDb(df dataframe.DataFrame, cheminProjet string, filname string, tableName string) error {
	adb := aquabase.InitDB_Extraction(cheminProjet)

	// Add a filePath column to save the GetThis filename

	colvalueList := strings.Split(strings.Split(filname, "::")[0], "\\")
	colvalue := filepath.Join(colvalueList[len(colvalueList)-2], colvalueList[len(colvalueList)-1])
	colvalue = strings.Replace(colvalue, "\\", "/", -1)
	df = DfAddColumn(df, colnameFileName, colvalue)
	fmt.Println("Import GetThis to DB: " + colvalue)

	// Columns filter/selection if columns exist

	columns := listItemsInList(columnSelection, df.Names())
	if len(columnSelection) != len(columns) {
		return fmt.Errorf("ERROR: exportDfToDb(E01): [Wrong columns size]")
	}
	// Select the specified columns
	df = df.Select(columns)
	//columns := df.Names()         // For no columns filter
	//columnSelection := df.Names() // For no columns filter

	// Check the table exist
	err := adb.CreateTableIfNotExist(tableName, df.Names())
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(_) [createTableIfNotExist]: %w", err)
	}
	// Add the filename column
	df = DfAddColumn(df, colnameFileName, colvalue)
	// Remove all preview values for this file
	where := colnameFileName + "='" + colvalue + "'"
	err = adb.RemoveFromWhere(tableName, where)
	if err != nil {
		return fmt.Errorf("SaveDf: %w", err)
	}
	// export data
	err = adb.SaveDf(df, tableName)
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(_) [AjoutEvenementDansBDD]: %w", err)
	}
	return nil
}

/* ******************************************************************** */
/* *********************** Utils Methods ****************************** */
/* ******************************************************************** */

func splitEndPath(fullPath string, pattern string) (string, string) {
	index := strings.LastIndex(fullPath, pattern)
	if index == -1 {
		fmt.Printf("WARNING: : no patern \"%s\" found\n", pattern)
		return fullPath, ""
	}
	part1 := fullPath[:index]
	part2 := fullPath[index+len(pattern):]
	return part1, part2
}

/** Find and return all elements of smallList existing in bigList **/
func listItemsInList(smallList, bigList []string) []string {
	var result []string
	for _, smallItem := range smallList {
		for _, bigItem := range bigList {
			if strings.Contains(bigItem, smallItem) {
				result = append(result, smallItem)
				break
			}
		}
	}
	return result
}

func DfHead(df dataframe.DataFrame, nFirstRows int) dataframe.DataFrame {
	indices := make([]int, nFirstRows)
	for i := 0; i < nFirstRows; i++ {
		indices[i] = i
	}
	return df.Subset(indices)
}

func DfAddColumn(df dataframe.DataFrame, colname string, value string) dataframe.DataFrame {
	sourceColumn := make([]string, df.Nrow())
	for i := range sourceColumn {
		sourceColumn[i] = value
	}
	return df.Mutate(series.New(sourceColumn, series.String, colname))
}

func viderTableGetThis(cheminProjet string) error {
	base := aquabase.InitDB_Extraction(cheminProjet)
	return base.RemoveFromWhere("getthis", "1=1")
}
