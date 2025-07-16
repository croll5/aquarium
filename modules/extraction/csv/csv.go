/*
Copyright ou © ou Copr. Charles Mailley, (21 janvier 2025)
Modification : Cécile Rolland, (13 mai 2025)

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

package csv

// compiler du go         : go build csv.go
// execution du programme : ./csv.exe
import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	_ "modernc.org/sqlite"
)

type Csv struct{}

/* ******************************************************************** */
/* ********************** Csv Methods ***************************** */
/* ******************************************************************** */

func (gt Csv) Extraction(cheminProjet string, fichier bytes.Buffer, cheminFichierAExtraire string, configExtraction config.ConfigExtraction) error {
	var df dataframe.DataFrame
	var err error
	// On lit les données du fichier CSV
	df = dataframe.ReadCSV(&fichier)
	if err != nil {
		return err
	}
	err = exportDfToDb(df, cheminProjet, cheminFichierAExtraire, configExtraction.Table.Nom, configExtraction.Table.Colonnes)
	return err
}

/* **************************************************************************** */
/* *********************** Csv Utils Methods ****************************** */
/* **************************************************************************** */

func exportDfToDb(df dataframe.DataFrame, cheminProjet string, filname string, tableName string, colonnesTable []config.ConfigColonneBDD) error {
	adb := aquabase.InitDB_Extraction(cheminProjet)

	// On crée un dictionnaire des colonnes de la table
	var nouveauNomColonne map[string]string = map[string]string{}
	var listeContenuColonnes []string = []string{}
	for _, colonne := range colonnesTable {
		// La colonne source est un peu particulière car elle n'est pas dans le csv
		if colonne.Contenu == "source" {
			df = DfAddColumn(df, "source", filname)
		}
		nouveauNomColonne[colonne.Contenu] = colonne.Nom
		listeContenuColonnes = append(listeContenuColonnes, colonne.Contenu)
	}
	fmt.Println("Import Csv to DB: " + filname)

	// On filtre sur les colonnes qui doivent être prises
	columns := listItemsInList(listeContenuColonnes, df.Names())
	// Select the specified columns
	df = df.Select(columns)
	// On renomme les colonnes
	for _, ancien := range df.Names() {
		// TODO: Vérifier que le tableau a cette colonne
		if ancien != nouveauNomColonne[ancien] {
			df = df.Rename(nouveauNomColonne[ancien], ancien)
		}
	}
	//columns := df.Names()         // For no columns filter
	//columnSelection := df.Names() // For no columns filter

	// Check the table exist
	err := adb.CreateTableIfNotExist1(tableName, df.Names(), true)
	if err != nil {
		log.Println("Erreur dans la création de la table :", err)
		return fmt.Errorf("ERROR: exportDfToDb(_) [createTableIfNotExist]: %w", err)
	}
	// export data
	err = adb.SaveDf(df, tableName)
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDb(_) [AjoutEvenementDansBDD]: %w", err)
	}
	return nil
}

/* ******************************************************************** */
/* *********************** Utils Methods ****************************** */
/* ******************************************************************** */

/** Find and return all elements of smallList existing in bigList **/
func listItemsInList(smallList, bigList []string) []string {
	var result []string
	for _, smallItem := range smallList {
		for _, bigItem := range bigList {
			if strings.Contains(bigItem, smallItem) {
				result = append(result, smallItem)
				break
			}
		}
	}
	return result
}

func DfHead(df dataframe.DataFrame, nFirstRows int) dataframe.DataFrame {
	indices := make([]int, nFirstRows)
	for i := 0; i < nFirstRows; i++ {
		indices[i] = i
	}
	return df.Subset(indices)
}

func DfAddColumn(df dataframe.DataFrame, colname string, value string) dataframe.DataFrame {
	sourceColumn := make([]string, df.Nrow())
	for i := range sourceColumn {
		sourceColumn[i] = value
	}
	return df.Mutate(series.New(sourceColumn, series.String, colname))
}
