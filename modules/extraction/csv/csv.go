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
	"aquarium/modules/config"
	"aquarium/modules/extraction/utilitaires"
	"bufio"
	"encoding/csv"
	"io"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	_ "modernc.org/sqlite"

	"github.com/pkg/errors"
)

type Csv struct{}

/* ******************************************************************** */
/* ********************** Csv Methods ***************************** */
/* ******************************************************************** */

func (gt Csv) Extraction(cheminProjet string, fichier io.Reader, cheminFichierAExtraire string, configExtraction config.ConfigExtraction, idMachine string) error {
	// on initialise la requête
	requeteInsertion := utilitaires.CreerRequeteInstertionDepuisConfig(cheminProjet, idMachine, &configExtraction.Table[0])
	scanner := bufio.NewReader(fichier)
	// On lit le fichier CSV
	lecteurCSV := csv.NewReader(scanner)
	enTete, err := lecteurCSV.Read()
	if err != nil {
		return errors.WithStack(err)
	}
	fonctionTraitement := getFonctionTraitementLigne(enTete, configExtraction.Table[0].Colonnes, cheminFichierAExtraire, idMachine)
	for i := 0; err == nil; i++ {
		ligne, err := lecteurCSV.Read()
		if err == io.EOF {
			break
		}
		if i > 0 && i%100_000 == 0 {
			requeteInsertion.Executer()
			requeteInsertion = utilitaires.CreerRequeteInstertionDepuisConfig(cheminProjet, idMachine, &configExtraction.Table[0])
		}
		requeteInsertion.AjouterDansRequete(fonctionTraitement(ligne)...)
	}
	err = requeteInsertion.Executer()
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

/* **************************************************************************** */
/* *********************** Csv Utils Methods ****************************** */
/* **************************************************************************** */

func getFonctionTraitementLigne(enTete []string, colonnesTable []config.ConfigColonneBDD, source string, idMachine string) func([]string) []interface{} {
	var fonctionsTraitement []func(valeurs []string) interface{} = make([]func([]string) interface{}, len(colonnesTable))
	for i, colonne := range colonnesTable {
		trouve := false
		for j, colonneCSV := range enTete {
			if colonneCSV == colonne.Contenu {
				fonctionsTraitement[i] = func(valeurs []string) interface{} {
					return valeurs[j]
				}
				trouve = true
			}
		}
		if !trouve {
			switch colonne.Contenu {
			case "aqua_source":
				fonctionsTraitement[i] = func(valeurs []string) interface{} {
					return source
				}
			case config.AQUA_MACHINE:
				fonctionsTraitement[i] = func(valeurs []string) interface{} {
					return idMachine
				}
			default:
				fonctionsTraitement[i] = func(valeurs []string) interface{} {
					return "[AQUA_ERR] Colonne " + colonne.Contenu + " non trouvée"
				}
			}
		}
	}
	return func(valeurs []string) []interface{} {
		var resultat []interface{} = make([]interface{}, len(colonnesTable))
		for i, fonctionTraitement := range fonctionsTraitement {
			resultat[i] = fonctionTraitement(valeurs)
		}
		return resultat
	}
}

func exportDfToDb(df dataframe.DataFrame, cheminProjet string, filname string, configTable config.ConfigTableBDD, idMachine string) error {
	// On initialise la requête
	requeteInsertion := utilitaires.CreerRequeteInstertionDepuisConfig(cheminProjet, idMachine, &configTable)
	// On parcourt le dataframe
	for _, ligneCSV := range df.Maps() {
		var valeursAAjouter []interface{} = make([]interface{}, len(configTable.Colonnes))
		for i, colonne := range configTable.Colonnes {
			valeur, ok := ligneCSV[colonne.Contenu]
			if !ok {
				switch colonne.Contenu {
				case "aqua_source":
					valeur = filname
				case config.AQUA_MACHINE:
					valeur = idMachine
				default:
					valeur = "[AQUA_ERREUR] - Nom de colonne non reconnu"
				}
			}
			valeursAAjouter[i] = valeur
		}
		err := requeteInsertion.AjouterDansRequete(valeursAAjouter...)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	err := requeteInsertion.Executer()
	if err != nil {
		return errors.WithStack(err)
	}
	return err
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
