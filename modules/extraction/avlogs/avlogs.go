/*
Copyright ou © ou Copr. Valentyna Pronina, (21 janvier 2025)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à [rappeler les
caractéristiques techniques de votre logiciel].

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

package avlogs

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ulikunitz/xz"

	"aquarium/modules/aquabase"
)

var pourcentageChargement float32 = -1

var colonnesTableAVLog []string = []string{"timestamp", "eventType", "severity", "description", "source"}

type AvLogs struct{}

func traiterInfosLog(line string, dejaFait *map[string]bool, source string, requete *aquabase.RequeteInsertion) {
	fields := strings.Split(line, ",") // Assuming CSV format
	if len(fields) < 4 {
		log.Println("Invalid log entry:", line)
		return
	}

	timestamp, eventType, severity, description := fields[0], fields[1], fields[2], fields[3]
	if _, exists := (*dejaFait)[line]; exists {
		return // Skip duplicate entries
	}

	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Println("Invalid timestamp format:", timestamp)
		return
	}

	requete.AjouterDansRequete(parsedTime.String(), eventType, severity, description, source)
	(*dejaFait)[line] = true
}

func (a AvLogs) Description() string {
	return "Parsage des journaux d'antivirus dans le fichier avlogs"
}

func (a AvLogs) Extraction(cheminProjet string) error {

	// Parcourir les fichiers du répertoire
	pourcentageChargement = 0
	logFilePath := filepath.Join(cheminProjet, "collecteAV", "logs", "av_log.xz")
	r, err := os.Open(logFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	xzReader, err := xz.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to decompress log file: %w", err)
	}

	var dejaFait = map[string]bool{}
	var abase = *aquabase.InitDB_Extraction(cheminProjet)
	var requete = abase.InitRequeteInsertionExtraction("av_log", colonnesTableAVLog)

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, xzReader); err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	lines := strings.Split(buffer.String(), "\n")
	for idx, line := range lines {
		traiterInfosLog(line, &dejaFait, logFilePath, &requete)
		pourcentageChargement = float32(idx) * 100 / float32(len(lines))
	}

	requete.Executer()
	pourcentageChargement = 101
	return nil
}

func (av AvLogs) CreationTable(cheminProjet string) error {
	aqua := aquabase.InitDB_Extraction(cheminProjet)
	aqua.CreateTableIfNotExist1("av_log", colonnesTableAVLog, true)
	return nil
}

func (av AvLogs) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if pourcentageChargement == -1 {
		bdd := aquabase.InitDB_Extraction(cheminProjet)
		if !bdd.EstTableVide("av_log") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (av AvLogs) Annuler() bool {
	return pourcentageChargement >= 100
}

func (a AvLogs) PrerequisOK(projectLink string) bool {
	logDir := filepath.Join(projectLink, "logs")
	files, err := os.ReadDir(logDir)
	if err != nil {
		return false
	}
	for _, file := range files {
		if file.Name() == "av_log.xz" {
			return true
		}
	}
	return false
}

func (a AvLogs) DetailsEvenement(idEvt int) string {
	return "Aucune information supplémentaire n'est disponible"
}

func (a AvLogs) SQLChronologie() string {
	return "SELECT id, \"av_log\", \"av_log\", source, timestamp, \"Event: \" || eventType || \" (Severity: \" || severity || \"), Description: \" || description FROM av_log"
}
