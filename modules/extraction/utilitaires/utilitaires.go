package utilitaires

import (
	"database/sql"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bodgit/sevenzip"
)

/*
* Fonction permettant l'insersion d'un évènement dans la table “chronologie“ de la base de données
@cheminProjet : la racine du projet aqua
@extracteur : l'identifiant de l'extracteur
@horodatage : la date à laquelle l'évènement a eu lieu
@source : le fichier duquel a été extrait l'évènement
@message : la destription de l'évènement
@return : une erreur s'il y en a eu une
*
*/
func AjoutEvenementDansBDD(cheminProjet string, extracteur string, horodatage time.Time, source string, message string) error {
	bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
	//log.Println(filepath.Join(cheminProjet, "analyse", "extractions.db"))
	if err != nil {
		return err
	}
	defer bd.Close()
	requete, err := bd.Prepare("INSERT INTO chronologie(extracteur, horodatage, source, message) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = requete.Exec(extracteur, horodatage, source, message)
	return err
}

func AjoutDiversEvenementDansBDD(cheminProjet string, horodatage time.Time, source string, typeOperation string, startSessionTime time.Time, endSessionTime time.Time, exitStatut string, results string) error {
	bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
	//log.Println(filepath.Join(cheminProjet, "analyse", "extractions.db"))
	if err != nil {
		return err
	}
	defer bd.Close()
	requete, err := bd.Prepare("CREATE TABLE IF NOT EXISTS divers(horodatage DATETIME NOT NULL, source VARCHAR(25),typeOperation TEXT, startSessionTime DATETIME NOT NULL, endSessionTime DATETIME NOT NULL, exitStatut TEXT, results TEXT[]); INSERT INTO divers(horodatage, source , typeOperation , startSessionTime , endSessionTime , exitStatut , results) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = requete.Exec(horodatage, source, typeOperation, startSessionTime, endSessionTime, exitStatut, results)
	return err
}

func FileTimeVersGo(date []byte) time.Time {
	var dateInt = int64(binary.LittleEndian.Uint64(date))
	var difference = dateInt / 10000000
	var complement = dateInt % 10000000
	var referentiel = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(referentiel+difference, complement)
}

/* Fonction permettant d'extraire un fichier d'un dossier compressé en 7z
 */
func ExtraireFichierDepuis7z(file *sevenzip.File, destination string) error {
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
