package arborescence

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
)

type Arborescence struct {
	Nom     string         `json:"nom"`
	Enfants []Arborescence `json:"enfants,omitempty"`
}

func GetArborescence(cheminProjet string) (Arborescence, error) {
	var resultatArbo Arborescence
	donneesFichier, err := os.ReadFile(filepath.Join(filepath.Dir(cheminProjet), "analyse", "arborescence.json"))
	if err != nil {
		log.Println("WARN | Le fichier d'arborescence n'existe pas ou n'a pas pu être ouvert : ", err.Error())
		return Arborescence{}, nil
	}
	err = json.Unmarshal(donneesFichier, &resultatArbo)
	return resultatArbo, err
}

func ExtraireArborescence(cheminProjet string) (Arborescence, error) {
	var resultatArbo Arborescence = Arborescence{}
	// On parcourt les fichiers GetTHis
	collectes, err := os.ReadDir(filepath.Join(filepath.Dir(cheminProjet), "collecteORC"))
	if err != nil {
		log.Println(err.Error())
		return resultatArbo, err
	}
	for _, collecte := range collectes {
		if collecte.IsDir() {
			log.Println("INFO | Parcourt du dossier ", collecte.Name())
			filepath.Walk(filepath.Join(filepath.Dir(cheminProjet), "collecteORC", collecte.Name()), func(path string, info fs.FileInfo, err error) error {
				log.Println("INFO | Ouverture de l'archive ", path, " qui a pour extension ", filepath.Ext(path))
				if filepath.Ext(path) != ".7z" {
					return nil
				}
				r, err := sevenzip.OpenReaderWithPassword(path, "avproof")
				if err != nil {
					return err
				}
				for _, f := range r.File {
					if f.Name == "GetThis.csv" {
						fmt.Println(f.FileInfo().Name())
					}
				}
				return nil
			})
		}
	}
	return resultatArbo, nil
}
