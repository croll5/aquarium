package extraction

import (
	"aquarium/modules/extraction/divers"
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/getthis"
	"aquarium/modules/extraction/navigateur"
	"aquarium/modules/extraction/sam"
	"aquarium/modules/extraction/werr"
	"errors"
	"path/filepath"
)

type Extracteur interface {
	Extraction(string) error
	Description() string
	PrerequisOK(string) bool
	CreationTable(string) error
	PourcentageChargement() int
}

var liste_extracteurs map[string]Extracteur = map[string]Extracteur{
	"evtx":       evtx.Evtx{},
	"navigateur": navigateur.Navigateur{},
	"werr":       werr.Werr{},
	"sam":        sam.Sam{},
	"getthis":    getthis.Getthis{},
	"divers":     divers.Divers{},
}

func ListeExtracteursHtml(cheminProjet string) (map[string]string, error) {
	// On itère sur tous les extracteurs
	var resultat map[string]string = map[string]string{}
	for k, v := range liste_extracteurs {
		//log.Println(filepath.Join(cheminProjet, "collecteORC"))
		if v.PrerequisOK(filepath.Join(cheminProjet, "collecteORC")) {
			resultat[k] = v.Description()
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
}
