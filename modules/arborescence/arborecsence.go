package arborescence

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

type Arborescence struct {
	Nom     string         `json:"nom"`
	Enfants []Arborescence `json:"enfants,omitempty"`
}

func remplitArborescenceDepuisCSV(fichierCSV *sevenzip.File, arbo *Arborescence) error {
	log.Println("INFO | Traitement du fichier", fichierCSV.FileInfo().Name())
	contenuBrute, err := fichierCSV.Open()
	if err != nil {
		return err
	}
	var contenuCSV *csv.Reader = csv.NewReader(contenuBrute)
	for {
		ligne, err := contenuCSV.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		var chemin []string = strings.Split(ligne[4], "\\")
		var vousetesici *[]Arborescence = &arbo.Enfants

		for _, dossier := range chemin {
			var dossierExiste bool = false
			for i := range *vousetesici {
				if (*vousetesici)[i].Nom == dossier {
					vousetesici = &(*vousetesici)[i].Enfants
					dossierExiste = true
					break
				}
			}
			if !dossierExiste {
				nouveauDossier := Arborescence{
					Nom:     dossier,
					Enfants: []Arborescence{},
				}
				*vousetesici = append(*vousetesici, nouveauDossier)
				// Mettre à jour le pointeur pour le nouveau dossier
				vousetesici = &(*vousetesici)[len(*vousetesici)-1].Enfants
			}
		}

	}
	return nil
}

func enregistrerArborescenceJson(arbo *Arborescence, chemin string) error {
	donneesArbo, err := json.Marshal(arbo)
	if err != nil {
		return err
	}
	fichier, err := os.Create(chemin)
	if err != nil {
		return err
	}
	defer fichier.Close()
	_, err = fichier.Write(donneesArbo)
	return err
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
	resultatArbo.Nom = "racine"
	resultatArbo.Enfants = []Arborescence{}
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
				for _, fichierCSV := range r.File {
					if fichierCSV.Name == "GetThis.csv" {
						remplitArborescenceDepuisCSV(fichierCSV, &resultatArbo)
					}
				}
				return nil
			})
		}
	}
	err = enregistrerArborescenceJson(&resultatArbo, filepath.Join(filepath.Dir(cheminProjet), "analyse", "arborescence.json"))
	return resultatArbo, err
}
