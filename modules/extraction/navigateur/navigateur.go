package navigateur

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/extraction/utilitaires"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
	_ "modernc.org/sqlite"
)

const (
	req_Firefox = "SELECT url, title, rev_host, datetime(last_visit_date / 1000000, 'unixepoch'), visit_count FROM moz_places;"
	req_Chrome  = "SELECT url, title, 'NONE', datetime(last_visit_time / 1000000, 'unixepoch'), visit_count FROM urls;"
	req_Edge    = "SELECT url, title, 'NONE', datetime(last_visit_time / 1000000, 'unixepoch'), visit_count FROM urls;"
)

type Navigateur struct{}

func (n Navigateur) Extraction(chemin_projet string) error {

	// Dézipper le dossier Browsers_history.7z
	path := filepath.Join(chemin_projet, "collecteORC", "Browsers", "Browsers_history.7z")
	if err := os.Mkdir(filepath.Join(chemin_projet, "collecteORC", "Browsers", "History"), os.ModeDir); err != nil {
		log.Fatal(err)
	}
	destPath := filepath.Join(chemin_projet, "collecteORC", "Browsers", "History")
	extractArchive(path, destPath)

	//Init tab logs
	var logs []Log
	logs = make([]Log, 0)

	//List of files

	files, err := ioutil.ReadDir(destPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		if f.IsDir() {
			extractFiles, err := ioutil.ReadDir(filepath.Join(destPath, f.Name()))
			if err != nil {
				log.Fatal(err)
			}
			for _, extractFile := range extractFiles {
				var pathNavigator = filepath.Join(destPath, f.Name())
				switch f.Name() {
				case "Firefox_Vista_History":
					openDataFiles(filepath.Join(pathNavigator, extractFile.Name()), req_Firefox, &logs)
					break
				case "Chrome_Vista_History":
					openDataFiles(filepath.Join(pathNavigator, extractFile.Name()), req_Chrome, &logs)
					break
				case "Edge_Anhaeim_History":
					openDataFiles(filepath.Join(pathNavigator, extractFile.Name()), req_Edge, &logs)
					break
				default:
					fmt.Println("Navigateur non pris en charge")
				}
			}
		}
	}

	for _, log := range logs {
		utilitaires.AjoutLogsNavigateur(chemin_projet, log.Time_date, log.Url, log.Title, log.Domain_name, log.Visit_count)
	}

	return nil
}

func (n Navigateur) Description() string {
	return "Historique de navigation"
}

func (n Navigateur) PrerequisOK(cheminORC string) bool {
	return true
}

func (n Navigateur) CreationTable(cheminProjet string) error {
	var base aquabase.Aquabase = aquabase.InitBDDExtraction(cheminProjet)
	base.CreateTableIfNotExist("navigateurs", []string{"horodatage", "url", "title", "domain_name", "visit_count"})
	return nil
}

func (n Navigateur) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	return -1
}

func (n Navigateur) Annuler() bool {
	return true
}

func (n Navigateur) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}

func openDataFiles(filePath string, requete string, logs *[]Log) {

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		fmt.Printf("Failed to open database\n")
		return
	}
	defer db.Close()

	data, err := db.Query(requete)
	if err != nil {
		fmt.Printf("Failed to retrieve data\n")
		return
	}

	for data.Next() {
		var log Log
		data.Scan(&log.Url, &log.Title, &log.Domain_name, &log.Time_string, &log.Visit_count)

		if len(log.Domain_name) > 0 && log.Domain_name != "NONE" {
			log.Reverse_domain()
		}

		log.ConvertStringToTime()
		*logs = append(*logs, log)

	}
}

/* Fonction permettant d'extraire un fichier d'un dossier compressé en 7z
 */
func extractFile(file *sevenzip.File, destination string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	log.Println("INFO | Dezippage du fichier ", file.Name)
	os.MkdirAll(filepath.Join(destination, filepath.Dir(file.Name)), 0755)
	fichierExtrait, err := os.Create(filepath.Join(destination, file.Name))
	if err != nil {
		log.Println("ERROR | Problème dans la création du fichier de copie : ", err.Error())
	}
	defer fichierExtrait.Close()

	_, err = io.Copy(fichierExtrait, rc)
	if err != nil {
		log.Println("ERROR | Problème dans l'extraction de l'ORC : ", err.Error())
	}

	return nil
}

/* Fonction permettant de décompresser un dossier compressé en 7z
 * Cette fonction utilise la bibliothèque sevenzip,
 * dont la documentation est présente ici : https://pkg.go.dev/github.com/bodgit/sevenzip
 */
func extractArchive(archive string, destination string) error {
	r, err := sevenzip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err = extractFile(f, destination); err != nil {
			return err
		}
	}

	return nil
}
