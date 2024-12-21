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
	return "Fichier Werr [NULL]"
}

func (w Werr) CreationTable(cheminProjet string) error {
	return nil
}

func (w Werr) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	return 60
}

func (w Werr) Annuler() bool {
	return true
}

func (w Werr) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}
