package werr

import "log"

type Werr struct{}

func (w Werr) Extraction(cheminProjet string) error {
	log.Println("Bonjour, je suis censé faire des extractions {Werr}")
	return nil
}

func (w Werr) PrerequisOK(cheminORC string) bool {
	return true
}

func (w Werr) Description() string {
	return "Fichier d'erreurs Windows"
}
