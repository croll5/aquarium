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
	var textSpliter []string = strings.Split(utf16LEToUtf8(texte), "\n")
	var typeDerreur string = "NaN"

	var horodatage time.Time
	var nomApp string = "NaN"
	var typeDevent string = "NaN"

	var stringEventTime string = "EventTime"
	var stringAppPath string = "AppPath"
	var stringEventType string = "EventType"

	for _, ligne := range textSpliter {
		valeurs := strings.Split(ligne, "=")
		if valeurs[0] == stringEventTime {
			val := valeurs[1]
			horodatage = FileTimeVersGo(val)
			log.Println(horodatage)

		}
		if valeurs[0] == stringAppPath {
			nomApp = valeurs[1]
			log.Println(nomApp)
		}
		if valeurs[0] == stringEventType {
			typeDevent = valeurs[1]
			log.Println(typeDevent)
		}

	}
	return horodatage, nomApp, typeDevent, typeDerreur
}
func utf16LEToUtf8(s string) string {
	initial := []byte(s)
	runes := []rune{}
	for i := 0; i < len(initial)-1; i += 2 {
		char := uint16(initial[i]) | uint16(initial[i+1])<<8
		runes = append(runes, rune(char))
	}
	return string(runes)
}

func FileTimeVersGo(dateString string) time.Time {
	// Nettoyer la chaîne pour supprimer les caractères de contrôle
	cleanedDate := strings.TrimSpace(dateString)

	date, err := strconv.ParseInt(cleanedDate, 10, 64)
	if err != nil {
		log.Printf("Erreur de conversion de l'horodatage : %s, erreur : %v\n", cleanedDate, err)
		return time.Time{}
	}
	var difference = date / 10000000
	var complement = date % 10000000
	var referentiel = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(referentiel+difference, int64(complement*100))
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
