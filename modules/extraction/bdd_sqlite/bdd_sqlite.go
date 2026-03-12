package bdd_sqlite

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"database/sql"
	"io"
	"os"

	"github.com/pkg/errors"
)

type SQLite struct{}

func (sq SQLite) Extraction(cheminProjet string, fichier io.Reader, cheminFichierAExtraire string, configExtraction config.ConfigExtraction, idMachine string) error {
	// On commence par créer un fichier temporaire qui contienra les données
	baseSqlite, err := os.CreateTemp(cheminProjet, "base_sqlite")
	if err != nil {
		return errors.WithStack(err)
	}
	// On ajoute les données dans la base SQLite
	_, err = io.Copy(baseSqlite, fichier)
	baseSqlite.Close()
	if err != nil {
		os.Remove(baseSqlite.Name())
		return errors.WithStack(err)
	}
	// On ouvre le fichier avec sqlite
	bdd, err := sql.Open("sqlite", baseSqlite.Name())
	if err != nil {
		os.Remove(baseSqlite.Name())
		return errors.WithStack(err)
	}
	defer supprimerBDDTemp(bdd, baseSqlite.Name())
	reponse, err := bdd.Query(configExtraction.Complement["requete_sql"])
	if err != nil {
		return errors.WithStack(err)
	}
	defer reponse.Close()
	// On liste les colonnes
	err = traiterReponseBDD(cheminProjet, cheminFichierAExtraire, reponse, configExtraction, idMachine)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func supprimerBDDTemp(bdd *sql.DB, cheminBDD string) error {
	err := bdd.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.Remove(cheminBDD)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func traiterReponseBDD(cheminProjet string, source string, reponse *sql.Rows, configExtraction config.ConfigExtraction, idMachine string) error {
	// On récupère les noms des colonnes
	listeColonnes, err := reponse.Columns()
	if err != nil {
		return errors.WithStack(err)
	}
	// On fait la liste des noms des colonnes
	var contenuColonnesAttendues []string = []string{}
	var nomColonnesAttendues []string = []string{}
	for _, colonneVoulue := range configExtraction.Table[0].Colonnes {
		contenuColonnesAttendues = append(contenuColonnesAttendues, colonneVoulue.Contenu)
		nomColonnesAttendues = append(nomColonnesAttendues, colonneVoulue.Nom)
	}
	// On prépare la requête d'insertion dans la BDD
	var abase = aquabase.InitDB_Extraction(cheminProjet)
	var requeteInsertion = abase.InitRequeteInsertionExtraction(configExtraction.Table[0].Nom, nomColonnesAttendues)
	// On remplit la requête
	var contenuLigne []interface{} = make([]interface{}, len(listeColonnes))
	var pointeursColonnes []interface{} = make([]interface{}, len(listeColonnes))
	for i := range listeColonnes {
		pointeursColonnes[i] = &contenuLigne[i]
	}
	for reponse.Next() {
		var valeursColonnes []interface{} = make([]interface{}, len(contenuColonnesAttendues))
		err := reponse.Scan(pointeursColonnes...)
		if err != nil {
			return errors.WithStack(err)
		}
		for i := range contenuColonnesAttendues {
			if contenuColonnesAttendues[i] == "aqua_source" {
				valeursColonnes[i] = source
				continue
			}
			if contenuColonnesAttendues[i] == config.AQUA_MACHINE {
				valeursColonnes[i] = idMachine
				continue
			}
			for j := range listeColonnes {
				if contenuColonnesAttendues[i] == listeColonnes[j] {
					valeursColonnes[i] = contenuLigne[j]
				}
			}
		}
		err = requeteInsertion.AjouterDansRequete(valeursColonnes...)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	// Instertion des valeurs dans la base de données
	err = requeteInsertion.Executer()
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}
