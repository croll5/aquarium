/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l'analyse des collectes
traçologiques effectuées avec le logiciel DFIR-ORC.

Ce logiciel est régi par la licence CeCILL soumise au droit français et
respectant les principes de diffusion des logiciels libres. Vous pouvez
utiliser, modifier et/ou redistribuer ce programme sous les conditions
de la licence CeCILL telle que diffusée par le CEA, le CNRS et l'INRIA
sur le site "http://www.cecill.info".

En contrepartie de l'accessibilité au code source et des droits de copie,
de modification et de redistribution accordés par cette licence, il n'est
offert aux utilisateurs qu'une garantie limitée.  Pour les mêmes raisons,
seule une responsabilité restreinte pèse sur l'auteur du programme,  le
titulaire des droits patrimoniaux et les concédants successifs.

A cet égard  l'attention de l'utilisateur est attirée sur les risques
associés au chargement,  à l'utilisation,  à la modification et/ou au
développement et à la reproduction du logiciel par l'utilisateur étant
donné sa spécificité de logiciel libre, qui peut le rendre complexe à
manipuler et qui le réserve donc à des développeurs et des professionnels
avertis possédant  des  connaissances  informatiques approfondies.  Les
utilisateurs sont donc invités à charger  et  tester  l'adéquation  du
logiciel à leurs besoins dans des conditions permettant d'assurer la
sécurité de leurs systèmes et ou de leurs données et, plus généralement,
à l'utiliser et l'exploiter dans les mêmes conditions de sécurité.

Le fait que vous puissiez accéder à cet en-tête signifie que vous avez
pris connaissance de la licence CeCILL, et que vous en avez accepté les
termes.
*/

package arborescence

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
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

type MetaDonnees struct {
	Nom             string
	ADesEnfants     bool
	EnfantsSuspects int
	EnfantsInconnus int
	Empreinte       string
}

var cacheArbo Arborescence

// FONCTIONS INTERNES

/*
* Fonction qui rencoie un pointeur vers un fichier à partir de son chemin dans une
arborescence
@param chemin : le chemin du fichier à cherche dans l'arborescence
@param modele : un pointeur vers l'arborescence dans laquelle il faut chercher le fichier
@return : un pointeur vers le fichier, ou nil si le fichier n'est pas dans l'arborescence
*/
func chercherCheminDansModele(chemin []string, modele *Arborescence) *Arborescence {
	var vousetesiciModele *[]Arborescence = &modele.Enfants
	var res *Arborescence
	// On parcourt les dossiers du chemin vers le fichier à chercher
	for _, dossier := range chemin {
		var trouve bool = false
		// On parcourt les dossier du répertoire courant de l'arborescence
		for i := range *vousetesiciModele {
			// Si on trouve le dossier recherché dans l'arborescence, on change le
			// répertoire courant de l'arborescence vers ce dossier
			if (*vousetesiciModele)[i].Nom == dossier {
				res = &(*vousetesiciModele)[i]
				vousetesiciModele = &(*vousetesiciModele)[i].Enfants
				trouve = true
				break
			}
		}
		// Si on n'a pas trouvé le dossier recherché dans l'arborescence, c'est que le fichier
		// n'existe pas. On renvoie la valeur nulle
		if !trouve {
			return nil
		}
	}
	// Si on a trouvé le fichier, on renvoie un pointeur vers celui-ci
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
	// legitimite : variable indiquant si le fichier concerné est également dans le modèle (1), et si leurs empreintes
	// sont identiques (2).
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
	// On ajoute le nouveau fichier et ses caractéristiques dans l'arborescence
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
func remplitArborescenceDepuisCSV(fichierCSV *sevenzip.File, arbo *Arborescence, modeleArbo *Arborescence, colonneChemin int, colonneMD5 int, colonneParent int, colonneNom int) error {
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
		var cheminFichier string
		if colonneChemin == -1 {
			cheminFichier = strings.Join([]string{ligne[colonneParent], ligne[colonneNom]}, "\\")
		} else {
			cheminFichier = ligne[colonneChemin]
		}
		ajoutCheminDansArborescence(cheminFichier, ligne[colonneMD5], arbo, modeleArbo)
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

/** Fonction qui à partir du chemin vers un projet renvoie une arborescence si elle existe
 */
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

/*
* Fonction permettant de faire l'arborescence du système de fichier de la machine analysée
@param cheminProjet : le chemin vers le projet aquarium
@cheminModele : le chemin vers le modèle d'ORC avec lequel on compare l'arborescence
*/
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
	// On commence par lire la liste de tous les fichiers
	var cheminFichierNTFS string = filepath.Join(cheminProjet, "collecteORC", "Detail", "NTFSInfo_detail.7z")
	r, err := sevenzip.OpenReaderWithPassword(cheminFichierNTFS, "avproof")
	if err != nil {
		return resultatArbo, err
	}
	for _, fichier := range r.File {
		if strings.HasPrefix(fichier.Name, "NTFSInfo") {
			remplitArborescenceDepuisCSV(fichier, &resultatArbo, &modeleArbo, -1, 25, 3, 2)
		}
	} /*
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
							remplitArborescenceDepuisCSV(fichierCSV, &resultatArbo, &modeleArbo, 4, 7, -1, -1)
						}
					}
					return nil
				})
			}
		}*/
	// On enregistre l'arborescence que l'on vient d'extraire
	err = enregistrerArborescenceJson(&resultatArbo, filepath.Join(cheminProjet, "analyse", "arborescence.json"))
	return resultatArbo, err
}

/*
* Fonction qui renvoie les caractéristiques des fichiers et dossiers contenus un dossier
@param cheminProjet : le chemin vers le projet ORC
@param cheminDossier : chemin vers le dossier duquel on veut les enfants
@return : une liste contenant les métadonnées des éléments contenus dans le dossier
*/
func RecupEnfantsArbo(cheminProjet string, cheminDossier []int) ([]MetaDonnees, error) {
	var fichiers []MetaDonnees = []MetaDonnees{}
	// Si l'arborescence nextraite du fichier .json, on l'extrait
	if len(cacheArbo.Enfants) == 0 {
		var err error
		cacheArbo, err = GetArborescence(cheminProjet)
		if err != nil {
			return fichiers, err
		}
	}
	var vousetesici *Arborescence = &cacheArbo
	// On suit le chemin donné en paramètres pour se placer dans le dossier
	// duquel on veut le contenu
	for _, pas := range cheminDossier {
		if len(vousetesici.Enfants) < pas {
			return fichiers, errors.New("Le chemin de dossier spécifié est incohérent avec l'arborecsence")
		}
		vousetesici = &(*vousetesici).Enfants[pas]
	}
	// On parcourt les éléments de ce dossier
	for i := range vousetesici.Enfants {
		var legitimite []int = []int{0, 0}
		if (*vousetesici).Enfants[i].Legitimite == 0 {
			legitimite = []int{1, 0}
		} else if (*vousetesici).Enfants[i].Legitimite == 1 {
			legitimite = []int{0, 1}
		}
		var metadonnees MetaDonnees = MetaDonnees{
			Nom:             (*vousetesici).Enfants[i].Nom,
			ADesEnfants:     len((*vousetesici).Enfants[i].Enfants) != 0,
			EnfantsInconnus: legitimite[0],
			EnfantsSuspects: legitimite[1],
			Empreinte:       (*vousetesici).Enfants[i].EmpreinteMD5,
		}
		fichiers = append(fichiers, metadonnees)
	}
	return fichiers, nil
}
