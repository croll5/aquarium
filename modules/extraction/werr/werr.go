package werr

type Werr struct{}

func (w Werr) Extraction(cheminProjet string) error {
	return nil
}

func (w Werr) PrerequisOK(cheminORC string) bool {
	return true
}

func (w Werr) Description() string {
	return "Fichier d'erreurs Windows"
}
