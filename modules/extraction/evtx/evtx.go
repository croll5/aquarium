package evtx

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/extraction/utilitaires"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xrawsec/golang-evtx/evtx"
	"github.com/bodgit/sevenzip"
	"github.com/pkg/errors"
)

type Evtx struct{}

var progressionChargement float32 = -1
var demandeInterruption bool = false
var interruptionReussie bool = false

var colonnesTableEvtx = []string{"horodatage", "eventID", "eventRecordID", "processID", "threadID", "level", "providerGuid", "providerName", "task", "message", "source"}

// -------------------------- FONCTIONS LOCALES -------------------------- //

/*
	Fonction qui, à partir d'un évènement, va ajouter à la requête ses caractéristiques

@param evenement : un pointeur vers l'évènement à ajouter
@param requeteInsertionEvtx : la requete de base de données à laquelle on veut l'ajouter
@param fichierSource : le chemin vers le fichier source
@return : une erreur s'il y a eu des problèmes dans l'extraction des caractéristiques de l'évènement
*/
func ajouterGoEvtxMapDansBDD(evenement *evtx.GoEvtxMap, requeteInsertionEvtx *aquabase.RequeteInsertion, fichierSource string) error {
	var message string = ""
	var nbErreurs int = 0
	var chemin = evtx.GoEvtxPath{"Event", "System", "Provider", "Name"}
	provider, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		provider = "NaN"
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
	} else {
		var infosEvenementJson []byte = evtx.ToJSON(infosEvenement)
		message = string(infosEvenementJson)
	}
	// Récupération du processID
	chemin = evtx.GoEvtxPath{"Event", "System", "Execution", "ProcessID"}
	processID, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		processID = "-1"
	}
	// Récupération du ThreadID
	chemin = evtx.GoEvtxPath{"Event", "System", "Execution", "ThreadID"}
	threadID, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		threadID = "-1"
	}
	// Récupération du niveau d'alerte
	chemin = evtx.GoEvtxPath{"Event", "System", "Level"}
	level, err := evenement.GetString(&chemin)
	if err != nil {
		nbErreurs++
		level = "-1"
	}
	// Récupération de l'identifiant du provider
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
		message = string(evtx.ToJSON(evenement))
	}
	// Concaténation de toutes les informations précédemment récupérées
	requeteInsertionEvtx.AjouterDansRequete(evenement.TimeCreated(), evenement.EventID(), evenement.EventRecordID(), processID, threadID, level, providerGuid, provider, task, message, fichierSource)
	return nil
}

/*
Fonction qui, à partir d'un fichier evtx zippé, ajoute tous ses évènements à la base de données

@param cheminProjet : le chemin de l'analyse aquarium
@param fichier : le fichier d'évènements
@param bd : un pointeur vers la base de données d'analyse
@param cheminTemp : le chemin vers un répertoire temporaire
@param fichierSource : le chemin du fichier evtx à extraire
*/
func (e Evtx) extraireEvenementsDepuisFichier(cheminProjet string, fichier *sevenzip.File, cheminTemp string, fichierSource string) error {
	// On commence par copier le fichier concerné dans le dossier temporaire
	err := utilitaires.ExtraireFichierDepuis7z(fichier, cheminTemp)
	if err != nil {
		return err
	}
	// On abandonne si ce n'est pas un fichier evtx
	if !strings.Contains(fichier.Name, ".evtx") {
		return nil
	}
	// On ouvre le fichier avec la bibliothèque evtx
	var fichierEvtx evtx.File
	fichierEvtx, err = evtx.OpenDirty(filepath.Join(cheminTemp, fichier.Name))
	if err != nil {
		return err
	}
	// On récupère la liste des évènements
	listeEvenements := fichierEvtx.FastEvents()
	var probleme error = nil
	var requeteInsertionEvtx aquabase.RequeteInsertion = aquabase.InitRequeteInsertionExtraction("Evtx", colonnesTableEvtx)
	for evenement := range listeEvenements {
		// On ajoute chaque évènement à la requete
		err := ajouterGoEvtxMapDansBDD(evenement, &requeteInsertionEvtx, fichierSource)
		if err != nil {
			probleme = err
		}
	}
	// On exécute la requete
	err = requeteInsertionEvtx.Executer(cheminProjet)
	// Si l'on n'a pas pu l'exécuter, on renvoie une erreur
	if err != nil {
		log.Println("ERROR - extraireEvenementDepuisFichier : ", err)
		return err
	}
	return probleme
}

/*
Fonction qui extrait tous les fichiers evtx d'un dossier donnée
@param cheminProjet : le chemin vers l'analyse aquarium
@param cheminTemp : le chemin vers un dossier temporaire dans lequel seront dézippés les fichiers evtx
@param nomDossier : le nom du dossier duquel extraire les évènements
@param bd : un pointeur vers la base de données dans laquelle écrire les évènements
*/
func (e Evtx) extraireEvementsDansDossier(cheminProjet string, cheminTemp string, nomDossier string) (error, error) {
	var probleme error = nil
	// On ouvre le dossier compressé avec la bibliothèque 7zip
	r, err := sevenzip.OpenReaderWithPassword(filepath.Join(cheminProjet, "collecteORC", nomDossier, "Event.7z"), "avproof")
	if err != nil {
		return err, nil
	}
	defer r.Close()
	// On parcourt tous les fichiers du dossier compressé et on les met dans la base de données
	var nbFichiers float32 = float32(len(r.File))
	for numFichier, fichierEvtx := range r.File {
		if demandeInterruption {
			err = interruptionExtracteur(cheminProjet)
			if err == nil {
				interruptionReussie = true
				return nil, errors.New("Operation Annulee")
			}
		}
		var fichierSource string = filepath.Join(nomDossier, "Event", fichierEvtx.Name)
		err = e.extraireEvenementsDepuisFichier(cheminProjet, fichierEvtx, cheminTemp, fichierSource)
		if err != nil {
			probleme = err
		}
		progressionChargement = float32(numFichier*100) / nbFichiers
	}
	return nil, probleme
}

func interruptionExtracteur(cheminProjet string) error {
	bdd := aquabase.InitDB_Extraction(cheminProjet)
	var err error = bdd.RemoveFromWhere("evtx", "1=1")
	progressionChargement = -1
	return err
}

// ------------------------- FONCTIONS GLOBALES ------------------------- //

/* Fonction d'extraction des fichiers evtx */
func (e Evtx) Extraction(cheminProjet string) error {
	progressionChargement = 0
	// Création d'une nouvelle table
	// Récupération du chemin de l'exécutable
	// emplacementExecutable, err := os.Executable()
	// if err != nil {
	// 	return errors.New("Impossible d'atteindre la base de donnée des messages.")
	// }
	// emplacementExecutable, err = filepath.EvalSymlinks(emplacementExecutable)
	// if err != nil {
	// 	return errors.New("Impossible d'atteindre la base de donnée des messages.")
	// }
	// Vérification de l'existence de la base de données
	// _, err = os.Stat(filepath.Join(filepath.Dir(emplacementExecutable), "ressources", "messages_evtx.db"))
	// if err != nil {
	// 	return errors.New("La base des messages evtx Windows n'est pas présente dans le dossier " + filepath.Join(filepath.Dir(emplacementExecutable), "ressources") + "Ou n'a pas le bon nom (messages_evtx.db)\nVous pouvez la télécharger à l'adresse suivante https://github.com/Velocidex/evtx-data ou réinstaller le logiciel (https://github.com/croll5/aquarium/releases)")
	// }
	// // Ouverture de la base de données des messages evtx
	// bdMessagesEvtx, err := sql.Open("sqlite", filepath.Join(filepath.Dir(emplacementExecutable), "ressources", "messages_evtx.db"))
	// if err != nil {
	// 	return errors.New("Problème dans l'ouverture de la base de données des messages windows. \nVérifiez que le fichier windows.10.enterprise.10.0.17763.amd64.db est bien présent dans votre dossier")
	// }
	// defer bdMessagesEvtx.Close()

	// requeteMessageEvtx, err := bdMessagesEvtx.Prepare(recuperationMessageEvtx)
	// if err != nil {
	// 	log.Println("Problème dans la préparation de la requête SQL de récupération du message windows", err)
	// 	return err
	// }
	var probleme error = nil
	var cheminTemp string = filepath.Join(cheminProjet, "temp", "evtx")
	err := os.MkdirAll(cheminTemp, os.ModeDir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(cheminTemp)
	pasDossier, probleme := e.extraireEvementsDansDossier(cheminProjet, cheminTemp, "General")
	if pasDossier != nil {
		// Tous les évènements de Little/evtx sont aussi dans General/evtx, il est donc inutile de les rééxtraire
		pasDossier, probleme = e.extraireEvementsDansDossier(cheminProjet, cheminTemp, "Little")
		if pasDossier != nil {
			return errors.New("Les fichiers General\\Event.7z et Little\\Event.7z ne peuvent être ouverts.")
		}
	}
	log.Println(probleme)
	if probleme == nil || probleme.Error() != "Operation Annulee" {

		progressionChargement = 101
	} else {
		progressionChargement = -1
	}
	return nil
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

func (e Evtx) CreationTable(cheminProjet string) error {
	base := aquabase.InitDB_Extraction(cheminProjet)
	err := base.CreateTableIfNotExist1("evtx", colonnesTableEvtx, true)
	return err
}

func (e Evtx) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if progressionChargement == -1 && verifierTableVide {
		// On véfifie que l'extracteur n'a pas déjà été chargé
		base := aquabase.InitDB_Extraction(cheminProjet)
		if base.EstTableVide("evtx") {
			return -1
		} else {
			return 100
		}
	}
	return progressionChargement
}

func (e Evtx) Annuler() bool {
	if demandeInterruption {
		if interruptionReussie {
			demandeInterruption = false
		}
		return interruptionReussie
	} else {
		interruptionReussie = false
		demandeInterruption = true
	}
	return false
}

func (e Evtx) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}
