package extraction

import (
	"aquarium/modules/extraction/divers"
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/getthis"
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
	"getthis":    getthis.Getthis{},
	"divers":     divers.Divers{},
}

func ListeExtracteursHtml(cheminProjet string) (map[string]string, error) {
	// On itère sur tous les extracteurs
	var resultat map[string]string = map[string]string{}
	bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
	if err != nil {
		log.Println(err)
		return map[string]string{}, err
	}
	defer func(bd *sql.DB) {
		err := bd.Close()
		if err != nil {

		}
	}(bd)
	requete, err := bd.Prepare("SELECT count(*) FROM chronologie WHERE extracteur=?;")
	if err != nil {
		return resultat, errors.New("Problème dans l'ouverture de la base de données d'analyse. \nAssurez vous que vous n'avez pas supprimé de fichiers ou recommencez une analyse. \n" + err.Error())
	}
	var nbLignes int
	for k, v := range liste_extracteurs {
		reponse, err := requete.Query(k)
		if err != nil {
			return resultat, errors.New("Problème dans l'ouverture de la base de données d'analyse. \nAssurez vous que vous n'avez pas supprimé de fichiers ou recommencez une analyse. \n" + err.Error())
		}
		defer func(reponse *sql.Rows) {
			err := reponse.Close()
			if err != nil {

			}
		}(reponse)
		reponse.Next()
		err = reponse.Scan(&nbLignes)
		if err != nil {
			return map[string]string{}, err
		}
		//log.Println(filepath.Join(cheminProjet, "collecteORC"))
		if v.PrerequisOK(filepath.Join(cheminProjet, "collecteORC")) && nbLignes == 0 {
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
