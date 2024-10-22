package extraction

import (
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/navigateur"
	"aquarium/modules/extraction/sam"
	"aquarium/modules/extraction/werr"
	"database/sql"
	"errors"
	"log"
	"path/filepath"
)

type Extracteur interface {
	Extraction(string) error
	Description() string
	PrerequisOK(string) bool
}

var liste_extracteurs map[string]Extracteur = map[string]Extracteur{
	"evtx":       evtx.Evtx{},
	"navigateur": navigateur.Navigateur{},
	"werr":       werr.Werr{},
	"sam":        sam.Sam{},
}

func ListeExtracteursHtml(cheminProjet string) (map[string]string, error) {
	// On itère sur tous les extracteurs
	var resultat map[string]string = map[string]string{}
	bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
	if err != nil {
		log.Println(err)
		return map[string]string{}, err
	}
	defer bd.Close()
	requete, err := bd.Prepare("SELECT count(*) FROM chronologie WHERE extracteur=?;")
	var nbLignes int
	for k, v := range liste_extracteurs {
		reponse, err := requete.Query(k)
		defer reponse.Close()
		reponse.Next()
		reponse.Scan(&nbLignes)
		if err != nil {
			return map[string]string{}, err
		}
		if v.PrerequisOK(filepath.Join(cheminProjet, "collecteORC")) && nbLignes == 0 {
			resultat[k] = v.Description()
		}
	}
	return resultat, nil
}

func Extraction(module string, chemin_projet string) error {
	if liste_extracteurs[module] == nil {
		return errors.New("Erreur : module " + module + " non reconnu")
	}
	err := liste_extracteurs[module].Extraction(chemin_projet)
	return errn
}
