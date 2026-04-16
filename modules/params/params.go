package params

import (
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
)

const (
	DOSSIER_CONFIG = "config"
	FICHIER_PARAMS = "params.xml"
)

type ParametresXML struct {
	XMLName        xml.Name `xml:"parametres" json:"-"`
	Contrastes     bool     `xml:"contrastes"`
	Dyslexie       bool     `xml:"dyslexie"`
	NonAuxBubulles bool     `xml:"non_aux_bubulles"`
}

func SauvegarderParametres(cheminBase string, contrastes bool, dyslexie bool, nonAuxBubulles bool) error {
	parametres := ParametresXML{
		Contrastes:     contrastes,
		Dyslexie:       dyslexie,
		NonAuxBubulles: nonAuxBubulles,
	}
	contenuXML, err := xml.MarshalIndent(parametres, "", "  ")
	if err != nil {
		return err
	}
	cheminConfig := filepath.Join(cheminBase, DOSSIER_CONFIG)
	if err = os.MkdirAll(cheminConfig, 0o755); err != nil {
		return err
	}
	cheminFichier := filepath.Join(cheminConfig, FICHIER_PARAMS)
	contenuFinal := append([]byte(xml.Header), contenuXML...)
	return os.WriteFile(cheminFichier, contenuFinal, 0o644)
}

func ChargerParametres(cheminBase string) (ParametresXML, error) {
	parametres := ParametresXML{}
	cheminFichier := filepath.Join(cheminBase, DOSSIER_CONFIG, FICHIER_PARAMS)
	contenu, err := os.ReadFile(cheminFichier)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return parametres, nil
		}
		return parametres, err
	}
	err = xml.Unmarshal(contenu, &parametres)
	if err != nil {
		return ParametresXML{}, err
	}
	return parametres, nil
}
