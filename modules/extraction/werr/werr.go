package werr

import (
	"aquarium/modules/aquabase"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bodgit/sevenzip" // Bibliothèque pour gérer les archives .7z
)

var colonnesTableWer []string = []string{"horodatage", "source", "nomApp", "typeDevent", "typeDerreur"}

type Werr struct{}

func (w Werr) Extraction(cheminProjet string) error {
	// Chemin du dossier WER
	cheminWER := filepath.Join(cheminProjet, "collecteORC", "General", "Errors.7z")
	r, err := sevenzip.OpenReader(cheminWER)
	if err != nil {
		return err
	}
	defer r.Close()

	// Parcourt des fichier Errors.7z

	for _, fichierWER := range r.File {
		rc, err := fichierWER.Open()
		if err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		defer rc.Close()
		// Copie du contenu du fichier dans un tampon, pour pouvoir l'ouvrir avec l'extracteur de registres
		var tampon bytes.Buffer
		if _, err := io.Copy(&tampon, rc); err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		tampon.Bytes()
		//Analyse les fichiers
		horodatage, nomApp, typeDevent, typeDerreur := analyserWER(tampon.Bytes())

		//implémentation dans la base de donnée
		var requeteInsertion = aquabase.InitRequeteInsertionExtraction("wer", colonnesTableWer)
		err = requeteInsertion.AjouterDansRequete(horodatage, fichierWER.Name, nomApp, typeDevent, typeDerreur)

		if err != nil {
			return fmt.Errorf("erreur lors l'ajout des valeurs WER dans la base de donnée - phase 1: %v", err)
		}
		err = requeteInsertion.Executer(cheminProjet)

		if err != nil {
			return fmt.Errorf("erreur lors l'ajout des valeurs WER dans la base de donnée - phase 2: %v", err)
		}

	}

	return nil
}

// Fonction d'analyse d'un fichier WER
func analyserWER(contenu []byte) (time.Time, string, string, string) {
	// Convertir le contenu en chaîne de caractères
	texte := string(contenu)

	// Extraire l'horodatage (recherche d'une ligne contenant "ReportTime")
	var textSpliter []string = strings.Split(texte, "\n")
	var typeDerreur string = "NaN"

	var horodatage time.Time
	var nomApp string = "NaN"
	var typeDevent string = "NaN"

	for _, ligne := range textSpliter {
		if strings.Contains(ligne, "EventTime") {
			val := strings.TrimPrefix(ligne, "EventTime=")
			horodatage = FileTimeVersGo(val)
			log.Println("coucou")

		}
		if strings.Contains(ligne, "AppPath") {
			nomApp = strings.TrimPrefix(ligne, "AppPath=")
			log.Println("coucou")
		}
		if strings.HasPrefix(ligne, "EventType=") {
			typeDevent = strings.TrimPrefix(ligne, "EventType=")
			log.Println("coucou")
		}

	}
	return horodatage, nomApp, typeDevent, typeDerreur
}

func FileTimeVersGo(dateString string) time.Time {
	date, err := strconv.Atoi(dateString)
	if err != nil {
		return time.Time{}
	}
	var difference = int64(date) / 10000000
	var complement = int64(date) % 10000000
	var referentiel = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(referentiel+difference, complement)
}

func (w Werr) PrerequisOK(cheminORC string) bool {
	// Vérifie si le dossier Errors/WER existe dans l'analyse
	dossierErrors := filepath.Join(cheminORC, "General", "Errors.7z")
	_, err := os.Stat(dossierErrors)
	return !os.IsNotExist(err)
}

func (w Werr) Description() string {
	return "Extraction des fichiers d'erreur Windows (WER)"
}

func (w Werr) CreationTable(cheminProjet string) error {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	var err error = base.CreateTableIfNotExist1("wer", colonnesTableWer, true)
	return err
}

func (w Werr) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	return -1
}

func (w Werr) Annuler() bool {
	return true
}

func (w Werr) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}
