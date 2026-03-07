package arborescence

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
)

func DetailsFichier(cheminProjet string, idFichier int64, idMachine string) ([]map[string]interface{}, error) {
	// On commence par récupérer la configuration du projet
	aquaConf, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	confMachine, err := config.GetConfigurationMachine(cheminProjet, idMachine, aquaConf.Machines[idMachine])
	if err != nil {
		return []map[string]interface{}{}, err
	}
	adb := aquabase.InitDB_Extraction(cheminProjet)
	return adb.SelectFrom(confMachine.Arborescence.RequeteMetadonnees, idFichier)
}
