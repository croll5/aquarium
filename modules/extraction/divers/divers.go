/*
Copyright ou © ou Copr. Abderahman Moad, (21 janvier 2025)

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

package divers

import (
	"aquarium/modules/aquabase"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bodgit/sevenzip"
)

var colonnesTableDivers []string = []string{"horodatage", "source", "typeOperation", "startSessionTime", "endSessionTime", "exitStatut", "results"}
var pourcentageChargement float32 = -1

var annulationDemandee bool = false
var annulationReussie bool = false

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

	for numFichier, file := range archive.File {
		if annulationDemandee {
			return annulerExtraction(cheminProjet)
		}
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
		pourcentageChargement = float32(numFichier*100) / float32(len(archive.File))
	}

	log.Println("Extraction divers terminée.")
	pourcentageChargement = 101
	return nil
}

func annulerExtraction(cheminProjet string) error {
	var base aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	err := base.RemoveFromWhere("divers", "1=1")
	if err == nil {
		pourcentageChargement = -1
		annulationReussie = true
	}
	return err
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
				sessionTime, err := time.ParseInLocation(dateFormat, parts[len(parts)-2]+" "+strings.TrimSuffix(parts[len(parts)-1], "]"), time.Local)
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
				startTime, err := time.ParseInLocation(dateFormat, parts[4]+" "+parts[5], time.Local)
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
				endTime, err := time.ParseInLocation(dateFormat, parts[4]+" "+parts[5], time.Local)
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
	var abase aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	var requeteInsertion aquabase.RequeteInsertion = abase.InitRequeteInsertionExtraction("divers", colonnesTableDivers)
	var err error
	for _, evt := range evenements {
		err = requeteInsertion.AjouterDansRequete(evt.Horodatage, nomFichier, evt.TypeOperation, evt.StartSession, evt.EndSession, evt.ExitStatus, evt.Results)
		if err != nil {
			return err
		}
	}
	err = requeteInsertion.Executer()
	if err != nil {
		return err
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
	var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	err := abase.CreateTableIfNotExist1("divers", colonnesTableDivers, true)
	return err
}

func (d Divers) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	var abase aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	if pourcentageChargement == -1 && verifierTableVide {
		if !abase.EstTableVide("divers") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (d Divers) Annuler() bool {
	if !annulationDemandee {
		annulationDemandee = true
		annulationReussie = false
	}
	if annulationReussie {
		annulationDemandee = false
	}
	return annulationReussie
}

func (d Divers) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}

func (d Divers) SQLChronologie() string {
	return "SELECT id, \"divers\", \"divers\", source, startSessionTime, \"Début de l’opération : \" || typeOperation FROM divers UNION SELECT id, \"divers\", \"divers\", source, endSessionTime, \"Fin de l'opération : \" || typeOperation || \", avec le statut : \" || exitStatut  FROM divers"
}
