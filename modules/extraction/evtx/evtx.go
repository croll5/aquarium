package evtx

import (
	"aquarium/modules/gestionprojet"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xrawsec/golang-evtx/evtx"
	"github.com/bodgit/sevenzip"
	"github.com/pkg/errors"
)

// Requêtes SQL
var creationTable = "CREATE TABLE IF NOT EXISTS evtx(horodatage DATETIME, eventID int, eventRecordID int, processID int, threadID int, level int, providerGuid VARCHAR(36), providerName VARCHAR(50), task int, message TEXT, source TEXT, PRIMARY KEY (horodatage, eventRecordID, source))"
var recuperationMessageEvtx = "SELECT messages.message FROM messages INNER JOIN providers ON messages.provider_id = providers.id WHERE providers.name = ? AND messages.event_id = ?"
var ajoutEvenementDansBDD = "INSERT INTO Evtx VALUES "

type Evtx struct{}

// -------------------------- FONCTIONS LOCALES -------------------------- //

func ajouterGoEvtxMapDansBDD(evenement *evtx.GoEvtxMap, requeteInsertionEvtx *[]string, requeteMessagesEvtx *sql.Stmt, fichierSource string) error {
	var message string = ""
	var nbErreurs int = 0
	var chemin = evtx.GoEvtxPath{"Event", "System", "Provider", "Name"}
	provider, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		provider = "NaN"
		// log.Println("Problème dans la récupération du nom de provider", err)
	}
	// else {
	// 	resultat, err := requeteMessagesEvtx.Query(provider, evenement.EventID())
	// 	if err != nil {
	// 		// log.Println("Problème dans l'exécution de la requête SQL de récupération du message windows", err)
	// 		return err
	// 	}
	// 	defer resultat.Close()
	// 	resultat.Next()
	// 	resultat.Scan(&message)
	// }

	// Completion du message avec les informations de l'evenement
	chemin = evtx.GoEvtxPath{"Event", "EventData"}
	infosEvenement, err := evenement.Get(&chemin)
	if err != nil {
		nbErreurs++
		// log.Println("Problème dans la récupération des informations de l'évènemenemt Windows")
	} else {
		var infosEvenementJson []byte = evtx.ToJSON(infosEvenement)
		message += "\n" + strings.ReplaceAll(string(infosEvenementJson), "\"", "“")
	}
	// Ajout de l'évènement dans la base de données
	chemin = evtx.GoEvtxPath{"Event", "System", "Execution", "ProcessID"}
	processID, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		processID = "-1"
	}
	chemin = evtx.GoEvtxPath{"Event", "System", "Execution", "ThreadID"}
	threadID, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		threadID = "-1"
	}
	chemin = evtx.GoEvtxPath{"Event", "System", "Level"}
	level, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		level = "-1"
	}
	chemin = evtx.GoEvtxPath{"Event", "System", "Provider", "Guid"}
	providerGuid, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		providerGuid = "NaN"
	}
	chemin = evtx.GoEvtxPath{"Event", "System", "Task"}
	task, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		task = "-1"
	}
	if nbErreurs > 1 {
		message += "\n" + strings.ReplaceAll(string(evtx.ToJSON(evenement)), "\"", "“")
	}
	//_, err = requeteInsertionEvtx.
	var valeurs string = strings.Join([]string{"\"" + evenement.TimeCreated().String() + "\"", fmt.Sprint(evenement.EventID()), fmt.Sprint(evenement.EventRecordID()), processID, threadID, level, "\"" + providerGuid + "\"", "\"" + strings.ReplaceAll(provider, "\"", "“") + "\"", task, "\"" + message + "\"", "\"" + strings.ReplaceAll(fichierSource, "\"", "“") + "\""}, ",")
	*requeteInsertionEvtx = append(*requeteInsertionEvtx, "("+valeurs+")")
	return nil
}

func (e Evtx) extraireEvenementsDepuisFichier(cheminProjet string, fichier *sevenzip.File, bd *sql.DB, requeteMessagesEvtx *sql.Stmt, cheminTemp string, fichierSource string) error {
	err := gestionprojet.ExtraireFichierDepuis7z(fichier, cheminTemp)
	if err != nil {
		return err
	}
	if !strings.Contains(fichier.Name, ".evtx") {
		return nil
	}
	var fichierEvtx evtx.File
	fichierEvtx, err = evtx.OpenDirty(filepath.Join(cheminTemp, fichier.Name))
	if err != nil {
		return err
	}
	listeEvenements := fichierEvtx.FastEvents()
	var probleme error = nil
	var requeteInsertionEvtx []string = []string{}
	for evenement := range listeEvenements {
		err := ajouterGoEvtxMapDansBDD(evenement, &requeteInsertionEvtx, requeteMessagesEvtx, fichierSource)
		if err != nil {
			probleme = err
		}
	}
	if len(requeteInsertionEvtx) == 0 {
		return nil
	}
	var requeteInsertionEvtxFinale string = ajoutEvenementDansBDD + strings.Join(requeteInsertionEvtx, ",")
	//log.Println(requeteInsertionEvtxFinale)
	_, err = bd.Exec(requeteInsertionEvtxFinale)
	if err != nil {
		log.Println("--------------------------------------------------------------------------------")
		log.Println("----------------------------- ERREUR -------------------------------------------")
		log.Println("--------------------------------------------------------------------------------")
		log.Println(requeteInsertionEvtxFinale)
		log.Println("--------------------------------------------------------------------------------")
		log.Println("----------------------------- ERREUR -------------------------------------------")
		log.Println("--------------------------------------------------------------------------------")
		return err
	}
	return probleme
}

func (e Evtx) extraireEvementsDansDossier(cheminProjet string, cheminTemp string, nomDossier string, bd *sql.DB, requeteMessageEvtx *sql.Stmt) (error, error) {
	var probleme error = nil
	r, err := sevenzip.OpenReaderWithPassword(filepath.Join(cheminProjet, "collecteORC", nomDossier, "Event.7z"), "avproof")
	if err != nil {
		return err, nil
	}
	defer r.Close()
	for _, fichierEvtx := range r.File {
		var fichierSource string = filepath.Join(nomDossier, "Event", fichierEvtx.Name)
		err = e.extraireEvenementsDepuisFichier(cheminProjet, fichierEvtx, bd, requeteMessageEvtx, cheminTemp, fichierSource)
		if err != nil {
			probleme = err
		}
	}
	return nil, probleme
}

// ------------------------- FONCTIONS GLOBALES ------------------------- //

func (e Evtx) Extraction(cheminProjet string) error {
	// Ouverture de la base de données et création d'une nouvelle table
	bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
	if err != nil {
		return errors.New("Problème dans l'ouverture de la base de données d'analyse. \nAssurez vous que vous n'avez pas supprimé de fichiers ou recommencez une analyse. \n" + err.Error())
	}
	requete, err := bd.Prepare(creationTable)
	if err != nil {
		return err
	}
	requete.Exec()
	defer bd.Close()
	// Récupération du chemin de l'exécutable
	emplacementExecutable, err := os.Executable()
	if err != nil {
		return errors.New("Impossible d'atteindre la base de donnée des messages.")
	}
	emplacementExecutable, err = filepath.EvalSymlinks(emplacementExecutable)
	if err != nil {
		return errors.New("Impossible d'atteindre la base de donnée des messages.")
	}
	// Vérification de l'existence de la base de données
	_, err = os.Stat(filepath.Join(filepath.Dir(emplacementExecutable), "ressources", "messages_evtx.db"))
	if err != nil {
		return errors.New("La base des messages evtx Windows n'est pas présente dans le dossier " + filepath.Join(filepath.Dir(emplacementExecutable), "ressources") + "Ou n'a pas le bon nom (messages_evtx.db)\nVous pouvez la télécharger à l'adresse suivante https://github.com/Velocidex/evtx-data ou réinstaller le logiciel (https://github.com/croll5/aquarium/releases)")
	}
	// Ouverture de la base de données des messages evtx
	bdMessagesEvtx, err := sql.Open("sqlite", filepath.Join(filepath.Dir(emplacementExecutable), "ressources", "messages_evtx.db"))
	if err != nil {
		return errors.New("Problème dans l'ouverture de la base de données des messages windows. \nVérifiez que le fichier windows.10.enterprise.10.0.17763.amd64.db est bien présent dans votre dossier")
	}
	defer bdMessagesEvtx.Close()

	requeteMessageEvtx, err := bdMessagesEvtx.Prepare(recuperationMessageEvtx)
	if err != nil {
		log.Println("Problème dans la préparation de la requête SQL de récupération du message windows", err)
		return err
	}
	var probleme error = nil
	var cheminTemp string = filepath.Join(cheminProjet, "temp", "evtx")
	os.MkdirAll(cheminTemp, os.ModeDir)
	defer os.RemoveAll(cheminTemp)
	pbGeneral, probleme := e.extraireEvementsDansDossier(cheminProjet, cheminTemp, "General", bd, requeteMessageEvtx)
	pbLittle, err := e.extraireEvementsDansDossier(cheminProjet, cheminTemp, "Little", bd, requeteMessageEvtx)
	if err != nil {
		log.Println("J'ai trouvé le coupable")
		probleme = err
	}
	if pbGeneral != nil && pbLittle != nil {
		return errors.New("Les fichiers General\\Event.7z et Little\\Event.7z ne peuvent être ouverts.")
	}
	return probleme
}

func (e Evtx) Description() string {
	return "Évènements Windows (fichier .evtx)"
}

func (e Evtx) PrerequisOK(cheminCollecte string) bool {
	dossierGeneral, errGeneral := os.ReadDir(filepath.Join(cheminCollecte, "General"))
	dossierLittle, errLittle := os.ReadDir(filepath.Join(cheminCollecte, "Little"))
	if errGeneral == nil {
		for _, fichier := range dossierGeneral {
			if fichier.Name() == "Event.7z" {
				return true
			}
		}
	}
	if errLittle == nil {
		for _, fichier := range dossierLittle {
			if fichier.Name() == "Event.7z" {
				return true
			}
		}
	}
	return false
}
