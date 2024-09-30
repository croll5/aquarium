package navigateur

type Navigateur struct{}

func (n Navigateur) Extraction(cheminProjet string) error {
	return nil
}

func (n Navigateur) Description() string {
	return "Historique de navigation"
}

func (n Navigateur) PrerequisOK(cheminORC string) bool {
	return true
}
