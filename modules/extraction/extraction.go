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

package extraction

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"aquarium/modules/extraction/bdd_sqlite"
	"aquarium/modules/extraction/csv"
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/journaux"
	"aquarium/modules/extraction/prefetch"
	"aquarium/modules/extraction/registre"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

type Extracteur interface {
	Extraction(string, io.Reader, string, config.ConfigExtraction, string) error
}

var progressionExtraction map[string]string = map[string]string{}
var listeExtractions map[string]ExtractionMachine = map[string]ExtractionMachine{}

var liste_extracteurs map[string]Extracteur = map[string]Extracteur{
	// "avs":        avlogs.AvLog{},
	"evtx":   evtx.Evtx{},
	"sqlite": bdd_sqlite.SQLite{},
	// "navigateur": navigateur.Navigateur{},
	// "werr":       werr.Werr{},
	"registre": registre.Registre{},
	"csv":      csv.Csv{},
	// "divers":     divers.Divers{},
	"prefetch": prefetch.Prefetch{},
	"journaux": journaux.Journaux{},
}

type ExtractionMachine struct {
	ConfigMachine    config.AquaConfigMachine
	ListeExtractions map[string]ParametresExtraction
}

type ParametresExtraction struct {
	Chemins         []config.ConfigChemin
	InfosExtraction config.InfosSommairesExtraction
}

/** Fonction qui renvoie des informations sur les extractions qui peuvent être effectuées
  * @param cheminProjet string : le chemin d'enregistrement de l'analyse aquarium
  * @return une liste de configurations d'extractions, et s'il y a lieu une erreur
**/
func ListeExtractionsHtml(cheminProjet string) (map[string]ExtractionMachine, error) {
	if len(listeExtractions) > 0 || progressionExtraction["chargement"] == "101" {
		return listeExtractions, nil
	}
	var configAnalyse config.AquaConfig
	configAnalyse, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return listeExtractions, err
	}
	for idMachine, machine := range configAnalyse.Machines {
		configMachine, err := config.GetConfigurationMachine(cheminProjet, idMachine, machine)
		if err != nil {
			return listeExtractions, err
		}
		var listeExtractionsMachine map[string]ParametresExtraction = map[string]ParametresExtraction{}
		for _, extraction := range configMachine.DetailsExtraction {
			// Récupérer la description de l’extraction
			infosExtraction, err := config.GetInfosSommairesExtraction(cheminProjet, extraction.Id)
			if err != nil {
				return listeExtractions, err
			}
			var parametresExtraction ParametresExtraction = ParametresExtraction{Chemins: extraction.Chemins, InfosExtraction: infosExtraction}
			listeExtractionsMachine[extraction.Id] = parametresExtraction
		}
		listeExtractions[idMachine] = ExtractionMachine{ConfigMachine: machine, ListeExtractions: listeExtractionsMachine}
	}
	return listeExtractions, nil
}

func LancerExtractions(cheminProjet string, ordreExtractions []map[string]string) error {
	for _, idExtractionMachine := range ordreExtractions {
		progressionExtraction["idMachine"] = idExtractionMachine["idMachine"]
		progressionExtraction["idExtraction"] = idExtractionMachine["idExtraction"]
		progressionExtraction["chargement"] = "0"
		err := Extraction(idExtractionMachine["idExtraction"], cheminProjet, idExtractionMachine["idMachine"])
		if err != nil {
			return err
		}
	}
	progressionExtraction["idMachine"] = ""
	progressionExtraction["idExtraction"] = ""
	return nil
}

/** Fonction qui exécute une extraction définie dans le fichier de configuration
  * @param idExtraction : l'identifiant de l'extraction
  * @cheminProjet string : le chemin d'enregistrement du projet
  * @return : une erreur s'il y a lieu
**/
func Extraction(idExtraction string, cheminProjet string, idMachine string) error {
	var probleme error
	// On récupère la configuration de l’extraction
	configExtr, err := config.GetConfigExtraction(cheminProjet, idExtraction)
	if err != nil {
		return err
	}
	// On vérifie que la table est bien créée
	err = creerTableExtraction(cheminProjet, configExtr)
	if err != nil {
		return err
	}
	// On liste les fichiers concernés par cette extraction
	listeFichiersAExtraire, probleme := config.ListeFichiersExtraction(listeExtractions[idMachine].ListeExtractions[idExtraction].Chemins, cheminProjet, idMachine, false)
	// On met la progression à 0 (début de l'extraction)
	progressionExtraction["chargement"] = "0"
	// On crée la table qui sera utilisée par l'extracteur
	err = creerTableExtraction(cheminProjet, configExtr)
	if err != nil {
		return err
	}
	// On compte le nombre de fichiers à extraire
	var nbFichiers int = 0
	for _, dossierAExtraire := range listeFichiersAExtraire {
		nbFichiers += len(dossierAExtraire.Elements)
	}
	var i int = 0
	// On boucle sur les fichiers à extraire
	for _, dossierAExtraire := range listeFichiersAExtraire {
		if dossierAExtraire.Est7Z {
			probleme = extrationAchive7z(cheminProjet, dossierAExtraire, idMachine, configExtr, &i, nbFichiers)
		} else {
			probleme = extractionDossier(cheminProjet, dossierAExtraire, idMachine, configExtr, &i, nbFichiers)
		}
	}
	progressionExtraction["chargement"] = "101"
	return probleme
}

/** Fonction qui crée toutes les tables de l'analyse
  * @param cheminProjet string : le chemin d'enregistrement du projet
**/
func CreationBaseAnalyse(cheminProjet string) error {
	// On cherche la configuration du projet
	aquaConfig, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return err
	}
	for idMachine, paramMachine := range aquaConfig.Machines {
		configMachine, err := config.GetConfigurationMachine(cheminProjet, idMachine, paramMachine)
		if err != nil {
			return err
		}
		for _, extraction := range configMachine.DetailsExtraction {
			configExtraction, err := config.GetConfigExtraction(cheminProjet, extraction.Id)
			if err != nil {
				return err
			}
			creerTableExtraction(cheminProjet, configExtraction)
		}
	}
	err = creerTableChronologie(cheminProjet)
	return err
}

/** Fonction qui renvoie le pourcentage de chargement de l'extraction
 ** @param cheminProjet string : le chemin de l'analyse aquarium
 ** @param idExtraction string : l'identifiant de l'extraction
**/
func ProgressionExtraction(cheminProjet string) map[string]string {
	return progressionExtraction
}

/** Fonction qui lance l'extraction de la table chronologie, qui contient un
  * résumé de tous les évènements
  * @param cheminprojet : le chemin d'enregistrament de l'analyse aquarium
  * @return : une erreur s'il y a lieu
**/
func ExtraireTableChronologie(cheminProjet string) error {
	aquaConfig, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return err
	}
	for idMachine, machine := range aquaConfig.Machines {
		configMachine, err := config.GetConfigurationMachine(cheminProjet, idMachine, machine)
		var listeRequetesChronologie []string = []string{}
		// On fait une liste des requêtes SQL à exécuter
		for _, extraction := range configMachine.DetailsExtraction {
			configExtraction, err := config.GetConfigExtraction(cheminProjet, extraction.Id)
			if err != nil {
				return err
			}
			if configExtraction.SQLChronologie != "" {
				listeRequetesChronologie = append(listeRequetesChronologie, configExtraction.SQLChronologie)
			}
		}
		// On liste les colonnes de la table chronologie
		var listeColonnesChrolonogie []string = []string{}
		for _, colonne := range configMachine.Chronologie.Colonnes {
			listeColonnesChrolonogie = append(listeColonnesChrolonogie, colonne.Nom)
		}
		var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
		err = base.RemplirTableDepuisRequetes(configMachine.Chronologie.Nom, listeColonnesChrolonogie, listeRequetesChronologie, true, "horodatage")
		if err != nil {
			return err
		}
	}
	return nil
}

/** Fonction qui renvoie le contenu de la tables "chronologie", contenant un résumé de l'ensemble des
  * évènemets
  * @param cheminprojet string : le chemin d'enregistrement de l'analyse aquarium
  * @param debut int : l'index à partir duquel on veut récupérer les valeurs
  * @param taill int : le nombre de valeurs que l'on veut récupérer
  * @return une liste des lignes de la table
**/ /*
func ValeursTableChronologie(cheminProjet string, debut int, taille int) []map[string]interface{} {
	var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	configAnalyse, err := config.GetConfigurationProjet(cheminProjet)
	if err != nil {
		return []map[string]interface{}{}
	}
	var listeColonnesChrolonogie []string = []string{}
	for _, colonne := range configAnalyse.Chronologie.Colonnes {
		listeColonnesChrolonogie = append(listeColonnesChrolonogie, colonne.Nom)
	}
	return abase.RecupererValeursTable("chronologie", listeColonnesChrolonogie, debut, taille)
}*/

/** --------------------- FONCTIONS À USAGE INTERNE --------------------- **/

func creerTableExtraction(cheminProjet string, extraction config.ConfigExtraction) error {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	for _, table := range extraction.Table {
		var listeColonnes map[string]string = map[string]string{}
		for _, colonne := range table.Colonnes {
			listeColonnes[colonne.Nom] = colonne.Type
		}
		err := base.CreateTableIfNotExist2(table.Nom, listeColonnes, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func extrationAchive7z(cheminProjet string, configArchive config.DossierAExtraire, idMachine string, configExtraction config.ConfigExtraction, i *int, total int) error {
	// On commence par ouvrir l'archive
	archive, err := sevenzip.OpenReaderWithPassword(configArchive.Chemin, "avproof")
	if err != nil {
		return err
	}
	for _, numFichier := range configArchive.Elements {
		fichier, err := archive.File[numFichier].Open()
		if err != nil {
			continue
		}
		defer fichier.Close()
		var source string = strings.Replace(filepath.Join(configArchive.Chemin, archive.File[numFichier].Name), cheminProjet, "", 1)
		extracteur, ok := liste_extracteurs[configExtraction.Extracteur]
		if !ok {
			return errors.New("L’extracteur « " + configExtraction.Extracteur + " » n'existe pas. Vérifiez le fichier de configuration.")
		}
		extracteur.Extraction(cheminProjet, fichier, source, configExtraction, idMachine)
		// On change la progression du chargement
		*i++
		progressionExtraction["chargement"] = fmt.Sprintf("%f", (float32(*i) / float32(total) * 100))
		fichier.Close()
	}
	return nil
}

func extractionDossier(cheminProjet string, configDossier config.DossierAExtraire, idMachine string, configExtraction config.ConfigExtraction, i *int, total int) error {
	// On commence par lire le dossier
	contenuDossier, err := os.ReadDir(configDossier.Chemin)
	if err != nil {
		return err
	}
	for _, numFichier := range configDossier.Elements {
		var cheminFichier string = filepath.Join(configDossier.Chemin, contenuDossier[numFichier].Name())
		fichier, err := os.Open(cheminFichier)
		if err != nil {
			continue
		}
		defer fichier.Close()
		extracteur, ok := liste_extracteurs[configExtraction.Extracteur]
		if !ok {
			return errors.New("L’extracteur « " + configExtraction.Extracteur + " » n'existe pas. Vérifiez le fichier de configuration.")
		}
		extracteur.Extraction(cheminProjet, fichier, strings.Replace(cheminFichier, cheminProjet, "", 1), configExtraction, idMachine)
		fichier.Close()
	}
	return nil
}

/*
*
Fonction qui crée les tables de chronologie en fonction de la configuration des machines de l’analyse
S’il y a deux tables de chronologie du même nom non identiques, une table contenant les colonnes des deux
tables sera créée.
@param cheminProjet : le chemin du dossier d'analyse
*/
func creerTableChronologie(cheminProjet string) error {
	// Récupération du fichier de configuration
	aquaconfig, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return err
	}
	// Établissement de la liste des colonnes
	var listeColonnes map[string]map[string]string = map[string]map[string]string{}
	for idMachine, configMachine := range aquaconfig.Machines {
		config, err := config.GetConfigurationMachine(cheminProjet, idMachine, configMachine)
		_, tableExiste := listeColonnes[config.Chronologie.Nom]
		if !tableExiste {
			listeColonnes[config.Chronologie.Nom] = map[string]string{}
		}
		if err != nil {
			return err
		}
		for _, colonne := range config.Chronologie.Colonnes {
			typeColonne, existe := listeColonnes[config.Chronologie.Nom][colonne.Nom]
			if !existe {
				listeColonnes[config.Chronologie.Nom][colonne.Nom] = colonne.Type
			}
			if typeColonne != colonne.Type {
				listeColonnes[config.Chronologie.Nom][colonne.Nom] = "TEXT"
			}
		}
	}
	// Création de la base de donnees
	var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	var probleme error = nil
	for nomTable, colonnesTable := range listeColonnes {
		err = abase.CreateTableIfNotExist2(nomTable, colonnesTable, false)
		if err != nil {
			probleme = err
		}
	}
	return probleme
}
