package rapport

import (
	"aquarium/modules/aquabase"
	"path/filepath"
)

type Rapport struct {
	cheminProjet string
	bdd          *aquabase.Aquabase
}

var colonnesTableEtapes []string = []string{"commentaire", "idPiste"}
var colonnesTablePistes []string = []string{"titre", "description", "conslusion"}

func InitRapport(cheminProjet string) *Rapport {
	// On s'assure que toutes les bases nécessaires ont été créées
	var base *aquabase.Aquabase = aquabase.Init(filepath.Join(cheminProjet, "analyse", "rapport.db"))
	base.CreateTableIfNotExist1("pistes", colonnesTablePistes, true)
	base.CreateTableIfNotExist1("etapes", colonnesTableEtapes, true)
	// On crée l'objet rapport
	var rapport Rapport = Rapport{cheminProjet: cheminProjet, bdd: base}
	return &rapport
}

func (rapport *Rapport) AjouterCommentaire(idPiste int, commentaire string) error {
	var requeteInsersion aquabase.RequeteInsertion = rapport.bdd.InitRequeteInsertionExtraction("etapes", colonnesTableEtapes)
	err := requeteInsersion.AjouterDansRequete(commentaire, idPiste)
	if err != nil {
		return err
	}
	err = requeteInsersion.Executer()
	return err
}

func (rapport *Rapport) CreerTables() error {
	err := rapport.bdd.CreateTableIfNotExist1("pistes", colonnesTablePistes, true)
	if err != nil {
		return err
	}
	err = rapport.bdd.CreateTableIfNotExist1("etapes", colonnesTableEtapes, true)
	return err
}
