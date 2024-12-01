package getthis //getthis main // TODO: use package getthis for debugging in this folder
// compiler du go         : go build getthis.go
// execution du programme : ./getthis.exe
import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/bodgit/sevenzip"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"io"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strings"
)

type Getthis struct {
}

/* ******************************************************************** */
/* ********************** Getthis Methods ***************************** */
/* ******************************************************************** */

func (gt Getthis) Extraction(cheminProjet string) error {
	log.Println("Bonjour, je suis censé faire des extractions {Getthis}")
	log.Println("dbPath:" + filepath.Join(cheminProjet, "analyse", "extractions.db"))
	fileToSearch := "GetThis.csv"
	zipExtension := ".7z"
	zipPassword := "avproof"
	// For each unzipped CSV files found
	list_GetThis, err := searchFilesInFolder(fileToSearch, cheminProjet)
	if err != nil {
		return err
	}
	for _, getThis := range list_GetThis {
		// Extract data
		df, err := readCsv(getThis)
		if err != nil {
			return err
		}
		// Export data
		err = exportDfToDb(df, cheminProjet, getThis, "getthis")
		if err != nil {
			return err
		}
	}
	// For each zipped .7z files found
	list_7zFile, err := searchFilesInFolder(zipExtension, cheminProjet)
	if err != nil {
		return err
	}
	for _, archivePath := range list_7zFile {
		// Search CSV files in zip file
		list7z_GetThis, err := searchFilesIn7z(fileToSearch, archivePath, zipPassword)
		if err != nil {
			return err
		}
		for _, getThis7z := range list7z_GetThis {
			zipPath, csvName := splitEndPath(getThis7z, "::")
			// Extract data for all CSV found
			df, err := readCsvIn7zFile(zipPath, csvName, zipPassword)
			if err != nil {
				return err
			}
			// Add a filePath column
			err = exportDfToDb(df, cheminProjet, getThis7z, "getthis")
			if err != nil {
				return err
			}
			//return nil // TODO: Used for save just one file
		}
		// Search zip file in zip file
		z, err := searchFilesIn7z(zipExtension, archivePath, zipPassword)
		if err != nil {
			return err
		}
		if len(z) > 0 {
			fmt.Println("WARNING: files under two archive layers : ", z)
		}
	}
	// No problem in the function
	return nil
}

func (gt Getthis) PrerequisOK(cheminORC string) bool {
	return true
}

func (gt Getthis) Description() string {
	return "Fichier Getthis"
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
	//DB infos
	dbPath := filepath.Join(cheminProjet, "analyse", "extractions.db")

	// Add a filePath column to save the GetThis filename
	colnameFileName := "filePath"
	colvalueList := strings.Split(strings.Split(filname, "::")[0], "\\")
	colvalue := filepath.Join(colvalueList[len(colvalueList)-2], colvalueList[len(colvalueList)-1])
	colvalue = strings.Replace(colvalue, "\\", "/", -1)
	df = dfNewColumn(df, colnameFileName, colvalue)
	fmt.Println(colvalue)

	// Columns filter/selection if columns exist
	/*columnSelection := []string{colnameFileName, "FullName", "MD5", "CreationDate", "LastModificationDate", "LastAccessDate"}
	columns := listItemsInList(columnSelection, df.Names())
	if len(columnSelection) != len(columns) {
		return fmt.Errorf("ERROR: exportDfToDb(E01): [Wrong columns size]")
	}
	// Select the specified columns
	df = df.Select(columns)*/
	columns := df.Names()         // For no columns filter
	columnSelection := df.Names() // For no columns filter

	// Open or create the sqliteDB
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(EO2) [Can't check the DB existance]: %w", err)
	}
	defer db.Close()

	// Check the table exist or create it //TODO A deplacer au début de gt.Extraction() car utile qu'une seule fois
	err = createTableIfNotExists(db, tableName, columnSelection)
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(_) [Can't check the Table existance]: %w", err)
	}

	// Remove all preview values for this GetThis file
	queryDelete := fmt.Sprintf(`DELETE FROM '%s' WHERE filePath = '%s'`, tableName, colvalue)
	_, err = db.Exec(queryDelete)
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(EO3) [Can't delete old values]: %w", err)
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(EO4) [Can't start a transaction]: %w", err)
	}

	// Prepare the insert query
	queryAdd := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ","),
		strings.Repeat("?,", len(columns)-1)+"?")
	stmt, err := tx.Prepare(queryAdd)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR: exportDfToDb(EO5) [Can't Prepare query]: %w", err)
	}
	defer stmt.Close()

	// Ajout des lignes dans la table
	for i := 0; i < df.Nrow(); i++ {
		//row := df.Subset([]int{i})
		values := make([]interface{}, df.Ncol())
		for j := 0; j < df.Ncol(); j++ {
			values[j] = df.Elem(0, j).String()
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ERROR: exportDfToDb(EO6) [Can't add data]: %w", err)
		}
	}

	// Commit la transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(EO7) [Can't Commit transaction]: %w", err)
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

func dfHead(df dataframe.DataFrame, nFirstRows int) dataframe.DataFrame {
	indices := make([]int, nFirstRows)
	for i := 0; i < nFirstRows; i++ {
		indices[i] = i
	}
	return df.Subset(indices)
} // TODO: useful for debugging

func dfNewColumn(df dataframe.DataFrame, colname string, value string) dataframe.DataFrame {
	sourceColumn := make([]string, df.Nrow())
	for i := range sourceColumn {
		sourceColumn[i] = value
	}
	return df.Mutate(series.New(sourceColumn, series.String, colname))
}

func placeholders(n int) string {
	p := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			p += ","
		}
		p += "?"
	}
	return p
} // TODO: TEMP, is for db saving

func checkIfRowExists(db *sql.DB, tableName string, row dataframe.DataFrame, uniqueCols []string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE ", tableName)
	conditions := ""
	values := []interface{}{}
	for _, col := range uniqueCols {
		if conditions != "" {
			conditions += " AND "
		}
		conditions += fmt.Sprintf("%s = ?", col)
		values = append(values, row.Col(col).Elem(0).String())
	}
	query += conditions

	var exists int
	err := db.QueryRow(query, values...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists == 1, nil
} // TODO: TEMP, is for db saving

func createTableIfNotExists(db *sql.DB, tableName string, columns []string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)
	for i, col := range columns {
		query += fmt.Sprintf("%s TEXT", col)
		if i < len(columns)-1 {
			query += ", "
		}
	}
	query += ");"

	_, err := db.Exec(query)
	return err
} // TODO: TEMP, is for db saving

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

/* ******************************************************************** */
/* **************************** Main  ********************************* */
/* ******************************************************************** */

// Fonction exportée, accessible de l'extérieur du package
func ExportedFunction() {
	fmt.Println("Ceci est une fonction exportée")
} // TODO: memo de la portabilitee des fonctions

// Fonction non exportée, inaccessible de l'extérieur du package
func unexportedFunction() {
	fmt.Println("Ceci est une fonction non exportée")
} // TODO: memo de la portabilitee des fonctions

func main() {
	directory := "C:\\Users\\charm\\Downloads\\Nouveau dossier"
	getthis := Getthis{}
	err := getthis.Extraction(directory)
	if err != nil {
		fmt.Println(err)
		return
	}
} // TODO: main() for debugging
