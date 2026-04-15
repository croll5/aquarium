/*
Copyright Cécile Rolland, (8 avril 2026)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l’analyse des collectes
traçologiques effectuées avec le logiciel DFIR-ORC.

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
