package extraction

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/extraction/divers"
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/getthis"
	"aquarium/modules/extraction/navigateur"
	"aquarium/modules/extraction/prefetch"
	"aquarium/modules/extraction/sam"
	"aquarium/modules/extraction/werr"
	"errors"
	"log"
	"path/filepath"
	"time"
)

type Extracteur interface {
	Extraction(string) error
	Description() string
	PrerequisOK(string) bool
	CreationTable(string) error
	PourcentageChargement(string, bool) float32
	Annuler() bool
	DetailsEvenement(int) string
	SQLChronologie() string
}

type InfosExtracteur struct {
	Description string
	Progression float32
}

var liste_extracteurs map[string]Extracteur = map[string]Extracteur{
	"evtx":       evtx.Evtx{},
	"navigateur": navigateur.Navigateur{},
	"werr":       werr.Werr{},
	"sam":        sam.Sam{},
	"getthis":    getthis.Getthis{},
	"divers":     divers.Divers{},
	"prefetch":   prefetch.Prefetch{},
}

var colonnesTableChronologie map[string]string = map[string]string{"idEvt": "INT", "extracteur": "TEXT", "nomTable": "TEXT", "source": "TEXT", "horodatage": "DATETIME", "message": "TEXT"}

func ListeExtracteursHtml(cheminProjet string) (map[string]InfosExtracteur, error) {
	// On itère sur tous les extracteurs
	var resultat map[string]InfosExtracteur = map[string]InfosExtracteur{}
	for k, v := range liste_extracteurs {
		//log.Println(filepath.Join(cheminProjet, "collecteORC"))
		if v.PrerequisOK(filepath.Join(cheminProjet, "collecteORC")) {
			resultat[k] = InfosExtracteur{Description: v.Description(), Progression: v.PourcentageChargement(cheminProjet, true)}
		}
	}
	return resultat, nil
}

func Extraction(module string, cheminProjet string) error {
	if liste_extracteurs[module] == nil {
		return errors.New("Erreur : module " + module + " non reconnu")
	}
	//err := liste_extracteurs[module].Extraction(filepath.Join(cheminProjet, "collecteORC")) // Master AbdelMoad: commit 04a90c8ebc005011aae072aa56441a6d656b68db
	err := liste_extracteurs[module].Extraction(cheminProjet)
	return err
}

func CreationBaseAnalyse(cheminProjet string) {
	for _, extracteur := range liste_extracteurs {
		extracteur.CreationTable(cheminProjet)
	}
	var base aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	base.CreateTableIfNotExist2("chronologie", colonnesTableChronologie, false)
}

func ProgressionExtraction(cheminProjet string, idExtracteur string) float32 {
	return liste_extracteurs[idExtracteur].PourcentageChargement(cheminProjet, false)
}

func AnnulerExtraction(idExtracteur string) bool {
	ticker := time.NewTicker(500 * time.Millisecond)
	for _ = range ticker.C {
		if liste_extracteurs[idExtracteur].Annuler() {
			ticker.Stop()
			return true
		}
	}
	time.Sleep(30 * time.Second)
	return false
}

func DetailsEvenement(idExtracteur string, idEvenement int) string {
	return liste_extracteurs[idExtracteur].DetailsEvenement(idEvenement)
}

func ExtraireTableChronologie(cheminProjet string) error {
	var listeRequetesChronologie []string = []string{}
	for _, extracteur := range liste_extracteurs {
		if extracteur.SQLChronologie() != "" {
			listeRequetesChronologie = append(listeRequetesChronologie, extracteur.SQLChronologie())
		}
	}
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	err := base.RemplirTableDepuisRequetes("chronologie", listeRequetesChronologie, true, "horodatage")
	log.Println(err)
	return nil
}
