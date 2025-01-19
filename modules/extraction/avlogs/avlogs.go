package avlogs

import (
	"bytes"
	"fmt"
	"github.com/ulikunitz/xz"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aquarium/modules/aquabase"
)

var pourcentageChargement float32 = -1

var colonnesTableAVLog []string = []string{"timestamp", "eventType", "severity", "description", "source"}

type AvLog struct{}

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

func (a AvLog) Description() string {
	return "Parsage des journaux d'antivirus dans le fichier avlogs"
}
func (a AvLog) Extraction(cheminProjet string) error {

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

	var dejaFait map[string]bool = map[string]bool{}
	var abase aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	var requete aquabase.RequeteInsertion = abase.InitRequeteInsertionExtraction("av_log", colonnesTableAVLog)

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

func (av AvLog) CreationTable(cheminProjet string) error {
	aqua := aquabase.InitDB_Extraction(cheminProjet)
	aqua.CreateTableIfNotExist1("av_log", colonnesTableAVLog, true)
	return nil
}

func (av AvLog) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if pourcentageChargement == -1 {
		bdd := aquabase.InitDB_Extraction(cheminProjet)
		if !bdd.EstTableVide("av_log") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (av AvLog) Annuler() bool {
	return pourcentageChargement >= 100
}

func (a AvLog) PrerequisOK(projectLink string) bool {
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

func (a AvLog) DetailsEvenement(idEvt int) string {
	return "Aucune information supplémentaire n'est disponible"
}

func (a AvLog) SQLChronologie() string {
	return "SELECT id, \"av_log\", \"av_log\", source, timestamp, \"Event: \" || eventType || \" (Severity: \" || severity || \"), Description: \" || description FROM av_log"
}
