/*
Copyright ou © ou Copr. Cynthia Calimouttoupoulle, (21 janvier 2025)

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

var pourcentageChargement float32 = -1

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

	for numFichier, fichierWER := range r.File {
		if filepath.Ext(fichierWER.Name) != ".data" {
			continue
		}
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
		adb := aquabase.InitDB_Extraction(cheminProjet)
		tableName := "wer"
		err = adb.CreateTableIfNotExist1(tableName, colonnesTableWer, true)
		var requeteInsertion = adb.InitRequeteInsertionExtraction(tableName, colonnesTableWer)
		err = requeteInsertion.AjouterDansRequete(horodatage, fichierWER.Name, nomApp, typeDevent, typeDerreur)

		if err != nil {
			return fmt.Errorf("erreur lors l'ajout des valeurs WER dans la base de donnée - phase 1: %v", err)
		}
		err = requeteInsertion.Executer()

		if err != nil {
			return fmt.Errorf("erreur lors l'ajout des valeurs WER dans la base de donnée - phase 2: %v", err)
		}
		pourcentageChargement = float32(numFichier) * 100 / float32(len(r.File))
	}
	pourcentageChargement = 101
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

		}
		if valeurs[0] == stringAppPath {
			nomApp = valeurs[1]
		}
		if valeurs[0] == stringEventType {
			typeDevent = valeurs[1]
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
	if verifierTableVide && pourcentageChargement == -1 {
		var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
		if !abase.EstTableVide("wer") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (w Werr) Annuler() bool {
	return pourcentageChargement >= 100
}

func (w Werr) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}

func (w Werr) SQLChronologie() string {
	return "SELECT id, \"wer\", \"wer\", source, horodatage, \"L'application \" || nomApp || \" a généré l'évènement \" || typeDevent FROM wer"
}
