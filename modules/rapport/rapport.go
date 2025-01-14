package rapport

import (
	"aquarium/modules/aquabase"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
)

type Rapport struct {
	cheminProjet string
	bdd          *aquabase.Aquabase
}

type EtapeAnalyse struct {
	RequeteSQL  interface{}
	Commentaire interface{}
	NomTable    string
	//LignesTable []map[string]interface{}
}

var colonnesTableEtapes []string = []string{"commentaire", "idPiste", "requeteSQL"}
var colonnesTablePistes []string = []string{"titre", "description", "conslusion"}

func InitRapport(cheminProjet string) *Rapport {
	// On s'assure que toutes les bases nécessaires ont été créées
	var base *aquabase.Aquabase = aquabase.Init(filepath.Join(cheminProjet, "analyse", "rapport.db"))
	// On crée l'objet rapport
	var rapport Rapport = Rapport{cheminProjet: cheminProjet, bdd: base}
	return &rapport
}

func (rapport *Rapport) AjouterEtape(idPiste string, commentaire string, requeteSQL string, lignesAEnregistrer []map[string]interface{}) error {
	var requeteInsersionEtape aquabase.RequeteInsertion = rapport.bdd.InitRequeteInsertionExtraction("etapes", colonnesTableEtapes)
	err := requeteInsersionEtape.AjouterDansRequete(commentaire, idPiste, requeteSQL)
	if err != nil {
		return err
	}
	err = requeteInsersionEtape.Executer()
	if err != nil {
		return err
	}
	var idEtape int = rapport.bdd.TailleRequeteSQL("SELECT * FROM etapes")
	log.Println(idEtape)
	// Création de la table à enregistrer
	log.Println(lignesAEnregistrer)
	if len(lignesAEnregistrer) < 1 {
		return errors.New("la table à enregistrer ne contient aucune ligne")
	}
	var nomTableAEnregistrer string = "enregistrement_" + strconv.Itoa(idEtape)
	var listeColonnesTableEtudiee []string = []string{}
	for cle := range lignesAEnregistrer[0] {
		listeColonnesTableEtudiee = append(listeColonnesTableEtudiee, cle)
	}
	err = rapport.bdd.CreateTableIfNotExist1(nomTableAEnregistrer, listeColonnesTableEtudiee, false)
	if err != nil {
		return err
	}
	// Ajout des valeurs à enregistrer
	var requeteInsertionValeurs aquabase.RequeteInsertion = rapport.bdd.InitRequeteInsertionExtraction(nomTableAEnregistrer, listeColonnesTableEtudiee)
	for _, valeurs := range lignesAEnregistrer {
		var valeursAInserer []interface{} = make([]interface{}, len(listeColonnesTableEtudiee))
		for i := 0; i < len(listeColonnesTableEtudiee); i++ {
			valeursAInserer[i] = valeurs[listeColonnesTableEtudiee[i]]
		}
		err = requeteInsertionValeurs.AjouterDansRequete(valeursAInserer...)
		if err != nil {
			return err
		}
	}
	err = requeteInsertionValeurs.Executer()
	//err = rapport.bdd.EnregistrerTableDepuisMap(requeteSQL, lignesAEnregistrer, "enregistrement"+strconv.Itoa(idEtape))
	return err
}

/* FONCTIONS DE MODIFICATION DU RAPPORT */

func (rapport *Rapport) AjouterPiste(titre string, description string) error {
	var requeteInsertion aquabase.RequeteInsertion = rapport.bdd.InitRequeteInsertionExtraction("pistes", []string{"titre", "description"})
	err := requeteInsertion.AjouterDansRequete(titre, description)
	if err != nil {
		return err
	}
	err = requeteInsertion.Executer()
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

/* FONCION D4AFFICHAGE DU RAPPORT */

func (rapport *Rapport) GetPistes() []map[string]interface{} {
	return rapport.bdd.ResultatRequeteSQL("SELECT * FROM pistes ORDER BY id")
}

func (rapport *Rapport) GetEtapesPiste(idPiste int) []EtapeAnalyse {
	var listeEtapes []map[string]interface{} = rapport.bdd.ResultatRequeteSQL("SELECT * FROM etapes WHERE idPiste=" + strconv.Itoa(idPiste))
	var listeEtapesAnalyse []EtapeAnalyse = []EtapeAnalyse{}
	for _, valeur := range listeEtapes {
		//var lignesTable []map[string]interface{} = rapport.bdd.SelectAllFrom(, 10000000)
		//listeTables = append(listeTables, lignesTable)
		var etapeAnalyse EtapeAnalyse = EtapeAnalyse{RequeteSQL: valeur["requeteSQL"], Commentaire: valeur["commentaire"], NomTable: fmt.Sprintf("enregistrement_%v", valeur["id"])}
		listeEtapesAnalyse = append(listeEtapesAnalyse, etapeAnalyse)
	}
	return listeEtapesAnalyse
}

func (rapport *Rapport) GetDonnesTableSauvegardee(nomTable string) []map[string]interface{} {
	return rapport.bdd.SelectAllFrom(nomTable, 1000000)
}
