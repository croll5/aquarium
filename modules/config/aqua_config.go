package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"
)

type AquaConfig struct {
	DebutAnalyse time.Time                    `json:"debut_analyse"`
	Description  string                       `json:"description"`
	FinAnalyse   time.Time                    `json:"fin_analyse"`
	Auteur       string                       `json:"nom_auteur"`
	Machines     map[string]AquaConfigMachine `json:"machines"`
}

type AquaConfigMachine struct {
	Nom    string `json:"nom"`
	Config string `json:"config"`
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
