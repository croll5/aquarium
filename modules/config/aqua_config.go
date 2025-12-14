package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"
)

type AquaConfig struct {
	DebutAnalyse     time.Time                    `json:"debut_analyse"`
	EntitesImpactees string                       `json:"entites_impactees"`
	FinAnalyse       time.Time                    `json:"fin_analyse"`
	Analyste         string                       `json:"analyste"`
	Machines         map[string]AquaConfigMachine `json:"machines"`
	MainCourante     []AquaConfigEvenement        `json:"main_courante"`
	DonneesFuitees   []AquaConfigFuiteDP          `json:"donnees_fuitees"`
	Contacts         []AquaConfigContact          `json:"contacts"`
	DossierAnalyse   string
}

type AquaConfigEvenement struct {
	Horodatage  time.Time `json:"horodatage"`
	Description string    `json:"description"`
}

type AquaConfigFuiteDP struct {
	NatureDonnees string `json:"nature_donnees"`
	Emplacement   string `json:"emplacement"`
}

type AquaConfigContact struct {
	Genre     string `json:"genre"`
	Prenom    string `json:"prenom"`
	Nom       string `json:"nom"`
	Qualite   string `json:"qualite"`
	Telephone string `json:"telephone"`
	Courriel  string `json:"courriel"`
}

type AquaConfigMachine struct {
	Nom              string              `json:"nom"`
	NatureEquipement string              `json:"nature_equipement"`
	Adresses         []AquaConfigAdresse `json:"adresses"`
	Fichiers         []string
	AAnalyser        bool   `json:"a_analyser"`
	Config           string `json:"config"`
}

type AquaConfigAdresse struct {
	Adresse     string `json:"adresse"`
	Commentaire string `json:"commentaire"`
}

var ANALYSE_AQUA string = "analyse.aqua"

func GetAquaConfig(cheminProjet string) (AquaConfig, error) {
	// Ouverture du fichier de configuration
	fichier, err := os.Open(filepath.Join(cheminProjet, ANALYSE_AQUA))
	if err != nil {
		return AquaConfig{}, err
	}
	defer fichier.Close()
	// Lecture des données du fichier
	donneesAqua, err := io.ReadAll(fichier)
	if err != nil {
		return AquaConfig{}, err
	}
	// Interprétation du contenu du fichier
	var aquaConfig AquaConfig
	err = json.Unmarshal(donneesAqua, &aquaConfig)
	return aquaConfig, err
}

func ListeMachinesAnalysees(cheminProjet string) (map[string]AquaConfigMachine, error) {
	aquaconfig, err := GetAquaConfig(cheminProjet)
	if err != nil {
		return map[string]AquaConfigMachine{}, err
	}
	return aquaconfig.Machines, nil
}

func EnregistrerAquaConfig(cheminProjet string, config AquaConfig) error {
	// Ouverture du fichier
	var fichier *os.File
	fichier, err := os.Open(filepath.Join(cheminProjet, ANALYSE_AQUA))
	if err != nil {
		fichier, err = os.Create(filepath.Join(cheminProjet, ANALYSE_AQUA))
		if err != nil {
			return err
		}
	}
	defer fichier.Close()
	// Génération des données
	donnees, err := json.Marshal(config)
	if err != nil {
		return err
	}
	// Enregistrement des données
	_, err = fichier.Write(donnees)
	return err
}
