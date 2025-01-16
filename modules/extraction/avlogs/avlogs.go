package avlogs

import (
	"aquarium/modules/extraction/utilitaires"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var pourcentageChargement float32 = -1

var colonnesTableAVLog []string = []string{"timestamp", "eventType", "severity", "description", "source"}

type AvLog struct{}

func (a AvLog) Parse(projectLink string) {

}

func (a AvLog) Description() string {
	return "Parsage des journaux d'antivirus dans le fichier avlogs"
}
func (a AvLog) Extraction(logDir string) error {

	// Parcourir les fichiers du répertoire
	err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erreur lors de l'accès au chemin %s: %v", path, err)
		}

		// Vérifier si c'est un fichier
		if !info.IsDir() {
			fmt.Printf("Traitement du fichier : %s\n", path)

			// Lire le fichier
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("impossible d'ouvrir le fichier %s: %v", path, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				// Analyser la ligne pour détecter une erreur
				timestamp, errorMsg := parseLogLine(line)
				if timestamp != "" && errorMsg != "" {
					// Préparer l'événement à envoyer
					event := map[string]interface{}{
						"timestamp": timestamp,
						"error":     errorMsg,
						"filename":  filepath.Base(path),
					}

					// Enregistrer dans la base de données
					err := utilitaires.AjoutEvenementDansBDD(event)
					if err != nil {
						return fmt.Errorf("échec de l'ajout à la base de données: %v", err)
					}
				}
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("erreur lors de la lecture du fichier %s: %v", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("échec de l'extraction: %v", err)
	}

	fmt.Println("Extraction terminée avec succès.")
	return nil
}

func parseLogLine(line string) (timestamp string, errorMsg string) {
	// Exemple de regex pour extraire l'horodatage et le message d'erreur
	regex := regexp.MustCompile(`\[(.*?)\]\sERROR:\s(.*)`)
	matches := regex.FindStringSubmatch(line)

	if len(matches) == 3 {
		timestamp = matches[1]
		errorMsg = matches[2]
		return timestamp, errorMsg
	}

	return "", ""
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
