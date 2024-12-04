package navigateur

import (
    "database/sql"
    "io"
    "aquarium/modules/extraction/utilitaires"
    "log"
    "fmt"
    "github.com/bodgit/sevenzip"
    "os"
    "path/filepath"
    _ "modernc.org/sqlite"
)


const (
    req_Firefox = "SELECT url, title, rev_host, datetime(last_visit_date / 1000000, 'unixepoch'), visit_count FROM moz_places;"
    req_Chrome = "SELECT url, title, 'NONE', datetime(last_visit_time / 1000000, 'unixepoch'), visit_count FROM urls;"
    req_Edge = "SELECT url, title, 'NONE', datetime(last_visit_time / 1000000, 'unixepoch'), visit_count FROM urls;"
)


type Navigateur struct {
	extrait bool
}

func (n Navigateur) Extraction(chemin_projet string) error {
    
    // Dézipper le dossier Browsers_history.7z
    path := chemin_projet + "\\collecteORC\\Browsers\\Browsers_history.7z" 
    if err := os.Mkdir(filepath.Join(chemin_projet, "\\collecteORC\\Browsers\\History"), os.ModeDir); err != nil {
            log.Fatal(err)
    }    
    dest_path := chemin_projet + "\\collecteORC\\Browsers\\History"
    extractArchive(path, dest_path)
    
    //Récupérer les logs
    var logs []Log
    logs = make([]Log, 0)

    openDataFiles(dest_path + "\\Firefox_Vista_History\\B030D72430D6F078_190000001A7F07_E0000001A9CF0_4_places.sqlite_{00000000-0000-0000-0000-000000000000}.data", req_Firefox, &logs)
    fmt.Printf("Taille du tableau de log : %d\n", len(logs))
    openDataFiles(dest_path + "\\Chrome_Vista_History\\B030D72430D6F078_4200000017F074_1000000285C28_3_History_{00000000-0000-0000-0000-000000000000}.data", req_Chrome, &logs)
    fmt.Printf("Taille du tableau de log : %d\n", len(logs))
    openDataFiles(dest_path + "\\Edge_Anhaeim_History\\B030D72430D6F078_5000000029BBA_A000000029BCC_3_History_{00000000-0000-0000-0000-000000000000}.data", req_Edge, &logs)
    fmt.Printf("Taille du tableau de log : %d\n", len(logs))
    
    fmt.Printf("Log : %s",logs[10].Time_string)
    
    for _, log := range logs {
        utilitaires.AjoutLogsNavigateur(chemin_projet, log.Time_date, log.Url, log.Title, log.Domain_name, log.Visit_count)
    }
    
	return nil
}

func (n Navigateur) Description() string {
	return "Historique de navigation [NULL]"
}

func (n Navigateur) PrerequisOK(cheminORC string) bool {
	return true
}

func openDataFiles(filePath string, requete string, logs *[]Log){    

    db, err := sql.Open("sqlite", filePath)
    if err != nil {
        fmt.Printf("Failed to open database")
        return
    }
    defer db.Close()

    data, err := db.Query(requete)
    if err != nil{
        fmt.Printf("Failed to retrieve data")
        return
    }
    
    for data.Next(){
        var log Log
        data.Scan(&log.Url, &log.Title, &log.Domain_name, &log.Time_string, &log.Visit_count)
        log.ConvertStringToTime();
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