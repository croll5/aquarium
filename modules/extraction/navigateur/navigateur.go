package navigateur

type Navigateur struct {
	extrait bool
}

func (n Navigateur) Extraction() error {
	return nil
}

func New() Navigateur {
	return Navigateur{extrait: false}
}

func (n Navigateur) Description() string {
	return "Historique de navigation"
}
