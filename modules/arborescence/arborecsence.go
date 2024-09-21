package arborescence

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
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
	return Arborescence{}, nil
}
