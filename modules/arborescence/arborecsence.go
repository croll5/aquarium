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
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const NOM_CHEMIN_FICHIER string = "chemin_fichier"

const DOSSIER_ARBO string = "arborescences"

type Arborescence struct {
	Enfants      map[string]*Arborescence `json:"enfants,omitempty"`
	EmpreinteMD5 string                   `json:"md5,omitempty"`
	Legitimite   int                      `json:"legitimite,omitempty"`
	Fichiers     []DonneesFichier         `json:"fichiers,omitempty"`
}

type DonneesFichier struct {
	Source int64 `json:"source"`
	Copie  int64 `json:"copie,omitempty"`
}

type MetaDonnees struct {
	Nom             string
	ADesEnfants     bool
	EnfantsSuspects int
	EnfantsInconnus int
	Empreinte       string
	IdSource        int64
	IdCopie         int64
}

var AQUA_MACHINE = "[AQUA_MACHINE]"

var cacheArbo Arborescence

var nomArboEnCache = ""
var enCoursDextraction bool = false

func ExtractionEnCours() bool {
	return enCoursDextraction
}

func ArborescenceEnCache() string {
	return nomArboEnCache
}

func ExtraireArborescence(cheminProjet string, cheminModele string, idMachine string) error {
	enCoursDextraction = true
	var cache Arborescence = Arborescence{}
	// On cherche la configuration du projet
	configAnalyse, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return err
	}
	// On lit la configuration de la machine
	configMachine, err := config.GetConfigurationMachine(cheminProjet, idMachine, configAnalyse.Machines[idMachine])
	if err != nil {
		return err
	}
	// On commence par lister les fichiers
	var abd *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	requeteSQL := strings.ReplaceAll(configMachine.Arborescence.RequeteCreation, AQUA_MACHINE, idMachine)
	resultatRequete, err := abd.ExecuterRequeteSQL(requeteSQL)
	if err != nil {
		return err
	}
	// On parcourt les résultats de la requête
	i := 0
	for {
		i++
		fichier, ok := resultatRequete.Suivant()
		if !ok {
			break
		}
		ajouterFichierDansArbo(fichier, configMachine.Arborescence, &cache)
		if i%10_000 == 0 {
			log.Println(i)
		}
	}
	// Enregistrement de l’analyse
	err = os.MkdirAll(filepath.Join(cheminProjet, config.DOSSIER_ANALYSE, DOSSIER_ARBO), os.ModeAppend)
	if err != nil {
		return err
	}
	err = enregistrerArborescenceJson(&cache, filepath.Join(cheminProjet, config.DOSSIER_ANALYSE, DOSSIER_ARBO, idMachine)+".json")
	cacheArbo = cache
	nomArboEnCache = idMachine
	enCoursDextraction = false
	return err
}

func RecupEnfantsArbo(cheminProjet string, cheminDossier []string, idMachine string) ([]MetaDonnees, error) {
	// On commence par récupérer la configuration de la machine
	aquaConfigMachine, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return []MetaDonnees{}, err
	}
	configAnalyse, err := config.GetConfigurationMachine(cheminProjet, idMachine, aquaConfigMachine.Machines[idMachine])
	if err != nil {
		return []MetaDonnees{}, err
	}
	// Si l'arborescence en cache ne correspond pas, on l’extrait à nouveau
	if idMachine != nomArboEnCache {
		// On essaie d’ouvrir le fichier de l’arborescence
		contenuFichierArbo, err := os.ReadFile(filepath.Join(cheminProjet, config.DOSSIER_ANALYSE, DOSSIER_ARBO, idMachine) + ".json")
		if err != nil {
			err = ExtraireArborescence(cheminProjet, "", idMachine)
			if err != nil {
				return []MetaDonnees{}, err
			}
			return getContenuDossier(cheminProjet, cheminDossier, configAnalyse.Arborescence)
		}
		// Si l’arborescence existe déjà, on la charge dans le cache
		err = json.Unmarshal(contenuFichierArbo, &cacheArbo)
		if err != nil {
			log.Println("ERREUR : ", err)
			return []MetaDonnees{}, err
		}
		nomArboEnCache = idMachine
	}
	return getContenuDossier(cheminProjet, cheminDossier, configAnalyse.Arborescence)
}

/* ------------------------------------ Fonctions internes ------------------------------------ */

func ajouterFichierDansArbo(infosFichier map[string]interface{}, configArbo config.ConfigArborescence, cache *Arborescence) error {
	// On commence par couper le chemin du fichier en dossiers
	var dossiers []string = make([]string, 0)
	switch infosFichier[NOM_CHEMIN_FICHIER].(type) {
	case string:
		dossiers = strings.Split(infosFichier[NOM_CHEMIN_FICHIER].(string), configArbo.SeparateurDossier)
	default:
		return errors.New("Problème dans la requête de sélection de l’arborescence")
	}
	// On parcourt ensuite les dossiers
	debut := 0
	if len(dossiers) > 0 && dossiers[0] == "" {
		debut = 1
	}
	var positionDansCache *Arborescence = cache
	for i := debut; i < len(dossiers)-1; i++ {
		_, trouve := (*positionDansCache).Enfants[dossiers[i]]
		// Si l’on n’a pas trouvé le dossier, on le crée
		if trouve {
			positionDansCache = (*positionDansCache).Enfants[dossiers[i]]
		} else {
			nouveauDossier := Arborescence{Enfants: map[string]*Arborescence{}}
			if (*positionDansCache).Enfants == nil {
				(*positionDansCache).Enfants = map[string]*Arborescence{}

			}
			(*positionDansCache).Enfants[dossiers[i]] = &nouveauDossier
			positionDansCache = &nouveauDossier
		}
	}
	// On crée les données du fichier
	donneesFichier := DonneesFichier{}
	// On ajoute le fichier
	switch infosFichier["idSource"].(type) {
	case int:
		donneesFichier.Source = int64(infosFichier["idSource"].(int))
	case int64:
		donneesFichier.Source = infosFichier["idSource"].(int64)
	default:
		return errors.New("[AQUA_ERR] : FOnction ajouterFichierDansArbo")
	}
	// On ajoute éventuellement le fichier collecté
	switch infosFichier["idCopie"].(type) {
	case int:
		donneesFichier.Copie = int64(infosFichier["idCopie"].(int))
	case int64:
		donneesFichier.Copie = infosFichier["idCopie"].(int64)
	default:
		break
	}
	(*positionDansCache).Fichiers = append((*positionDansCache).Fichiers, donneesFichier)
	return nil
}

func getContenuDossier(cheminProjet string, cheminDossier []string, configArborescence config.ConfigArborescence) ([]MetaDonnees, error) {
	var positionDansArborescence *Arborescence = &cacheArbo
	var existe bool
	for _, dossier := range cheminDossier {
		positionDansArborescence, existe = positionDansArborescence.Enfants[dossier]
		if !existe {
			return []MetaDonnees{}, os.ErrNotExist
		}
	}
	var metaDonnees []MetaDonnees = make([]MetaDonnees, len(positionDansArborescence.Enfants)+len(positionDansArborescence.Fichiers))
	var i int = 0
	// On commence par ajouter les dossiers
	for nomDossier := range positionDansArborescence.Enfants {
		metaDonnees[i] = MetaDonnees{Nom: nomDossier, ADesEnfants: true}
		i++
	}
	// On ouvre la base de données
	adb := aquabase.InitDB_Extraction(cheminProjet)
	for _, fichier := range positionDansArborescence.Fichiers {
		resultat := adb.SelectFrom(configArborescence.RequeteMetadonnees, fichier.Source)
		if len(resultat) < 1 {
			return []MetaDonnees{}, errors.ErrUnsupported
		}
		switch resultat[0]["Nom"].(type) {
		case string:
			metaDonnees[i] = MetaDonnees{Nom: resultat[0]["Nom"].(string), IdSource: fichier.Source, IdCopie: fichier.Copie}
		default:
			return []MetaDonnees{}, errors.ErrUnsupported
		}
		i++
	}
	return metaDonnees, nil
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

/* /////////////////////////////////////////////////////////////////////////////////////////////////////// */
/* /////////////////////////////////////////////////////////////////////////////////////////////////////// */
/* ///////////////////////////////////////// ANCIENNES FONCTIONS ///////////////////////////////////////// */
/* /////////////////////////////////////////////////////////////////////////////////////////////////////// */
/* /////////////////////////////////////////////////////////////////////////////////////////////////////// */

// FONCTIONS INTERNES

/*
* Fonction qui rencoie un pointeur vers un fichier à partir de son chemin dans une
arborescence
@param chemin : le chemin du fichier à cherche dans l'arborescence
@param modele : un pointeur vers l'arborescence dans laquelle il faut chercher le fichier
@return : un pointeur vers le fichier, ou nil si le fichier n'est pas dans l'arborescence
// */
// func archive_chercherCheminDansModele(chemin []string, modele *Arborescence) *Arborescence {
// 	var vousetesiciModele *[]Arborescence = &modele.Enfants
// 	var res *Arborescence
// 	// On parcourt les dossiers du chemin vers le fichier à chercher
// 	for _, dossier := range chemin {
// 		var trouve bool = false
// 		// On parcourt les dossier du répertoire courant de l'arborescence
// 		for i := range *vousetesiciModele {
// 			// Si on trouve le dossier recherché dans l'arborescence, on change le
// 			// répertoire courant de l'arborescence vers ce dossier
// 			if (*vousetesiciModele)[i].Nom == dossier {
// 				res = &(*vousetesiciModele)[i]
// 				vousetesiciModele = &(*vousetesiciModele)[i].Enfants
// 				trouve = true
// 				break
// 			}
// 		}
// 		// Si on n'a pas trouvé le dossier recherché dans l'arborescence, c'est que le fichier
// 		// n'existe pas. On renvoie la valeur nulle
// 		if !trouve {
// 			return nil
// 		}
// 	}
// 	// Si on a trouvé le fichier, on renvoie un pointeur vers celui-ci
// 	return res
// }

/*
* Fonction permettant d'ajouter un fichier dans une arborescence à partir de son chemin
@cheminFichier : le chemin du fichier à ajouter à l'arborescence
@arbo : un pointeur vers l'arborescence
@return : rien, modification de l'arborescence
// *
// */
// func archive_ajoutCheminDansArborescence(cheminFichier string, md5Fichier string, arbo *Arborescence, modeleArbo *Arborescence) {
// 	// On supprime le \ au début du chemin pour éviter d'avoir une racine vide
// 	nomChemin, _ := strings.CutPrefix(cheminFichier, "\\")
// 	// On coupe le chemin en une liste de dossiers
// 	var chemin []string = strings.Split(nomChemin, "\\")
// 	var vousetesici *[]Arborescence = &arbo.Enfants
// 	// Pour chaque dossier/fichier du chemin, on l'ajoute s'il n'est pas encore
// 	// dans l'arborescence
// 	for _, dossier := range chemin[:int(math.Max(float64(len(chemin))-1, 0))] {
// 		var dossierExiste bool = false
// 		// On regarde dans tous les sous-dossiers du répertoire s'il y en a un du nom
// 		// du dossier que l'on veut rajouter (pour voir si le dossier existe déjà)
// 		for i := range *vousetesici {
// 			if (*vousetesici)[i].Nom == dossier {
// 				vousetesici = &(*vousetesici)[i].Enfants
// 				dossierExiste = true
// 				break
// 			}
// 		}
// 		// Si le dossier n'existe pas encore, on le crée
// 		if !dossierExiste {
// 			nouveauDossier := Arborescence{
// 				Nom:     dossier,
// 				Enfants: []Arborescence{},
// 			}
// 			*vousetesici = append(*vousetesici, nouveauDossier)
// 			// On se place dans le dossier concerné
// 			vousetesici = &(*vousetesici)[len(*vousetesici)-1].Enfants
// 		}
// 	}
// 	// legitimite : variable indiquant si le fichier concerné est également dans le modèle (1), et si leurs empreintes
// 	// sont identiques (2).
// 	var legitimite int = -1
// 	var fichierDansModele *Arborescence = archive_chercherCheminDansModele(chemin, modeleArbo)
// 	if fichierDansModele == nil {
// 		legitimite = 0
// 	} else {
// 		if fichierDansModele.EmpreinteMD5 == md5Fichier {
// 			legitimite = 2
// 		} else {
// 			legitimite = 1
// 		}
// 	}
// 	// On ajoute le nouveau fichier et ses caractéristiques dans l'arborescence
// 	var nouveauFichier Arborescence = Arborescence{
// 		Nom:          chemin[len(chemin)-1],
// 		Enfants:      []Arborescence{},
// 		EmpreinteMD5: md5Fichier,
// 		Legitimite:   legitimite,
// 	}
// 	*vousetesici = append(*vousetesici, nouveauFichier)
// }

/*
Fonction qui ajoute les fichiers contenus dans un fichier GetThis.csv à un arborescence
@fichierCSV : pointeur vers un fichier CSV duquel on veut récupérer les noms de fichiers
@arbo : arborescence que l'on veut remplir
@return : une erreur s'il y en a eu une
// */
// func archive_remplitArborescenceDepuisCSV(fichierCSV *sevenzip.File, arbo *Arborescence, modeleArbo *Arborescence) error {
// 	// On commence par ouvrir le fichier CSV
// 	contenuBrute, err := fichierCSV.Open()
// 	if err != nil {
// 		return err
// 	}
// 	// On lit son contenu
// 	var contenuCSV *csv.Reader = csv.NewReader(contenuBrute)
// 	// On ignore la première ligne, qui correspond au titre des colonnes
// 	contenuCSV.Read()
// 	// On ajoute les fichiers dans l'arborescence (4ème colonne du fichier GetThis.csv)
// 	for {
// 		ligne, err := contenuCSV.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		archive_ajoutCheminDansArborescence(ligne[4], ligne[7], arbo, modeleArbo)
// 	}
// 	return nil
// }

// FONCTIONS EXTERNES

/** Fonction qui à partir du chemin vers un projet renvoie une arborescence si elle existe
//  */
// func archive_GetArborescence(cheminProjet string) (Arborescence, error) {
// 	var resultatArbo Arborescence
// 	donneesFichier, err := os.ReadFile(filepath.Join(cheminProjet, "analyse", "arborescence.json"))
// 	if err != nil {
// 		log.Println("WARN | Le fichier d'arborescence n'existe pas ou n'a pas pu être ouvert : ", err.Error())
// 		return Arborescence{}, nil
// 	}
// 	err = json.Unmarshal(donneesFichier, &resultatArbo)
// 	return resultatArbo, err
// }

/*
* Fonction permettant de faire l'arborescence du système de fichier de la machine analysée
@param cheminProjet : le chemin vers le projet aquarium
@cheminModele : le chemin vers le modèle d'ORC avec lequel on compare l'arborescence
// */
// func archive_ExtraireArborescence(cheminProjet string, cheminModele string) (Arborescence, error) {
// 	// Si un modèle a été donné en argument, on le récupère
// 	var modeleArbo Arborescence
// 	if cheminModele != "" {
// 		var err error
// 		modeleArbo, err = archive_GetArborescence(cheminModele)
// 		if err != nil {
// 			return Arborescence{}, err
// 		}
// 	}
// 	var resultatArbo Arborescence = Arborescence{}
// 	resultatArbo.Nom = "racine"
// 	resultatArbo.Enfants = []Arborescence{}
// 	// On parcourt les fichiers GetTHis
// 	collectes, err := os.ReadDir(filepath.Join(cheminProjet, "collecteORC"))
// 	if err != nil {
// 		log.Println(err.Error())
// 		return resultatArbo, err
// 	}
// 	for _, collecte := range collectes {
// 		if collecte.IsDir() {
// 			filepath.Walk(filepath.Join(cheminProjet, "collecteORC", collecte.Name()), func(path string, info fs.FileInfo, err error) error {
// 				if filepath.Ext(path) != ".7z" {
// 					return nil
// 				}
// 				log.Println("INFO | Ouverture de l'archive ", path)
// 				r, err := sevenzip.OpenReaderWithPassword(path, "avproof")
// 				if err != nil {
// 					return err
// 				}
// 				for _, fichierCSV := range r.File {
// 					if fichierCSV.Name == "GetThis.csv" {
// 						archive_remplitArborescenceDepuisCSV(fichierCSV, &resultatArbo, &modeleArbo)
// 					}
// 				}
// 				return nil
// 			})
// 		}
// 	}
// 	// On enregistre l'arborescence que l'on vient d'extraire
// 	err = enregistrerArborescenceJson(&resultatArbo, filepath.Join(cheminProjet, "analyse", "arborescence.json"))
// 	return resultatArbo, err
// }

/*
* Fonction qui renvoie les caractéristiques des fichiers et dossiers contenus un dossier
@param cheminProjet : le chemin vers le projet ORC
@param cheminDossier : chemin vers le dossier duquel on veut les enfants
@return : une liste contenant les métadonnées des éléments contenus dans le dossier
// */
// func archive_RecupEnfantsArbo(cheminProjet string, cheminDossier []int) ([]MetaDonnees, error) {
// 	var fichiers []MetaDonnees = []MetaDonnees{}
// 	// Si l'arborescence nextraite du fichier .json, on l'extrait
// 	if len(cacheArbo.Enfants) == 0 {
// 		var err error
// 		cacheArbo, err = archive_GetArborescence(cheminProjet)
// 		if err != nil {
// 			return fichiers, err
// 		}
// 	}
// 	var vousetesici *Arborescence = &cacheArbo
// 	// On suit le chemin donné en paramètres pour se placer dans le dossier
// 	// duquel on veut le contenu
// 	for _, pas := range cheminDossier {
// 		if len(vousetesici.Enfants) < pas {
// 			return fichiers, errors.New("Le chemin de dossier spécifié est incohérent avec l'arborecsence")
// 		}
// 		vousetesici = &(*vousetesici).Enfants[pas]
// 	}
// 	// On parcourt les éléments de ce dossier
// 	for i := range vousetesici.Enfants {
// 		var legitimite []int = []int{0, 0}
// 		if (*vousetesici).Enfants[i].Legitimite == 0 {
// 			legitimite = []int{1, 0}
// 		} else if (*vousetesici).Enfants[i].Legitimite == 1 {
// 			legitimite = []int{0, 1}
// 		}
// 		var metadonnees MetaDonnees = MetaDonnees{
// 			Nom:             (*vousetesici).Enfants[i].Nom,
// 			ADesEnfants:     len((*vousetesici).Enfants[i].Enfants) != 0,
// 			EnfantsInconnus: legitimite[0],
// 			EnfantsSuspects: legitimite[1],
// 			Empreinte:       (*vousetesici).Enfants[i].EmpreinteMD5,
// 		}
// 		fichiers = append(fichiers, metadonnees)
// 	}
// 	return fichiers, nil
// }
