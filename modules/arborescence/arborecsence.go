package arborescence

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

type Arborescence struct {
	Nom          string         `json:"nom"`
	Enfants      []Arborescence `json:"enfants,omitempty"`
	EmpreinteMD5 string         `json:"md5,omitempty"`
	Legitimite   int            `json:"legitimite,omitempty"`
}

// FONCTIONS INTERNES

func chercherCheminDansModele(chemin []string, modele *Arborescence) *Arborescence {
	var vousetesiciModele *[]Arborescence = &modele.Enfants
	var res *Arborescence
	for _, dossier := range chemin {
		var trouve bool = false
		for i := range *vousetesiciModele {
			if (*vousetesiciModele)[i].Nom == dossier {
				res = &(*vousetesiciModele)[i]
				vousetesiciModele = &(*vousetesiciModele)[i].Enfants
				trouve = true
				break
			}
		}
		if !trouve {
			return nil
		}
	}
	return res
}

/*
* Fonction permettant d'ajouter un fichier dans une arborescence à partir de son chemin
@cheminFichier : le chemin du fichier à ajouter à l'arborescence
@arbo : un pointeur vers l'arborescence
@return : rien, modification de l'arborescence
*
*/
func ajoutCheminDansArborescence(cheminFichier string, md5Fichier string, arbo *Arborescence, modeleArbo *Arborescence) {
	// On supprime le \ au début du chemin pour éviter d'avoir une racine vide
	nomChemin, _ := strings.CutPrefix(cheminFichier, "\\")
	// On coupe le chemin en une liste de dossiers
	var chemin []string = strings.Split(nomChemin, "\\")
	var vousetesici *[]Arborescence = &arbo.Enfants
	// Pour chaque dossier/fichier du chemin, on l'ajoute s'il n'est pas encore
	// dans l'arborescence
	for _, dossier := range chemin[:int(math.Max(float64(len(chemin))-1, 0))] {
		var dossierExiste bool = false
		// On regarde dans tous les sous-dossiers du répertoire s'il y en a un du nom
		// du dossier que l'on veut rajouter (pour voir si le dossier existe déjà)
		for i := range *vousetesici {
			if (*vousetesici)[i].Nom == dossier {
				vousetesici = &(*vousetesici)[i].Enfants
				dossierExiste = true
				break
			}
		}
		// Si le dossier n'existe pas encore, on le crée
		if !dossierExiste {
			nouveauDossier := Arborescence{
				Nom:     dossier,
				Enfants: []Arborescence{},
			}
			*vousetesici = append(*vousetesici, nouveauDossier)
			// On se place dans le dossier concerné
			vousetesici = &(*vousetesici)[len(*vousetesici)-1].Enfants
		}
	}
	var legitimite int = -1
	var fichierDansModele *Arborescence = chercherCheminDansModele(chemin, modeleArbo)
	if fichierDansModele == nil {
		legitimite = 0
	} else {
		if fichierDansModele.EmpreinteMD5 == md5Fichier {
			legitimite = 2
		} else {
			legitimite = 1
		}
	}
	var nouveauFichier Arborescence = Arborescence{
		Nom:          chemin[len(chemin)-1],
		Enfants:      []Arborescence{},
		EmpreinteMD5: md5Fichier,
		Legitimite:   legitimite,
	}
	*vousetesici = append(*vousetesici, nouveauFichier)
}

/*
Fonction qui ajoute les fichiers contenus dans un fichier GetThis.csv à un arborescence
@fichierCSV : pointeur vers un fichier CSV duquel on veut récupérer les noms de fichiers
@arbo : arborescence que l'on veut remplir
@return : une erreur s'il y en a eu une
*/
func remplitArborescenceDepuisCSV(fichierCSV *sevenzip.File, arbo *Arborescence, modeleArbo *Arborescence) error {
	// On commence par ouvrir le fichier CSV
	contenuBrute, err := fichierCSV.Open()
	if err != nil {
		return err
	}
	// On lit son contenu
	var contenuCSV *csv.Reader = csv.NewReader(contenuBrute)
	// On ignore la première ligne, qui correspond au titre des colonnes
	contenuCSV.Read()
	// On ajoute les fichiers dans l'arborescence (4ème colonne du fichier GetThis.csv)
	for {
		ligne, err := contenuCSV.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ajoutCheminDansArborescence(ligne[4], ligne[7], arbo, modeleArbo)
	}
	return nil
}

/*
Fonction qui permet d'enrgistrer une arborescence dans un fichier json
@arbo : un pointeur vers l'arborescence à enregistrer
@chemin : le chemin du fichier json dans lequel enregistrer l'arborescence
*/
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

// FONCTIONS EXTERNES

func GetArborescence(cheminProjet string) (Arborescence, error) {
	var resultatArbo Arborescence
	donneesFichier, err := os.ReadFile(filepath.Join(cheminProjet, "analyse", "arborescence.json"))
	if err != nil {
		log.Println("WARN | Le fichier d'arborescence n'existe pas ou n'a pas pu être ouvert : ", err.Error())
		return Arborescence{}, nil
	}
	err = json.Unmarshal(donneesFichier, &resultatArbo)
	return resultatArbo, err
}

func ExtraireArborescence(cheminProjet string, cheminModele string) (Arborescence, error) {
	// Si un modèle a été donné en argument, on le récupère
	var modeleArbo Arborescence
	if cheminModele != "" {
		var err error
		modeleArbo, err = GetArborescence(cheminModele)
		if err != nil {
			return Arborescence{}, err
		}
	}
	var resultatArbo Arborescence = Arborescence{}
	resultatArbo.Nom = "racine"
	resultatArbo.Enfants = []Arborescence{}
	// On parcourt les fichiers GetTHis
	collectes, err := os.ReadDir(filepath.Join(cheminProjet, "collecteORC"))
	if err != nil {
		log.Println(err.Error())
		return resultatArbo, err
	}
	for _, collecte := range collectes {
		if collecte.IsDir() {
			filepath.Walk(filepath.Join(cheminProjet, "collecteORC", collecte.Name()), func(path string, info fs.FileInfo, err error) error {
				if filepath.Ext(path) != ".7z" {
					return nil
				}
				log.Println("INFO | Ouverture de l'archive ", path)
				r, err := sevenzip.OpenReaderWithPassword(path, "avproof")
				if err != nil {
					return err
				}
				for _, fichierCSV := range r.File {
					if fichierCSV.Name == "GetThis.csv" {
						remplitArborescenceDepuisCSV(fichierCSV, &resultatArbo, &modeleArbo)
					}
				}
				return nil
			})
		}
	}
	err = enregistrerArborescenceJson(&resultatArbo, filepath.Join(cheminProjet, "analyse", "arborescence.json"))
	return resultatArbo, err
}
