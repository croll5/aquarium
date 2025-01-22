/*
Copyright ou © ou Copr. @EroOlf (pseudo GitHub, contacter Cécile Rolland pour toute remarque sur ce code), (21 janvier 2025)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l'analyse des collectes
traçologiques effectuées avec le logiciel DFIR-ORC.

Ce logiciel est régi par la licence CeCILL soumise au droit français et
respectant les principes de diffusion des logiciels libres. Vous pouvez
utiliser, modifier et/ou redistribuer ce programme sous les conditions
de la licence CeCILL telle que diffusée par le CEA, le CNRS et l'INRIA
sur le site "http://www.cecill.info".

En contrepartie de l'accessibilité au code source et des droits de copie,
de modification et de redistribution accordés par cette licence, il n'est
offert aux utilisateurs qu'une garantie limitée.  Pour les mêmes raisons,
seule une responsabilité restreinte pèse sur l'auteur du programme,  le
titulaire des droits patrimoniaux et les concédants successifs.

A cet égard  l'attention de l'utilisateur est attirée sur les risques
associés au chargement,  à l'utilisation,  à la modification et/ou au
développement et à la reproduction du logiciel par l'utilisateur étant
donné sa spécificité de logiciel libre, qui peut le rendre complexe à
manipuler et qui le réserve donc à des développeurs et des professionnels
avertis possédant  des  connaissances  informatiques approfondies.  Les
utilisateurs sont donc invités à charger  et  tester  l'adéquation  du
logiciel à leurs besoins dans des conditions permettant d'assurer la
sécurité de leurs systèmes et ou de leurs données et, plus généralement,
à l'utiliser et l'exploiter dans les mêmes conditions de sécurité.

Le fait que vous puissiez accéder à cet en-tête signifie que vous avez
pris connaissance de la licence CeCILL, et que vous en avez accepté les
termes.
*/

package navigateur

import (
	"aquarium/modules/aquabase"
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

var pourcentageChargement float32 = -1
var colonnesTableNavigateurs []string = []string{"horodatage", "source", "url", "title", "domain_name", "visit_count"}

const (
	req_Firefox = "SELECT url, title, rev_host, datetime(last_visit_date / 1000000, 'unixepoch'), visit_count FROM moz_places;"
	req_Chrome  = "SELECT url, title, 'NONE', datetime(last_visit_time / 1000000, 'unixepoch'), visit_count FROM urls;"
	req_Edge    = "SELECT url, title, 'NONE', datetime(last_visit_time / 1000000, 'unixepoch'), visit_count FROM urls;"
)

type Navigateur struct{}

func (n Navigateur) Extraction(cheminProjet string) error {

	pourcentageChargement = 0

	// Dézipper le dossier Browsers_history.7z
	path := filepath.Join(cheminProjet, "collecteORC", "Browsers", "Browsers_history.7z")
	if err := os.Mkdir(filepath.Join(cheminProjet, "collecteORC", "Browsers", "History"), 0766); err != nil {
		return err
	}
	destPath := filepath.Join(cheminProjet, "collecteORC", "Browsers", "History")
	extractArchive(path, destPath)

	//Init tab logs
	var abase aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	var requeteInsersion aquabase.RequeteInsertion = abase.InitRequeteInsertionExtraction("navigateurs", colonnesTableNavigateurs)

	//List of files

	files, err := ioutil.ReadDir(destPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println(f.Name())
		if f.IsDir() {
			extractFiles, err := ioutil.ReadDir(filepath.Join(destPath, f.Name()))
			if err != nil {
				return err
			}
			for numFichier, extractFile := range extractFiles {
				var pathNavigator = filepath.Join(destPath, f.Name())
				switch f.Name() {
				case "Firefox_Vista_History":
					openDataFiles(filepath.Join(pathNavigator, extractFile.Name()), req_Firefox, &requeteInsersion)
					break
				case "Chrome_Vista_History":
					openDataFiles(filepath.Join(pathNavigator, extractFile.Name()), req_Chrome, &requeteInsersion)
					break
				case "Edge_Anhaeim_History":
					openDataFiles(filepath.Join(pathNavigator, extractFile.Name()), req_Edge, &requeteInsersion)
					break
				default:
					fmt.Println("Navigateur non pris en charge")
				}
				pourcentageChargement = float32(numFichier*100) / float32(len(extractFiles))
			}
		}
	}

	err = requeteInsersion.Executer()
	if err != nil {
		pourcentageChargement = -1
		return err
	}
	pourcentageChargement = 101
	return nil
}

func (n Navigateur) Description() string {
	return "Historique de navigation"
}

func (n Navigateur) PrerequisOK(cheminORC string) bool {
	return true
}

func (n Navigateur) CreationTable(cheminProjet string) error {
	base := aquabase.InitDB_Extraction(cheminProjet)
	base.CreateTableIfNotExist1("navigateurs", colonnesTableNavigateurs, true)
	return nil
}

func (n Navigateur) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if pourcentageChargement == -1 {
		base := aquabase.InitDB_Extraction(cheminProjet)
		if !base.EstTableVide("navigateurs") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (n Navigateur) Annuler() bool {
	// Trop peu de fichiers pour que cela ne soit pertinent
	return pourcentageChargement >= 100
}

func (n Navigateur) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}

func (n Navigateur) SQLChronologie() string {
	return "SELECT id, \"navigateurs\", \"navigateurs\", source, horodatage, \"L'utilisateur a visité la page « \" || title || \" » à l'URL : \" ||  url || \". Nombre total de visites : \" || visit_count FROM navigateurs WHERE NOT horodatage = \"0001-01-01 01:00:00 +0100 CET\""
}

func openDataFiles(filePath string, requete string, requeteInsertion *aquabase.RequeteInsertion) {

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
		var logObf Log
		data.Scan(&logObf.Url, &logObf.Title, &logObf.Domain_name, &logObf.Time_string, &logObf.Visit_count)

		if len(logObf.Domain_name) > 0 && logObf.Domain_name != "NONE" {
			logObf.Reverse_domain()
		}

		logObf.ConvertStringToTime()
		requeteInsertion.AjouterDansRequete(logObf.Time_date.Local(), filePath, logObf.Url, logObf.Title, logObf.Domain_name, logObf.Visit_count)

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
