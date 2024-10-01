package utilitaires

import (
	"database/sql"
	"path/filepath"
	"time"
)

func AjoutEvenementDansBDD(cheminProjet string, extracteur string, horodatage time.Time, source string, message string) error {
	bd, err := sql.Open("sqlite3", filepath.Join(cheminProjet, "analyse", "extractions.db"))
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
