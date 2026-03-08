package arborescence

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"

	"github.com/pkg/errors"
)

func DetailsFichier(cheminProjet string, idFichier int64, idMachine string) ([]map[string]interface{}, error) {
	// On commence par récupérer la configuration du projet
	aquaConf, err := config.GetAquaConfig(cheminProjet)
	if err != nil {
		return []map[string]interface{}{}, errors.WithStack(err)
	}
	confMachine, err := config.GetConfigurationMachine(cheminProjet, idMachine, aquaConf.Machines[idMachine])
	if err != nil {
		return []map[string]interface{}{}, errors.WithStack(err)
	}
	adb := aquabase.InitDB_Extraction(cheminProjet)
	resultat, err := adb.SelectFrom(confMachine.Arborescence.RequeteMetadonnees, idFichier)
	if err != nil {
		return resultat, errors.WithStack(err)
	}
	return resultat, err
}
