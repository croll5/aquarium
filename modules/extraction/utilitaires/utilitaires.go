package utilitaires

import (
	"database/sql"
	"encoding/binary"
	"path/filepath"
	"time"
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

func AjoutLogsNavigateur(cheminProjet string, horodatage time.Time, url string, title string, domain_name string, visit_count int) error {
    bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
    if err != nil {
    	return err
    }
    defer bd.Close()
    requete, err := bd.Prepare ("INSERT INTO navigateurs(horodatage, url, title, domain_name, visit_count) VALUES (?, ?, ?, ?, ?)")
    if err != nil {
    		return err
    	}
    	_, err = requete.Exec(horodatage, url, title, domain_name, visit_count)
    	return err
}

func FileTimeVersGo(date []byte) time.Time {
	var dateInt = int64(binary.LittleEndian.Uint64(date))
	var difference = dateInt / 10000000
	var complement = dateInt % 10000000
	var referentiel = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(referentiel+difference, complement)
}
