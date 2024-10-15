package werr

type Werr struct {
	extrait bool
}

func (w Werr) Extraction() error {
	return nil
}

func New() Werr {
	return Werr{extrait: false}
}

func (w Werr) Description() string {
	return "Fichier d'erreurs Windows"
}
