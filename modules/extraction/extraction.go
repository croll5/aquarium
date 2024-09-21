package extraction

import (
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/navigateur"
	"errors"
)

type Extracteur interface {
	Extraction() error
	Description() string
}

var liste_extracteurs map[string]Extracteur = map[string]Extracteur{
	"evtx":       evtx.New(),
	"navigateur": navigateur.New(),
}

func ListeExtracteursHtml() map[string]string {
	// On itère sur tous les extracteurs
	var resultat map[string]string = map[string]string{}
	for k, v := range liste_extracteurs {
		resultat[k] = v.Description()
	}
	return resultat
}

func Extraction(module string) error {
	if liste_extracteurs[module] == nil {
		return errors.New("Erreur : module " + module + " non reconnu")
	}
	err := liste_extracteurs[module].Extraction()
	return err
}
