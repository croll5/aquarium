/*
Copyright Cécile Rolland, (17 août 2025)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l'analyse des collectes
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

package chronologie

import (
	"aquarium/modules/aquabase"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var machineEnCache = ""
var cacheChronologie []map[string]interface{} = []map[string]interface{}{}

func ContenuEvenementsChronologie(cheminProjet string, idMachine string, indexDebut int64, taille int64) ([]map[string]interface{}, error) {
	var resultat []map[string]interface{} = []map[string]interface{}{}
	adb := aquabase.InitDB_Extraction(cheminProjet)
	// Récupération des identifiants des évènements classés par horodatage
	listeIdEvenements, err := adb.ResultatRequeteSQLAvecFiltres(&aquabase.ParametresRequeteSQL{
		NomTable:        aquabase.TABLE_CHRONOLOGIE_GLOBALE,
		Colonnes:        []string{"id_evenement", "horodatage", "nom_table"},
		FiltresColonnes: []aquabase.FiltreRequeteSQL{aquabase.FiltreRequeteSQL{NomColonne: "id_machine", Valeur: idMachine}},
		Limit:           taille,
		Offset:          indexDebut,
		OrderBy:         "horodatage",
	})
	if err != nil {
		return resultat, errors.WithStack(err)
	}
	// Pour chaque identifiant, récupération de l’évènement associé
	resultat, err = aquaChronologieVersTableComplete(cheminProjet, listeIdEvenements)
	if err != nil {
		return resultat, errors.WithStack(err)
	}
	go adb.CreerIndex(aquabase.TABLE_CHRONOLOGIE_GLOBALE, []string{"id_evenement", "nom_table"}, aquabase.TABLE_CHRONOLOGIE_GLOBALE+"index_2")
	return resultat, nil
}

func PositionDateDansChronologie(cheminProjet string, idMachine string, dateSelectionne time.Time) (int64, error) {
	adb := aquabase.InitDB_Extraction(cheminProjet)
	resultatTaille, err := adb.ResultatRequeteSQLAvecFiltres(&aquabase.ParametresRequeteSQL{
		NomTable: aquabase.TABLE_CHRONOLOGIE_GLOBALE,
		Colonnes: []string{"count(*) AS nb"},
		FiltresColonnes: []aquabase.FiltreRequeteSQL{aquabase.FiltreRequeteSQL{
			NomColonne: "horodatage",
			Valeur:     dateSelectionne,
			Inferieur:  true,
		}},
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return resultatTaille[0]["nb"].(int64), nil
}

func ChercherIndicateurDansBase(cheminProjet string, valeurRecherchee string, nbEvenements int) ([]map[string]interface{}, error) {
	var resultat []map[string]interface{}
	adb := aquabase.InitDB_Extraction(cheminProjet)
	listeTables, err := adb.SchemaBDD()
	if err != nil {
		return resultat, err
	}
	requete := "SELECT * FROM aqua_chronologie WHERE "
	var listeConditions []string = make([]string, len(listeTables))
	i := 0
	for nomTable, listeColonnes := range listeTables {
		listeConditions[i] = "nom_table=\"" + nomTable + "\" AND id_evenement IN (SELECT id FROM " + nomTable + " WHERE "
		listeConditions[i] += strings.Join(listeColonnes, " LIKE \"%"+strings.ReplaceAll(valeurRecherchee, "\"", "\"\"")+"%\" OR ")
		listeConditions[i] += " LIKE \"%" + strings.ReplaceAll(valeurRecherchee, "\"", "\"\"") + "%\")"
		i++
	}
	requete += strings.Join(listeConditions, " OR ")
	requete += " ORDER BY horodatage"
	listeIdEvenements, err := adb.ResultatRequeteSQL(requete)
	if err != nil {
		return resultat, errors.WithStack(err)
	}
	if len(listeIdEvenements) > nbEvenements {
		cacheChronologie = listeIdEvenements
		listeIdEvenements = listeIdEvenements[:nbEvenements]
	}
	resultat, err = aquaChronologieVersTableComplete(cheminProjet, listeIdEvenements)
	if err != nil {
		return resultat, errors.WithStack(err)
	}
	return resultat, nil
}

func ChercherIndicateurEvenementSuivant(cheminProjet string, debut int, fin int) ([]map[string]interface{}, error) {
	if fin < debut {
		return []map[string]interface{}{}, errors.WithStack(errors.Errorf("Indice de début %d supérieur à l’indice de fin %d", debut, fin))
	}
	if fin >= len(cacheChronologie) {
		return []map[string]interface{}{}, nil
	}
	return aquaChronologieVersTableComplete(cheminProjet, cacheChronologie[debut:fin])
}

func extraireIndicateurs(cheminProjet string, evenements []map[string]interface{}) error {
	dicoRegex := make(map[string]*regexp.Regexp)
	var err error
	dicoRegex["sha1"], err = regexp.Compile(`[A-Fa-f0-9]{40}`)
	if err != nil {
		return err
	}
	dicoRegex["ip"], err = regexp.Compile(`[1-2]{0,1}[0-9]{1,2}\.[1-2]{0,1}[0-9]{1,2}\.[1-2]{0,1}[0-9]{1,2}\.[1-2]{0,1}[0-9]{1,2}`)
	dicoRegex["md5"], err = regexp.Compile(`[A-Fa-f0-9]{32}`)
	for _, evenement := range evenements {
		for nomChamp, valeurChamp := range evenement {
			valeurChampStr := fmt.Sprint(valeurChamp)
			for typeRegex, moteurRegex := range dicoRegex {
				res := moteurRegex.FindAllString(valeurChampStr, -1)
				indexTrouves := moteurRegex.FindAllStringIndex(valeurChampStr, -1)
				for _, indexTrouve := range indexTrouves {
					if len(indexTrouve) != 2 {
						break
					}
					debutOk := indexTrouve[0] == 0 || valeurChampStr[indexTrouve[0]-1] == ' '
					finOk := indexTrouve[1] == len(valeurChampStr) || valeurChampStr[indexTrouve[1]] == ' '
					if debutOk && finOk {
						log.Println("Pour le champs ", nomChamp, " : ", typeRegex, " : ", res, "\n", valeurChampStr, "\n", indexTrouve)
					}
				}
			}
		}
	}
	return nil
}

func aquaChronologieVersTableComplete(cheminProjet string, listeIdEvenements []map[string]interface{}) ([]map[string]interface{}, error) {
	var resultat []map[string]interface{} = []map[string]interface{}{}
	adb := aquabase.InitDB_Extraction(cheminProjet)
	for _, idEvenement := range listeIdEvenements {
		if idEvenement["id_evenement"] == nil || idEvenement["id_evenement"] == "" || idEvenement["nom_table"] == nil || idEvenement["nom_table"] == "" {
			continue
		}
		resultatEvenement, err := adb.ResultatRequeteSQLAvecFiltres(&aquabase.ParametresRequeteSQL{
			NomTable:        fmt.Sprintf("%s", idEvenement["nom_table"]),
			FiltresColonnes: []aquabase.FiltreRequeteSQL{aquabase.FiltreRequeteSQL{NomColonne: "id", Valeur: idEvenement["id_evenement"]}},
		})
		if err != nil {
			return resultat, err
		}
		for i := range resultatEvenement {
			resultatEvenement[i]["aqua_horodatage"] = idEvenement["horodatage"]
			resultatEvenement[i]["aqua_table"] = idEvenement["nom_table"]
		}
		resultat = append(resultat, resultatEvenement...)
		// go extraireIndicateurs(cheminProjet, resultatEvenement)
	}
	return resultat, nil
}
