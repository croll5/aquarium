/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025)

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

package prefetch

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"bytes"
	"strings"

	"www.velocidex.com/golang/go-prefetch"
)

type Prefetch struct{}

func (p Prefetch) Extraction(cheminProjet string, fichier bytes.Buffer, nomFichier string, configExtraction config.ConfigExtraction, idMachine string) error {
	readerAt := bytes.NewReader(fichier.Bytes())
	infosPrechargement, err := prefetch.LoadPrefetch(readerAt)
	if err != nil {
		return err
	}
	var adb aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	for _, table := range configExtraction.Table {
		var requeteInsertion aquabase.RequeteInsertion = adb.InitRequeteInsertionExtraction(table.Nom, listeColonnesTable(table))
		if table.Condition == "" {
			var valeurs []interface{} = make([]interface{}, 0)
			for _, colonne := range table.Colonnes {
				valeurs = append(valeurs, getValeurColonne(infosPrechargement, colonne, nomFichier, idMachine))
			}
			requeteInsertion.AjouterDansRequete(valeurs...)
		} else {
			extraireValeursMultiples(infosPrechargement, table, nomFichier, &requeteInsertion, idMachine)
		}
		requeteInsertion.Executer()
	}
	return nil
}

/* FONCTIONS LOCALES */

func listeColonnesTable(table config.ConfigTableBDD) []string {
	var nomsColonnes []string = []string{}
	for _, colonne := range table.Colonnes {
		nomsColonnes = append(nomsColonnes, colonne.Nom)
	}
	return nomsColonnes
}

func getValeurColonne(fichierPrefetch *prefetch.PrefetchInfo, colonne config.ConfigColonneBDD, source string, idMachine string) interface{} {
	switch colonne.Contenu {
	case "executable":
		return fichierPrefetch.Executable
	case "empreinte":
		return fichierPrefetch.Hash
	case "version":
		return fichierPrefetch.Version
	case "taille_fichier":
		return fichierPrefetch.FileSize
	case "nb_executions":
		return fichierPrefetch.RunCount
	case "aqua_source":
		return source
	case config.AQUA_MACHINE:
		return idMachine
	case "date_execution":
		var datesExecutions []string = make([]string, len(fichierPrefetch.LastRunTimes))
		for i, date := range fichierPrefetch.LastRunTimes {
			datesExecutions[i] = date.Format("02/01/2006 ")
		}
		return strings.Join(datesExecutions, "\n")
	case "ressource":
		return strings.Join(fichierPrefetch.FilesAccessed, "\n")
	default:
		return "[AQUA] Contenu de colonne inconnu"
	}
}

func extraireValeursMultiples(infosPrechargement *prefetch.PrefetchInfo, table config.ConfigTableBDD, source string, requeteInsertion *aquabase.RequeteInsertion, idMachine string) {
	var valeursARepeter = getListeValeurs(infosPrechargement, table.Condition)
	for _, valeurARepeter := range valeursARepeter {
		var valeurs []interface{} = make([]interface{}, 0)
		for _, colonne := range table.Colonnes {
			if colonne.Contenu == table.Condition {
				valeurs = append(valeurs, valeurARepeter)
			} else {
				valeurs = append(valeurs, getValeurColonne(infosPrechargement, colonne, source, idMachine))
			}
		}
		requeteInsertion.AjouterDansRequete(valeurs...)
	}
}

func getListeValeurs(fichierPrefetch *prefetch.PrefetchInfo, contenuColonne string) []interface{} {
	switch contenuColonne {
	case "ressource":
		var resultat []interface{} = make([]interface{}, len(fichierPrefetch.FilesAccessed))
		for i, fichier := range fichierPrefetch.FilesAccessed {
			resultat[i] = fichier
		}
		return resultat
	case "date_execution":
		var resultat []interface{} = make([]interface{}, len(fichierPrefetch.LastRunTimes))
		for i, execution := range fichierPrefetch.LastRunTimes {
			resultat[i] = execution
		}
		return resultat
	}
	var resultat []interface{} = make([]interface{}, 1)
	resultat[0] = "[AQUA] Erreur dans l’extraction."
	return resultat
}
