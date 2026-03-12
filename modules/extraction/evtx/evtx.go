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

package evtx

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"bytes"
	"io"

	"github.com/0xrawsec/golang-evtx/evtx"

	"github.com/pkg/errors"
)

type Evtx struct{}

// -------------------------- FONCTIONS LOCALES -------------------------- //

/*
	Fonction qui, à partir d'un évènement, va ajouter à la requête ses caractéristiques

@param evenement : un pointeur vers l'évènement à ajouter
@param requeteInsertionEvtx : la requete de base de données à laquelle on veut l'ajouter
@param fichierSource : le chemin vers le fichier source
@return : une erreur s'il y a eu des problèmes dans l'extraction des caractéristiques de l'évènement
*/
func ajouterGoEvtxMapDansBDD(evenement *evtx.GoEvtxMap, requeteInsertionEvtx *aquabase.RequeteInsertion, fichierSource string, configExtraction config.ConfigExtraction, idMachine string) error {
	var listeContenuColonnes []interface{} = make([]interface{}, 0)
	for _, colonne := range configExtraction.Table[0].Colonnes {
		if colonne.Contenu == "horodatage" {
			listeContenuColonnes = append(listeContenuColonnes, evenement.TimeCreated())
		} else if colonne.Contenu == "aqua_source" {
			listeContenuColonnes = append(listeContenuColonnes, fichierSource)
		} else if colonne.Contenu == config.AQUA_MACHINE {
			listeContenuColonnes = append(listeContenuColonnes, idMachine)
		} else if colonne.Contenu == "message" {
			chemin := evtx.GoEvtxPath{"Event", "EventData"}
			infosEvenement, err := evenement.Get(&chemin)
			if err != nil {
				listeContenuColonnes = append(listeContenuColonnes, "[AQUA] Problème dans l'extraction du message : "+err.Error())
				continue
			}
			var infosEvenementJson []byte = evtx.ToJSON(infosEvenement)
			listeContenuColonnes = append(listeContenuColonnes, string(infosEvenementJson))
		} else {
			var chemin evtx.GoEvtxPath = evtx.Path(colonne.Contenu)
			valeur, err := evenement.GetString(&chemin)
			if err != nil {
				valeur = "NaN"
			}
			listeContenuColonnes = append(listeContenuColonnes, valeur)
		}
	}
	err := requeteInsertionEvtx.AjouterDansRequete(listeContenuColonnes...)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

/*
Fonction qui, à partir d'un fichier evtx zippé, ajoute tous ses évènements à la base de données

@param cheminProjet : le chemin de l'analyse aquarium
@param fichier : le fichier d'évènements
@param bd : un pointeur vers la base de données d'analyse
@param cheminTemp : le chemin vers un répertoire temporaire
@param fichierSource : le chemin du fichier evtx à extraire
*/
func (e Evtx) extraireEvenementsDepuisTampon(cheminProjet string, fichier io.Reader, fichierSource string, configExtraction config.ConfigExtraction, idMachine string) error {
	// On copie le contenu du fichier dans un tampon
	var tampon bytes.Buffer
	if _, err := io.Copy(&tampon, fichier); err != nil {
		return errors.WithStack(err)
	}
	// On ouvre le tampon avec la bibliothèque evtx
	readerAt := bytes.NewReader(tampon.Bytes())
	var fichierEvtx evtx.File
	fichierEvtx, err := evtx.New(readerAt)
	if err != nil {
		return errors.WithStack(err)
	}
	// On récupère la liste des évènements
	listeEvenements := fichierEvtx.FastEvents()
	var probleme error = nil
	// On prépare la requête d'insertion dans la BDD
	var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	var listeColonnesEvtx []string = []string{}
	// On répurère la liste des colonnes à extraire
	for _, colonne := range configExtraction.Table[0].Colonnes {
		listeColonnesEvtx = append(listeColonnesEvtx, colonne.Nom)
	}
	// On prépare le contenu qui sera inséré dans la table
	var requeteInsertionEvtx aquabase.RequeteInsertion = abase.InitRequeteInsertionExtraction("Evtx", listeColonnesEvtx)
	for evenement := range listeEvenements {
		// On ajoute chaque évènement à la requete
		err := ajouterGoEvtxMapDansBDD(evenement, &requeteInsertionEvtx, fichierSource, configExtraction, idMachine)
		if err != nil {
			probleme = errors.WithStack(err)
		}
	}
	// On exécute la requete
	err = requeteInsertionEvtx.Executer()
	// Si l'on n'a pas pu l'exécuter, on renvoie une erreur
	if err != nil {
		return errors.WithStack(err)
	}
	return probleme
}

// ------------------------- FONCTIONS GLOBALES ------------------------- //

/* Fonction d'extraction des fichiers evtx */
func (e Evtx) Extraction(cheminProjet string, fichier io.Reader, nomFichier string, configExtraction config.ConfigExtraction, idMachine string) error {
	err := e.extraireEvenementsDepuisTampon(cheminProjet, fichier, nomFichier, configExtraction, idMachine)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}
