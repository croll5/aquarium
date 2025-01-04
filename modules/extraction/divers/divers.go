package divers

import (
	"aquarium/modules/extraction/utilitaires"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bodgit/sevenzip"
)

// Evenement représente les informations extraites d'un événement dans un log.
type Evenement struct {
	Horodatage    time.Time
	TypeOperation string
	StartSession  time.Time
	EndSession    time.Time
	ExitStatus    string
	Results       string
}

type Divers struct{}

// Extraction analyse le contenu de l'archive TextLogs.7z et extrait les événements pertinents.
func (d Divers) Extraction(cheminProjet string) error {
	textLogsPath := filepath.Join(cheminProjet, "collecteORC", "General", "TextLogs.7z")

	archive, err := sevenzip.OpenReader(textLogsPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de l'archive TextLogs.7z: %v", err)
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		if !strings.HasSuffix(file.Name, "/") && filepath.Base(filepath.Dir(file.Name)) == "divers" {
			log.Printf("Traitement du fichier : %s", file.Name)

			rc, err := file.Open()
			if err != nil {
				log.Printf("Erreur lors de l'ouverture du fichier %s: %v", file.Name, err)
				continue
			}

			var tampon bytes.Buffer
			if _, err := io.Copy(&tampon, rc); err != nil {
				log.Printf("Erreur lors de la copie du contenu du fichier %s dans le tampon : %v", file.Name, err)
				rc.Close()
				continue
			}
			rc.Close()

			if err := d.extraireEtReformater(tampon.String(), file.Name, cheminProjet); err != nil {
				log.Printf("Erreur lors de l'extraction et du reformattage pour le fichier %s: %v", file.Name, err)
			}
		}
	}

	log.Println("Extraction terminée.")
	return nil
}

// extraireEtReformater extrait les informations pertinentes d'un fichier log.
func (d Divers) extraireEtReformater(contenu string, nomFichier string, cheminProjet string) error {
	lines := strings.Split(contenu, "\n")
	var horodatageSession time.Time
	var evenements []Evenement

	var evenementCourant Evenement
	var collecterTexte bool
	var texteSection strings.Builder

	// Format attendu des dates
	const dateFormat = "2006/01/02 15:04:05.000"

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Détecte l'horodatage global de session (Boot Session)
		if strings.HasPrefix(trimmed, "[Boot Session") {
			parts := strings.Split(trimmed, " ")
			if len(parts) >= 3 {
				sessionTime, err := time.Parse(dateFormat, parts[len(parts)-2]+" "+strings.TrimSuffix(parts[len(parts)-1], "]"))
				if err != nil {
					log.Printf("Erreur de parsing de l'horodatage global : %v", err)
					continue
				}
				horodatageSession = sessionTime
			}
			continue
		}

		// Détecte le type d'opération
		if strings.HasPrefix(trimmed, ">>>  [") && strings.HasSuffix(trimmed, "]") {
			typeOperation := strings.Trim(trimmed, ">>>  []")
			evenementCourant = Evenement{
				Horodatage:    horodatageSession,
				TypeOperation: typeOperation,
			}
			continue
		}

		// Début de section
		if strings.HasPrefix(trimmed, ">>>  Section start") {
			parts := strings.Split(trimmed, " ")
			if len(parts) >= 4 {
				startTime, err := time.Parse(dateFormat, parts[4]+" "+parts[5])
				if err != nil {
					log.Printf("Erreur de parsing de la StartSession : %v", err)
					continue
				}
				evenementCourant.StartSession = startTime
				collecterTexte = true
				texteSection.Reset()
			}
			continue
		}

		// Fin de section
		if strings.HasPrefix(trimmed, "<<<  Section end") {
			parts := strings.Split(trimmed, " ")
			if len(parts) >= 4 {
				endTime, err := time.Parse(dateFormat, parts[4]+" "+parts[5])
				if err != nil {
					log.Printf("Erreur de parsing de la EndSession : %v", err)
					continue
				}
				evenementCourant.EndSession = endTime
				evenementCourant.Results = texteSection.String()
				collecterTexte = false
			}
			continue
		}

		// Statut de sortie
		if strings.HasPrefix(trimmed, "<<<  [Exit status:") {
			status := strings.TrimSuffix(strings.TrimPrefix(trimmed, "<<<  [Exit status: "), "]")
			evenementCourant.ExitStatus = status
			evenements = append(evenements, evenementCourant)
			continue
		}

		// Collecte le texte entre start et end
		if collecterTexte {
			texteSection.WriteString(trimmed + "\n")
		}
	}

	// Enregistrement des événements dans la base de données
	for _, evt := range evenements {
		if err := utilitaires.AjoutDiversEvenementDansBDD(
			cheminProjet,
			evt.Horodatage,
			nomFichier,
			evt.TypeOperation,
			evt.StartSession,
			evt.EndSession,
			evt.ExitStatus,
			evt.Results,
		); err != nil {
			log.Printf("Erreur lors de l'ajout dans la base de données: %v", err)
		}
	}

	log.Println("Extraction et enregistrement des événements terminés.")
	return nil
}

// Description retourne une description du module.
func (d Divers) Description() string {
	return "Journaux d'évenements divers"
}

// PrerequisOK vérifie si les prérequis pour l'extraction sont remplis.
func (d Divers) PrerequisOK(cheminORC string) bool {
	//log.Println(cheminORC)
	textLogsPath := filepath.Join(cheminORC, "General", "TextLogs.7z")
	if _, err := os.Stat(textLogsPath); os.IsNotExist(err) {
		//log.Println("blocage 1")
		return false
	}

	archive, err := sevenzip.OpenReader(textLogsPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de l'archive TextLogs.7z: %v", err)
		return false
	}
	defer archive.Close()

	for _, file := range archive.File {

		if filepath.Dir(file.Name) == "divers" {
			//log.Println("ca marche")
			return true
		}
	}

	return false
}

func (d Divers) CreationTable(cheminProjet string) error {
	return nil
}

func (d Divers) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	return -1
}

func (d Divers) Annuler() bool {
	return true
}

func (d Divers) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}
