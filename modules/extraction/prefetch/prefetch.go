package prefetch

import (
	"aquarium/modules/aquabase"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"www.velocidex.com/golang/go-prefetch"
)

type Prefetch struct{}

/* VARIABLE LOCALES */

var pourcentageChargement float32 = -1
var annulationDemandee bool = false
var annulationReussie bool = false

var colonnesTablePrefetch []string = []string{"executable", "fileSize", "hash", "runCount", "version"}
var colonnesTableFichierAccedesPrefetch []string = []string{"idFichier", "fileAccessed"}
var colonnesTableDernieresExecutionsPrefetch []string = []string{"idFichier", "dateExecution"}

/* FONCTIONS LOCALES */

func extraireInfosPrefetchDepuis7z(fichier *sevenzip.File, insertionPrefetch *aquabase.RequeteInsertion, insertionRessourcesPrefetch *aquabase.RequeteInsertion, insertionExecutionPrefetch *aquabase.RequeteInsertion, numFichier int) error {
	// On commence par ouvrir le fichier prefetch
	rc, err := fichier.Open()
	if err != nil {
		log.Println("Format de fichier non supporté : ", err.Error())
	}
	defer rc.Close()
	// Copie du contenu du fichier dans un tampon, pour pouvoir l'ouvrir avec l'extracteur de registres
	var tampon bytes.Buffer
	if _, err := io.Copy(&tampon, rc); err != nil {
		log.Println("Format de fichier non supporté : ", err.Error())
	}
	readerAt := bytes.NewReader(tampon.Bytes())
	// Ouverture du fichier avec la bibliothèque go-prefetch de Velocidex
	infosPrechargement, err := prefetch.LoadPrefetch(readerAt)
	if err != nil {
		return err
	}
	// On ajoute les informations sur le fichier dans la base de données
	insertionPrefetch.AjouterDansRequete(infosPrechargement.Executable, infosPrechargement.FileSize, infosPrechargement.Hash, infosPrechargement.RunCount, infosPrechargement.Version, numFichier)
	// On ajoute les dates d'execution dans la table des executions
	for _, execution := range infosPrechargement.LastRunTimes {
		insertionExecutionPrefetch.AjouterDansRequete(numFichier, execution.Local())
	}
	// On ajoute les ressources dans la table des ressources
	for _, ressource := range infosPrechargement.FilesAccessed {
		insertionRessourcesPrefetch.AjouterDansRequete(numFichier, ressource)
	}
	return nil
}

func annulerExtraction(cheminProjet string) error {
	base := aquabase.InitDB_Extraction(cheminProjet)
	err := base.RemoveFromWhere("prefetch", "1=1")
	if err != nil {
		return err
	}
	err = base.RemoveFromWhere("executionPrefetch", "1=1")
	if err != nil {
		return err
	}
	err = base.RemoveFromWhere("ressourcesPrefetch", "1=1")
	return err
}

/* FONCTIONS REQUISES PAR LE MODULE EXTRACTEUR */

func (pref Prefetch) Extraction(cheminProjet string) error {
	var insertionPrefetch aquabase.RequeteInsertion = aquabase.InitRequeteInsertionExtraction("prefetch", append(colonnesTablePrefetch, "id"))
	var insertionExecutionPrefetch aquabase.RequeteInsertion = aquabase.InitRequeteInsertionExtraction("executionPrefetch", colonnesTableDernieresExecutionsPrefetch)
	var insertionRessourcesPrefetch aquabase.RequeteInsertion = aquabase.InitRequeteInsertionExtraction("ressourcesPrefetch", colonnesTableFichierAccedesPrefetch)
	var numFichier = 0
	dossierArtefact, err := sevenzip.OpenReader(filepath.Join(cheminProjet, "CollecteORC", "General", "Artefacts.7z"))
	if err == nil {
		for fichiersTraites, artefact := range dossierArtefact.File {
			if annulationDemandee {
				err := annulerExtraction(cheminProjet)
				if err == nil {
					annulationReussie = true
					return nil
				}
			}
			if strings.Contains(artefact.Name, "Prefetch") {
				extraireInfosPrefetchDepuis7z(artefact, &insertionPrefetch, &insertionRessourcesPrefetch, &insertionExecutionPrefetch, numFichier)
				numFichier++
			}
			pourcentageChargement = float32(fichiersTraites*100) / float32(len(dossierArtefact.File))
		}
		dossierArtefact.Close()
	}
	insertionPrefetch.Executer(cheminProjet)
	insertionExecutionPrefetch.Executer(cheminProjet)
	insertionRessourcesPrefetch.Executer(cheminProjet)
	pourcentageChargement = 101
	return nil
}

func (pref Prefetch) Description() string {
	return "Fichiers de préchargement (contenant des informations sur l'exécution d'applications)"
}

func (pref Prefetch) PrerequisOK(cheminCollecte string) bool {
	dossierGeneral, err := os.ReadDir(filepath.Join(cheminCollecte, "General"))
	if err == nil {
		for _, fichier := range dossierGeneral {
			if fichier.Name() == "Artefacts.7z" {
				dossierArtefact, err := sevenzip.OpenReader(filepath.Join(cheminCollecte, "General", "Artefacts.7z"))
				if err != nil {
					return false
				}
				for _, artefact := range dossierArtefact.File {
					if strings.Contains(artefact.Name, "Prefetch") {
						return true
					}
				}
				dossierArtefact.Close()
			}
		}
	}
	return false
}

func (pref Prefetch) CreationTable(cheminProjet string) error {
	base := aquabase.InitDB_Extraction(cheminProjet)
	base.CreateTableIfNotExist1("prefetch", colonnesTablePrefetch, true)
	base.CreateTableIfNotExist1("executionPrefetch", colonnesTableDernieresExecutionsPrefetch, true)
	base.CreateTableIfNotExist1("ressourcesPrefetch", colonnesTableFichierAccedesPrefetch, true)
	return nil
}

func (pref Prefetch) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if pourcentageChargement == -1 && verifierTableVide {
		base := aquabase.InitDB_Extraction(cheminProjet)
		if !base.EstTableVide("prefetch") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (pref Prefetch) Annuler() bool {
	if !annulationDemandee {
		annulationDemandee = true
		annulationReussie = false
	}
	if annulationReussie {
		annulationDemandee = false
	}
	return annulationReussie
}

func (pref Prefetch) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}

func (pref Prefetch) SQLChronologie() string {
	return ""
}
