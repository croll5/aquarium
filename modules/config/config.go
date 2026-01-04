package config

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
)

var DOSSIER_FICHIERS_A_ANALYSER = "fichiers"
var AQUA_MACHINE = "aqua_machine"
var DOSSIER_EXTRACTIONS = "extractions"
var EXTENSION_XML = ".xml"

type ConfigColonneBDD struct {
	Nom     string `xml:",chardata"`
	Type    string `xml:"type,attr"`
	Contenu string `xml:"contenu,attr"`
}

type ConfigTableBDD struct {
	Nom       string             `xml:"nom,attr"`
	Condition string             `xml:"condition,attr"`
	Colonnes  []ConfigColonneBDD `xml:"colonne"`
}

type ConfigChemin struct {
	Dossiers []string `xml:"dossier"`
	Archive  string   `xml:"archive"`
	Fichier  string   `xml:"fichier"`
}

type ConfigExtraction struct {
	Id                 string           `xml:"id"`
	Nom                string           `xml:"nom"`
	Description        string           `xml:"description"`
	Extracteur         string           `xml:"extracteur"`
	Table              []ConfigTableBDD `xml:"table"`
	SQLChronologie     string           `xml:"sql_chronologie"`
	Progression        float32
	AnnulationDemandee bool
	ConfigComplement   []ComplementConfigXML `xml:"complements>parametre"`
	Complement         map[string]string
}

type InfosSommairesExtraction struct {
	Description string `xml:"description"`
	Nom         string `xml:"nom"`
}

type DetailsConfigExtraction struct {
	Id      string         `xml:"id,attr"`
	Chemins []ConfigChemin `xml:"chemin"`
}

type ConfigurationXML struct {
	Chronologie       ConfigTableBDD            `xml:"chronologie"`
	DetailsExtraction []DetailsConfigExtraction `xml:"extraction"`
}

type DossierAExtraire struct {
	Chemin   string
	Est7Z    bool
	Elements []int
}

type ComplementConfigXML struct {
	Cle    string `xml:"cle,attr"`
	Valeur string `xml:",innerxml"`
}

func GetConfigurationMachine(cheminProjet string, idMachine string, aquaConfigMachine AquaConfigMachine) (ConfigurationXML, error) {
	// On commence par récupérer le chemin du fichier de configuration
	var donneesConfig ConfigurationXML
	var cheminFichierConfigPrincipal string
	cheminFichierConfigPrincipal, err := cheminFichierConfig(cheminProjet, aquaConfigMachine.Config)
	if err != nil {
		return donneesConfig, err
	}
	fichierConf, err := os.Open(cheminFichierConfigPrincipal)
	if err != nil {
		return donneesConfig, err
	}
	bytesConfig, err := io.ReadAll(fichierConf)
	if err != nil {
		return donneesConfig, err
	}
	err = xml.Unmarshal(bytesConfig, &donneesConfig)
	// On supprime les extractions qui n’ont pas les fichiers nécessiares
	var listeExtractions []DetailsConfigExtraction = []DetailsConfigExtraction{}
	for _, extraction := range donneesConfig.DetailsExtraction {
		listeDossier, err := ListeFichiersExtraction(extraction.Chemins, cheminProjet, idMachine, false)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(listeDossier) > 0 {
			listeExtractions = append(listeExtractions, extraction)
		}
	}
	donneesConfig.DetailsExtraction = listeExtractions
	return donneesConfig, err
}

func GetConfigExtraction(cheminProjet string, nomFichierConfig string) (ConfigExtraction, error) {
	var donneesConfig ConfigExtraction
	var cheminConfig string
	cheminConfig, err := cheminFichierConfig(cheminProjet, filepath.Join(DOSSIER_EXTRACTIONS, nomFichierConfig)+EXTENSION_XML)
	if err != nil {
		return donneesConfig, err
	}
	fichierConfig, err := os.Open(cheminConfig)
	if err != nil {
		return donneesConfig, err
	}
	bytesConfig, err := io.ReadAll(fichierConfig)
	if err != nil {
		return donneesConfig, err
	}
	err = xml.Unmarshal(bytesConfig, &donneesConfig)
	return donneesConfig, err
}

func GetInfosSommairesExtraction(cheminProjet string, nomFichierConfig string) (InfosSommairesExtraction, error) {
	var infosExtraction InfosSommairesExtraction
	var cheminConfig string
	cheminConfig, err := cheminFichierConfig(cheminProjet, filepath.Join(DOSSIER_EXTRACTIONS, nomFichierConfig)+EXTENSION_XML)
	if err != nil {
		return infosExtraction, err
	}
	fichierConfig, err := os.Open(cheminConfig)
	if err != nil {
		return infosExtraction, err
	}
	bytesConfig, err := io.ReadAll(fichierConfig)
	if err != nil {
		return infosExtraction, err
	}
	err = xml.Unmarshal(bytesConfig, &infosExtraction)
	return infosExtraction, err
}

func lireConfigExtraction(cheminExtraction string) (ConfigExtraction, error) {
	var configExtraction ConfigExtraction
	// On ouvre le ficher de configuration
	fichierConfig, err := os.Open(cheminExtraction)
	if err != nil {
		return configExtraction, err
	}
	bytesConfig, err := io.ReadAll(fichierConfig)
	if err != nil {
		return configExtraction, err
	}
	err = xml.Unmarshal(bytesConfig, &configExtraction)
	return configExtraction, err
}

func ListeFichiersExtraction(chemins []ConfigChemin, cheminProjet string, dossierMachine string, fichierUnique bool) ([]DossierAExtraire, error) {
	var resultat []DossierAExtraire = []DossierAExtraire{}
	var probleme error
	for _, configChemin := range chemins {
		listeDossiers, err := listeCheminsDossiersRec(configChemin.Dossiers, filepath.Join(cheminProjet, DOSSIER_FICHIERS_A_ANALYSER, dossierMachine), fichierUnique)
		if err != nil {
			return resultat, err
		}
		for _, dossier := range listeDossiers {
			var err error
			var correspondances []DossierAExtraire
			if configChemin.Archive != "" {
				correspondances, err = listeCheminsArchives(dossier, configChemin.Archive, configChemin.Fichier, fichierUnique)
			} else {
				correspondances, err = listeCorrespondancesDansDossier(dossier, configChemin.Fichier)
			}
			if err != nil {
				probleme = err
				continue
			}
			resultat = append(resultat, correspondances...)
			if fichierUnique && len(correspondances) > 0 {
				return resultat, nil
			}
		}
	}

	return resultat, probleme
}

/**                           FONCTIONS LOCALES                           **/

func cheminFichierConfig(cheminProjet string, nomFichierConfig string) (string, error) {
	// On commence par regarder si le fichier est présent dans le dossier de l'analyse
	var cheminLocal string = filepath.Join(cheminProjet, "config", nomFichierConfig)
	_, err := os.Stat(cheminLocal)
	if err == nil {
		return cheminLocal, nil
	}
	emplacementExecutable, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(emplacementExecutable), "config", nomFichierConfig), nil
}

/** Fonction qui renvoie la liste des chemins qui parcourent les dossiers donnés en argument
**/
func listeCheminsDossiersRec(cheminCible []string, cheminActuel string, fichierUnique bool) ([]string, error) {
	if len(cheminCible) == 0 {
		return []string{cheminActuel}, nil
	} else {
		var resultat []string = []string{}
		var probleme error
		// On ouvre le chemin actuel pour voir les dossiers qui y sont présents
		dossiers, err := os.ReadDir(cheminActuel)
		if err != nil {
			return resultat, err
		}
		// On les parcours pour voir s'ils correspondent à la cible
		for _, dossier := range dossiers {
			correspond, err := filepath.Match(cheminCible[0], dossier.Name())
			if err != nil {
				probleme = err
				continue
			}
			if correspond && dossier.IsDir() {
				nvChemins, err := listeCheminsDossiersRec(cheminCible[1:], filepath.Join(cheminActuel, dossier.Name()), fichierUnique)
				if err != nil {
					probleme = err
					continue
				}
				resultat = append(resultat, nvChemins...)
				if fichierUnique {
					return resultat, nil
				}
			}
		}
		return resultat, probleme
	}
}

func listeCheminsArchives(cheminActuel string, nomArchive string, fichierCible string, fichierUnique bool) ([]DossierAExtraire, error) {
	var resultat []DossierAExtraire = []DossierAExtraire{}
	// On parcourt le dossier dans lequel se trouve l'archive à la recherche de celle-ci
	archives, err := os.ReadDir(cheminActuel)
	if err != nil {
		return resultat, err
	}
	for _, archive := range archives {
		correspond, err := filepath.Match(nomArchive, archive.Name())
		if err != nil {
			return resultat, err
		}
		if correspond {
			correspondances, err := listeCorrespondancesDansArchive(filepath.Join(cheminActuel, archive.Name()), fichierCible, "avproof")
			if err != nil {
				return resultat, err
			}
			resultat = append(resultat, correspondances)
			if fichierUnique {
				return resultat, nil
			}
		}
	}
	return resultat, nil
}

func listeCorrespondancesDansArchive(cheminArchive string, cheminCible string, mdp string) (DossierAExtraire, error) {
	var listePositions []int = []int{}
	var probleme error
	// On ouvre l'archive et on la parcourt
	archive, err := sevenzip.OpenReaderWithPassword(cheminArchive, mdp)
	if err != nil {
		return DossierAExtraire{}, err
	}
	for num, fichier := range archive.File {
		correspond, err := filepath.Match(cheminCible, fichier.Name)
		if err != nil {
			probleme = err
			continue
		}
		if correspond {
			listePositions = append(listePositions, num)
		}
	}
	return DossierAExtraire{Chemin: cheminArchive, Est7Z: true, Elements: listePositions}, probleme
}

func listeCorrespondancesDansDossier(cheminDossier string, nomCible string) ([]DossierAExtraire, error) {
	var numFichiersOk []int = []int{}
	var probleme error
	// On ouvre le dossier à explorer
	fichiers, err := os.ReadDir(cheminDossier)
	if err != nil {
		return []DossierAExtraire{}, err
	}
	for numFichier, fichier := range fichiers {
		correspond, err := filepath.Match(nomCible, fichier.Name())
		if err != nil {
			probleme = err
			continue
		}
		if correspond {
			numFichiersOk = append(numFichiersOk, numFichier)
		}
	}
	return []DossierAExtraire{DossierAExtraire{Chemin: cheminDossier, Est7Z: false, Elements: numFichiersOk}}, probleme
}

func (cftable ConfigTableBDD) GetNomsColonnes() []string {
	var res []string = []string{}
	for _, colonne := range cftable.Colonnes {
		res = append(res, colonne.Nom)
	}
	return res
}

func GetListeConfigurationsDisponibles() ([]string, error) {
	// Récupération de l’emplacement de l’exécutable
	emplacementExecutable, err := os.Executable()
	if err != nil {
		return []string{}, err
	}
	dossierConfig := filepath.Join(filepath.Dir(emplacementExecutable), "config")
	// Liste des fichiers du dossier
	fichiers, err := os.ReadDir(dossierConfig)
	if err != nil {
		return []string{}, err
	}
	// Énumération des ficheirs xml
	var listeConfigs []string = []string{}
	for _, fichier := range fichiers {
		ok, _ := filepath.Match("*.xml", fichier.Name())
		if ok {
			listeConfigs = append(listeConfigs, fichier.Name())
		}
	}
	return listeConfigs, nil
}

func ajouterColonneMachineDansTables(confTable []ConfigTableBDD) []ConfigTableBDD {
	for i := range confTable {
		confTable[i].Colonnes = append(confTable[i].Colonnes, ConfigColonneBDD{Nom: AQUA_MACHINE, Contenu: AQUA_MACHINE, Type: "VARCHAR(20)"})
	}
	return confTable
}
