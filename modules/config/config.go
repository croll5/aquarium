package config

import (
	"aquarium/modules/aquabase"
	"embed"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
	"github.com/pkg/errors"
)

const (
	DOSSIER_FICHIERS_A_ANALYSER = "fichiers"
	DOSSIER_ANALYSE             = "analyse"
	AQUA_MACHINE                = "aqua_machine"
	DOSSIER_CONFIG_EXTRACTIONS  = "config_extractions"
	DOSSIER_CONFIG_MACHINES     = "config_machines"
	EXTENSION_XML               = ".xml"
	DOSSIER_CONFIG              = "config"
	DOSSIER_RESSOURCES          = "ressources"
	DOSSIER_REGLES_DETECTIONS   = "regles_detection"
)

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
	Id               string                `xml:"id"`
	Nom              string                `xml:"nom"`
	Description      string                `xml:"description"`
	Extracteur       string                `xml:"extracteur"`
	Table            []ConfigTableBDD      `xml:"table"`
	SQLChronologie   string                `xml:"sql_chronologie"`
	ConfigComplement []ComplementConfigXML `xml:"complements>parametre"`
	Complement       map[string]string
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
	Arborescence      ConfigArborescence        `xml:"arborescence"`
}

type ConfigArborescence struct {
	RequeteCreation    string `xml:"requete_creation"`
	RequeteMetadonnees string `xml:"requete_metadonnees"`
	SeparateurDossier  string `xml:"separateur_dossiers"`
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

var cacheConfigMachines map[string]ConfigurationXML = map[string]ConfigurationXML{}
var cacheConfigExtractions map[string]ConfigExtraction = map[string]ConfigExtraction{}

func GetConfigurationMachine(cheminProjet string, idMachine string, aquaConfigMachine AquaConfigMachine) (ConfigurationXML, error) {
	cacheConfig, existe := cacheConfigMachines[idMachine]
	if existe {
		return cacheConfig, nil
	}
	// On commence par récupérer le chemin du fichier de configuration
	var donneesConfig ConfigurationXML
	var cheminFichierConfigPrincipal string
	cheminFichierConfigPrincipal, err := cheminFichierConfig(cheminProjet, filepath.Join(DOSSIER_CONFIG_MACHINES, aquaConfigMachine.Config))
	if err != nil {
		return donneesConfig, errors.WithStack(err)
	}
	fichierConf, err := os.Open(cheminFichierConfigPrincipal)
	if err != nil {
		return donneesConfig, errors.WithStack(err)
	}
	bytesConfig, err := io.ReadAll(fichierConf)
	if err != nil {
		return donneesConfig, errors.WithStack(err)
	}
	err = xml.Unmarshal(bytesConfig, &donneesConfig)
	// On supprime les extractions qui n’ont pas les fichiers nécessiares
	var listeExtractions []DetailsConfigExtraction = []DetailsConfigExtraction{}
	var problemeRencontre error
	for _, extraction := range donneesConfig.DetailsExtraction {
		listeDossier, err := ListeFichiersExtraction(extraction.Chemins, cheminProjet, idMachine, false)
		if err != nil && len(listeDossier) == 0 {
			problemeRencontre = errors.WithStack(err)
		}
		var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
		infosExtraction, err := GetConfigExtraction(cheminProjet, extraction.Id)
		if err != nil {
			problemeRencontre = errors.WithStack(err)
			continue
		}
		for _, tableExtraction := range infosExtraction.Table {
			if abase.EstTableVide(tableExtraction.Nom) {
				if len(listeDossier) > 0 {
					listeExtractions = append(listeExtractions, extraction)
				}
				break
			}
		}
	}
	donneesConfig.DetailsExtraction = listeExtractions
	cacheConfigMachines[idMachine] = donneesConfig
	return donneesConfig, problemeRencontre
}

func EnregistrerConfigMachine(cheminProjet string, nomConfig string, reutilisable bool, configuration ConfigurationXML) error {
	fmt.Println(configuration)
	return nil
}

func GetConfigExtraction(cheminProjet string, nomFichierConfig string) (ConfigExtraction, error) {
	cacheConfig, existe := cacheConfigExtractions[nomFichierConfig]
	if existe {
		return cacheConfig, nil
	}
	var donneesConfig ConfigExtraction
	var cheminConfig string
	cheminConfig, err := cheminFichierConfig(cheminProjet, filepath.Join(DOSSIER_CONFIG_EXTRACTIONS, nomFichierConfig)+EXTENSION_XML)
	if err != nil {
		return donneesConfig, errors.WithStack(err)
	}
	fichierConfig, err := os.Open(cheminConfig)
	if err != nil {
		return donneesConfig, errors.WithStack(err)
	}
	bytesConfig, err := io.ReadAll(fichierConfig)
	if err != nil {
		return donneesConfig, errors.WithStack(err)
	}
	err = xml.Unmarshal(bytesConfig, &donneesConfig)
	if err != nil {
		return donneesConfig, errors.WithStack(err)
	}
	// On récupère les données complémentaires
	var listeComplements map[string]string = map[string]string{}
	for _, complement := range donneesConfig.ConfigComplement {
		listeComplements[complement.Cle] = complement.Valeur
	}
	donneesConfig.Complement = listeComplements
	cacheConfigExtractions[nomFichierConfig] = donneesConfig
	// On ajoute le nom de la machine à la liste des colonnes
	for i := range donneesConfig.Table {
		donneesConfig.Table[i].Colonnes = append(donneesConfig.Table[i].Colonnes, ConfigColonneBDD{Nom: AQUA_MACHINE, Contenu: AQUA_MACHINE, Type: "TEXT"})
	}
	return donneesConfig, err
}

func GetInfosSommairesExtraction(cheminProjet string, nomFichierConfig string) (InfosSommairesExtraction, error) {
	var infosExtraction InfosSommairesExtraction
	var cheminConfig string
	cheminConfig, err := cheminFichierConfig(cheminProjet, filepath.Join(DOSSIER_CONFIG_EXTRACTIONS, nomFichierConfig)+EXTENSION_XML)
	if err != nil {
		return infosExtraction, errors.WithStack(err)
	}
	fichierConfig, err := os.Open(cheminConfig)
	if err != nil {
		return infosExtraction, errors.WithStack(err)
	}
	bytesConfig, err := io.ReadAll(fichierConfig)
	if err != nil {
		return infosExtraction, errors.WithStack(err)
	}
	err = xml.Unmarshal(bytesConfig, &infosExtraction)
	if err != nil {
		return infosExtraction, errors.WithStack(err)
	}
	return infosExtraction, err
}

func lireConfigExtraction(cheminExtraction string) (ConfigExtraction, error) {
	var configExtraction ConfigExtraction
	// On ouvre le ficher de configuration
	fichierConfig, err := os.Open(cheminExtraction)
	if err != nil {
		return configExtraction, errors.WithStack(err)
	}
	bytesConfig, err := io.ReadAll(fichierConfig)
	if err != nil {
		return configExtraction, errors.WithStack(err)
	}
	err = xml.Unmarshal(bytesConfig, &configExtraction)
	if err != nil {
		return configExtraction, errors.WithStack(err)
	}
	return configExtraction, err
}

func ListeFichiersExtraction(chemins []ConfigChemin, cheminProjet string, dossierMachine string, fichierUnique bool) ([]DossierAExtraire, error) {
	var resultat []DossierAExtraire = []DossierAExtraire{}
	var probleme error
	for _, configChemin := range chemins {
		listeDossiers, err := listeCheminsDossiersRec(configChemin.Dossiers, filepath.Join(cheminProjet, DOSSIER_FICHIERS_A_ANALYSER, dossierMachine), fichierUnique)
		if err != nil {
			return resultat, errors.WithStack(err)
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
				probleme = errors.WithStack(err)
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
	var cheminLocal string = filepath.Join(cheminProjet, DOSSIER_CONFIG, nomFichierConfig)
	_, err := os.Stat(cheminLocal)
	if err == nil {
		return cheminLocal, nil
	}
	emplacementExecutable, err := os.Executable()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(filepath.Dir(emplacementExecutable), DOSSIER_CONFIG, nomFichierConfig), nil
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
			return resultat, errors.WithStack(err)
		}
		// On les parcours pour voir s'ils correspondent à la cible
		for _, dossier := range dossiers {
			correspond, err := filepath.Match(cheminCible[0], dossier.Name())
			if err != nil {
				probleme = errors.WithStack(err)
				continue
			}
			if correspond && dossier.IsDir() {
				nvChemins, err := listeCheminsDossiersRec(cheminCible[1:], filepath.Join(cheminActuel, dossier.Name()), fichierUnique)
				if err != nil {
					probleme = errors.WithStack(err)
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
		return resultat, errors.WithStack(err)
	}
	var problemeRencontre error
	for _, archive := range archives {
		correspond, err := filepath.Match(nomArchive, archive.Name())
		if err != nil {
			problemeRencontre = errors.WithStack(err)
		}
		if correspond {
			correspondances, err := listeCorrespondancesDansArchive(filepath.Join(cheminActuel, archive.Name()), fichierCible, "avproof")
			if err != nil {
				problemeRencontre = errors.WithStack(err)
			}
			resultat = append(resultat, correspondances)
			if fichierUnique {
				return resultat, nil
			}
		}
	}
	return resultat, problemeRencontre
}

func listeCorrespondancesDansArchive(cheminArchive string, cheminCible string, mdp string) (DossierAExtraire, error) {
	var listePositions []int = []int{}
	var probleme error
	// On ouvre l'archive et on la parcourt
	archive, err := sevenzip.OpenReaderWithPassword(cheminArchive, mdp)
	if err != nil {
		return DossierAExtraire{}, errors.WithStack(err)
	}
	for num, fichier := range archive.File {
		correspond, err := filepath.Match(cheminCible, fichier.Name)
		if err != nil {
			probleme = errors.WithStack(err)
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
		return []DossierAExtraire{}, errors.WithStack(err)
	}
	for numFichier, fichier := range fichiers {
		correspond, err := filepath.Match(nomCible, fichier.Name())
		if err != nil {
			probleme = errors.WithStack(err)
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

func GetListeConfigurationsDisponibles(dossierExtrations bool) ([]string, error) {
	// Récupération de l’emplacement de l’exécutable
	emplacementExecutable, err := os.Executable()
	if err != nil {
		return []string{}, errors.WithStack(err)
	}
	dossierConfig := filepath.Join(filepath.Dir(emplacementExecutable), DOSSIER_CONFIG)
	if dossierExtrations {
		dossierConfig = filepath.Join(dossierConfig, DOSSIER_CONFIG_EXTRACTIONS)
	} else {
		dossierConfig = filepath.Join(dossierConfig, DOSSIER_CONFIG_MACHINES)
	}
	// Liste des fichiers du dossier
	fichiers, err := os.ReadDir(dossierConfig)
	if err != nil {
		err := initialiserDossierConfigDepuisEmbarque(filepath.Dir(emplacementExecutable))
		if err != nil {
			return []string{}, errors.WithStack(err)
		}
		fichiers, err = os.ReadDir(dossierConfig)
		if err != nil {
			return []string{}, errors.WithStack(err)
		}
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

func VerifierPrerequisDemarrage() ([]string, error) {
	emplacementExecutable, err := os.Executable()
	if err != nil {
		return []string{}, errors.WithStack(err)
	}
	base := filepath.Dir(emplacementExecutable)

	type prerequis struct {
		cheminRelatif string
		extension     string
	}

	prerequisAttendus := []prerequis{
		{cheminRelatif: "./config", extension: ""},
		{cheminRelatif: "./config/config_extractions", extension: ".xml"},
		{cheminRelatif: "./config/config_machines", extension: ".xml"},
		{cheminRelatif: "./ressources", extension: ""},
		{cheminRelatif: "./ressources/regles_detection", extension: ".json"},
	}

	manquants := []string{}
	for _, attendu := range prerequisAttendus {
		cheminComplet := filepath.Join(base, attendu.cheminRelatif)
		infos, err := os.Stat(cheminComplet)
		if err != nil || !infos.IsDir() {
			manquants = append(manquants, attendu.cheminRelatif)
			continue
		}
		if attendu.extension == "" {
			continue
		}
		fichiers, err := os.ReadDir(cheminComplet)
		if err != nil {
			manquants = append(manquants, attendu.cheminRelatif+" (aucun "+attendu.extension+")")
			continue
		}
		trouveFichier := false
		for _, fichier := range fichiers {
			if fichier.IsDir() {
				continue
			}
			ok, _ := filepath.Match("*"+attendu.extension, fichier.Name())
			if ok {
				trouveFichier = true
				break
			}
		}
		if !trouveFichier {
			manquants = append(manquants, attendu.cheminRelatif+" (aucun "+attendu.extension+")")
		}
	}

	return manquants, nil
}

//go:embed config_embarquee/*
var configEmbarquee embed.FS

func initialiserDossierConfigDepuisEmbarque(emplacementExecutable string) error {
	for _, nomConfig := range []string{DOSSIER_CONFIG_EXTRACTIONS, DOSSIER_CONFIG_MACHINES} {
		fichiersConfigExtraction, err := configEmbarquee.ReadDir("config_embarquee/" + nomConfig)
		if err != nil {
			return errors.WithStack(err)
		}
		cheminDossierExtractions := filepath.Join(emplacementExecutable, DOSSIER_CONFIG, nomConfig)
		os.MkdirAll(cheminDossierExtractions, 0o755)
		for _, config := range fichiersConfigExtraction {
			contenuFichierExtraction, err := configEmbarquee.ReadFile("config_embarquee/" + nomConfig + "/" + config.Name())
			if err != nil {
				return errors.WithStack(err)
			}
			os.WriteFile(filepath.Join(cheminDossierExtractions, config.Name()), contenuFichierExtraction, 0o755)
		}
	}

	return nil

}

func ajouterColonneMachineDansTables(confTable []ConfigTableBDD) []ConfigTableBDD {
	for i := range confTable {
		confTable[i].Colonnes = append(confTable[i].Colonnes, ConfigColonneBDD{Nom: AQUA_MACHINE, Contenu: AQUA_MACHINE, Type: "VARCHAR(20)"})
	}
	return confTable
}
