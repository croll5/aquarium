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
	"time"
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
		return resultat, err
	}
	// Pour chaque identifiant, récupération de l’évènement associé
	for _, idEvenement := range listeIdEvenements {
		resultatEvenement, err := adb.ResultatRequeteSQLAvecFiltres(&aquabase.ParametresRequeteSQL{
			NomTable:        fmt.Sprintf("%s", idEvenement["nom_table"]),
			FiltresColonnes: []aquabase.FiltreRequeteSQL{aquabase.FiltreRequeteSQL{NomColonne: "id", Valeur: idEvenement["id_evenement"]}},
		})
		if err != nil {
			return resultat, err
		}
		for i := range resultatEvenement {
			resultatEvenement[i]["aqua_horodatage"] = idEvenement["horodatage"]
		}
		resultat = append(resultat, resultatEvenement...)
	}
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
		return 0, err
	}
	return resultatTaille[0]["nb"].(int64), nil
}
