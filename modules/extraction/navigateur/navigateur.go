package navigateur

import "log"

type Navigateur struct{}

func (n Navigateur) Extraction(cheminProjet string) error {
	log.Println("Bonjour, je suis censé faire des extractions {Navigateur}")
	return nil
}

func (n Navigateur) Description() string {
	return "Historique de navigation [NULL]"
}

func (n Navigateur) PrerequisOK(cheminORC string) bool {
	return true
}
