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
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
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
		return errors.WithStack(err)
	}
	// On lit la configuration de la machine
	configMachine, err := config.GetConfigurationMachine(cheminProjet, idMachine, configAnalyse.Machines[idMachine])
	if err != nil {
		return errors.WithStack(err)
	}
	// On commence par lister les fichiers
	var abd *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	requeteSQL := strings.ReplaceAll(configMachine.Arborescence.RequeteCreation, AQUA_MACHINE, idMachine)
	resultatRequete, err := abd.ExecuterRequeteSQL(requeteSQL)
	if err != nil {
		return errors.WithStack(err)
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
	}
	// Enregistrement de l’analyse
	err = os.MkdirAll(filepath.Join(cheminProjet, config.DOSSIER_ANALYSE, DOSSIER_ARBO), os.ModeAppend)
	if err != nil {
		return errors.WithStack(err)
	}
	err = enregistrerArborescenceJson(&cache, filepath.Join(cheminProjet, config.DOSSIER_ANALYSE, DOSSIER_ARBO, idMachine)+".json")
	if err != nil {
		return errors.WithStack(err)
	}
	cacheArbo = cache
	nomArboEnCache = idMachine
	enCoursDextraction = false
	return err
}

func RecupEnfantsArbo(cheminProjet string, cheminDossier []string, idMachine string) ([]MetaDonnees, error) {
	// On commence par récupérer la configuration de la machine
	aquaConfigMachine, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return []MetaDonnees{}, errors.WithStack(err)
	}
	configAnalyse, err := config.GetConfigurationMachine(cheminProjet, idMachine, aquaConfigMachine.Machines[idMachine])
	if err != nil {
		return []MetaDonnees{}, errors.WithStack(err)
	}
	// Si l'arborescence en cache ne correspond pas, on l’extrait à nouveau
	if idMachine != nomArboEnCache {
		// On essaie d’ouvrir le fichier de l’arborescence
		contenuFichierArbo, err := os.ReadFile(filepath.Join(cheminProjet, config.DOSSIER_ANALYSE, DOSSIER_ARBO, idMachine) + ".json")
		if err != nil {
			err = ExtraireArborescence(cheminProjet, "", idMachine)
			if err != nil {
				return []MetaDonnees{}, errors.WithStack(err)
			}
			return getContenuDossier(cheminProjet, cheminDossier, configAnalyse.Arborescence)
		}
		// Si l’arborescence existe déjà, on la charge dans le cache
		cacheArbo = Arborescence{}
		err = json.Unmarshal(contenuFichierArbo, &cacheArbo)
		if err != nil {
			return []MetaDonnees{}, errors.WithStack(err)
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
		return errors.WithStack(errors.New("Problème dans la requête de sélection de l’arborescence"))
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
		return errors.WithStack(errors.New("Type de donnée non reconnu"))
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
			return []MetaDonnees{}, errors.WithStack(os.ErrNotExist)
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
		resultat, err := adb.SelectFrom(configArborescence.RequeteMetadonnees, fichier.Source)
		if err != nil {
			return []MetaDonnees{}, errors.WithStack(err)
		}
		if len(resultat) < 1 {
			return []MetaDonnees{}, errors.WithStack(errors.Errorf("Résultat vide"))
		}
		switch resultat[0]["Nom"].(type) {
		case string:
			metaDonnees[i] = MetaDonnees{Nom: resultat[0]["Nom"].(string), IdSource: fichier.Source, IdCopie: fichier.Copie}
		default:
			return []MetaDonnees{}, errors.WithStack(errors.Errorf("Résultat vide"))
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
		return errors.WithStack(err)
	}
	fichier, err := os.Create(chemin)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fichier.Close()
	_, err = fichier.Write(donneesArbo)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}
