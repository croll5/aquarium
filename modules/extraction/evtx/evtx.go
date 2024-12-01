package evtx

import "log"

type Evtx struct{}

func (e Evtx) Extraction(cheminProjet string) error {
	log.Println("Bonjour, je suis censé faire des extractions {Evtx}")
	return nil
}

func (e Evtx) Description() string {
	return "Évènements Windows (fichier .evtx) [NULL]"
}

func (e Evtx) PrerequisOK(cheminORC string) bool {
	return true
}
