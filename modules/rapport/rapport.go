/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à [rappeler les
caractéristiques techniques de votre logiciel].

Ce logiciel est régi par la licence CeCILL soumise au droit français et
respectant les principes de diffusion des logiciels libres. Vous pouvez
utiliser, modifier et/ou redistribuer ce programme sous les conditions
de la licence CeCILL telle que diffusée par le CEA, le CNRS et l'INRIA
sur le site "http://www.cecill.info".

En contrepartie de l'accessibilité au code source et des droits de copie,
de modification et de redistribution accordés par cette licence, il n'est
offert aux utilisateurs qu'une garantie limitée.  Pour les mêmes raisons,
seule une responsabilité restreinte pèse sur l'auteur du programme,  le
titulaire des droits patrimoniaux et les concédants successifs.

A cet égard  l'attention de l'utilisateur est attirée sur les risques
associés au chargement,  à l'utilisation,  à la modification et/ou au
développement et à la reproduction du logiciel par l'utilisateur étant
donné sa spécificité de logiciel libre, qui peut le rendre complexe à
manipuler et qui le réserve donc à des développeurs et des professionnels
avertis possédant  des  connaissances  informatiques approfondies.  Les
utilisateurs sont donc invités à charger  et  tester  l'adéquation  du
logiciel à leurs besoins dans des conditions permettant d'assurer la
sécurité de leurs systèmes et ou de leurs données et, plus généralement,
à l'utiliser et l'exploiter dans les mêmes conditions de sécurité.

Le fait que vous puissiez accéder à cet en-tête signifie que vous avez
pris connaissance de la licence CeCILL, et que vous en avez accepté les
termes.
*/

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
