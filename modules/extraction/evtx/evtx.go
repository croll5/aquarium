package evtx

import "log"

type Evtx struct {
	extrait bool
}

func (e Evtx) Extraction() error {
	log.Println("Bonjour, je suis censé faire des extractions")
	return nil
}

func New() Evtx {
	return Evtx{extrait: false}
}

func (e Evtx) Description() string {
	return "Évènements Windows (fichier .evtx)"
}
