package divers

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"aquarium/modules/extraction/utilitaires"

	"github.com/bodgit/sevenzip"
)

type Divers struct{}

// Extraction analyse le contenu de l'archive TextLogs.7z et extrait les événements pertinents.
func (d Divers) Extraction(cheminProjet string) error {
	// Définition du chemin de l'archive TextLogs.7z
	textLogsPath := filepath.Join(cheminProjet, "collecteORC", "General", "TextLogs.7z")

	// Ouverture de l'archive TextLogs.7z
	archive, err := sevenzip.OpenReader(textLogsPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de l'archive TextLogs.7z: %v", err)
		return err
	}
	defer archive.Close()

	// Parcours des fichiers dans TextLogs.7z pour trouver le dossier "divers"
	for _, file := range archive.File {
		// Vérifier si le fichier se trouve dans le dossier "divers"
		if !strings.HasSuffix(file.Name, "/") && filepath.Base(filepath.Dir(file.Name)) == "divers" {
			log.Printf("Traitement du fichier : %s", file.Name)

			// Ouverture du fichier pour lecture
			rc, err := file.Open()
			if err != nil {
				log.Printf("Erreur lors de l'ouverture du fichier %s: %v", file.Name, err)
				continue
			}

			// Copie du contenu du fichier dans un tampon pour traitement
			var tampon bytes.Buffer
			if _, err := io.Copy(&tampon, rc); err != nil {
				log.Printf("Erreur lors de la copie du contenu du fichier %s dans le tampon : %v", file.Name, err)
				rc.Close()
				continue
			}
			rc.Close() // Assurez-vous de libérer les ressources après lecture

			// Extraire et enregistrer les événements à partir du contenu
			if err := d.extraireEtEnregistrerEvenements(tampon.Bytes(), file.Name, cheminProjet); err != nil {
				log.Printf("Erreur lors de l'extraction et de l'enregistrement des événements pour le fichier %s: %v", file.Name, err)
			}
		}
	}

	log.Println("Extraction terminée.")
	return nil
}

// Méthode pour extraire et enregistrer les événements dans la base de données
func (d Divers) extraireEtEnregistrerEvenements(contenu []byte, nomFichier string, cheminProjet string) error {
	// Expressions régulières pour détecter les sections et événements
	sectionStartRegex := regexp.MustCompile(`Section start (\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3})`)
	eventRegex := regexp.MustCompile(`sto:\s+\{Unpublish Driver Package:\s+(.+?)}\s+(\d{2}:\d{2}:\d{2}\.\d{3})`)
	bootRegex := regexp.MustCompile(`Boot Session.*(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3})`)

	lines := strings.Split(string(contenu), "\n")
	var horodatageSession time.Time

	// Recherche de l'horodatage global de session
	for _, line := range lines {
		if match := bootRegex.FindStringSubmatch(line); match != nil {
			sessionTimestamp, err := time.Parse("2006/01/02 15:04:05.000", match[1])
			if err == nil {
				horodatageSession = sessionTimestamp
			}
			break
		}
	}

	// Parcours des lignes pour détecter et enregistrer les événements
	for _, line := range lines {
		// Début de section
		if match := sectionStartRegex.FindStringSubmatch(line); match != nil {
			sectionTimestamp, _ := time.Parse("2006/01/02 15:04:05.000", match[1])
			message := "Début de la section de mise à jour de pilotes"
			if err := utilitaires.AjoutEvenementDansBDD(cheminProjet, "divers", sectionTimestamp, nomFichier, message); err != nil {
				log.Printf("Erreur lors de l'ajout de l'événement de début de section: %v", err)
			}
		}

		// Événements spécifiques
		if match := eventRegex.FindStringSubmatch(line); match != nil {
			driverPath := match[1]
			eventTime := match[2]
			eventTimestamp, _ := time.Parse("15:04:05.000", eventTime)
			eventTimestamp = time.Date(horodatageSession.Year(), horodatageSession.Month(), horodatageSession.Day(),
				eventTimestamp.Hour(), eventTimestamp.Minute(), eventTimestamp.Second(), eventTimestamp.Nanosecond(), time.UTC)

			message := "Tentative de désinstallation du pilote: " + driverPath
			if err := utilitaires.AjoutEvenementDansBDD(cheminProjet, "divers", eventTimestamp, nomFichier, message); err != nil {
				log.Printf("Erreur lors de l'ajout de l'événement de pilote : %v", err)
			}
		}
	}

	log.Println("Extraction et enregistrement des événements terminés.")
	return nil
}

// Description retourne une description du module.
func (d Divers) Description() string {
	return "Extraction de divers logs"
}

// PrerequisOK vérifie si les prérequis pour l'extraction sont remplis.
func (d Divers) PrerequisOK(cheminORC string) bool {
	// Vérification de l'existence de "TextLogs.7z"
	textLogsPath := filepath.Join(cheminORC, "General", "TextLogs.7z")
	//log.Printf("Chemin construit pour TextLogs.7z : %s", textLogsPath)
	if _, err := os.Stat(textLogsPath); os.IsNotExist(err) {
		log.Println("blocage 1")
		return false
	}

	// Vérification du contenu de l'archive
	archive, err := sevenzip.OpenReader(textLogsPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de l'archive TextLogs.7z: %v", err)
		return false
	}
	defer archive.Close()

	// Recherche d'un dossier nommé "divers"
	for _, file := range archive.File {
		//log.Println(filepath.Dir(file.Name))
		// test
		if filepath.Dir(file.Name) == "divers" {
			return true
		}
	}

	return false
}

func (d Divers) CreationTable(cheminProjet string) error {
	return nil
}

func (d Divers) PourcentageChargement() int {
	return 0
}
