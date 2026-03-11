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

package registre

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"aquarium/modules/extraction/utilitaires"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"strconv"
	"strings"

	"www.velocidex.com/golang/regparser"

	"github.com/pkg/errors"
)

/* VARIABLES GLOBALES */
var pourcentageChargement float32 = -1

var colonnesTableSam []string = []string{"horodatage", "idCompte", "nomCompte", "operation", "source"}

type Registre struct{}

func traiterCle(cleDeRegistre *regparser.CM_KEY_NODE, source string, requete *aquabase.RequeteInsertion, configExtraction config.ConfigExtraction, idMachine string) error {
	var listeContenuColonnes []interface{} = make([]interface{}, 0)
	// On parcourt les colonnes à ajouter
	for _, configColonne := range configExtraction.Table[0].Colonnes {
		switch configColonne.Contenu {
		case "nomCle":
			listeContenuColonnes = append(listeContenuColonnes, cleDeRegistre.Name())
		case "aqua_source":
			listeContenuColonnes = append(listeContenuColonnes, source)
		case config.AQUA_MACHINE:
			listeContenuColonnes = append(listeContenuColonnes, idMachine)
		default:
			var contenu []string = strings.Split(configColonne.Contenu, ":")
			if len(contenu) != 5 {
				listeContenuColonnes = append(listeContenuColonnes, "[ERREUR Aquarium] - Configuration du champ incorrecte")
				log.Println("Problème : ", len(contenu))
				continue
			}
			// On commence par regarder chercher les position de début et de fin
			debut, err := strconv.ParseInt(strings.TrimPrefix(contenu[2], "0x"), 16, 64)
			if err != nil {
				return err
			}
			taille, err := strconv.ParseInt(strings.TrimPrefix(contenu[3], "0x"), 16, 64)
			if err != nil {
				return err
			}
			// On cherche la valeur de la clé
			numCle, err := strconv.ParseInt(contenu[0], 16, 64)
			if err != nil {
				return err
			}
			if len(cleDeRegistre.Values()) <= int(numCle) {
				listeContenuColonnes = append(listeContenuColonnes, "[AUQA_ERREUR] Valeur non décodée dans la clé : "+cleDeRegistre.Name())
			} else {
				var DonneesCle []byte = cleDeRegistre.Values()[numCle].ValueData().Data
				if contenu[1] == "variable" {
					debut = int64(binary.LittleEndian.Uint32(DonneesCle[debut:debut+8])) + 0xCC
					taille = int64(binary.LittleEndian.Uint32(DonneesCle[taille : taille+8]))
				}
				// On ajoute les valeurs de la clé à la table
				listeContenuColonnes = append(listeContenuColonnes, utilitaires.GetFonctionDecodageBytes(contenu[4])(DonneesCle[debut:debut+taille]))
			}
		}
	}
	err := requete.AjouterDansRequete(listeContenuColonnes...)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func (s Registre) Extraction(cheminProjet string, fichier io.Reader, source string, configExtraction config.ConfigExtraction, idMachine string) error {
	// On copie dans un fichier
	var tampon bytes.Buffer
	if _, err := io.Copy(&tampon, fichier); err != nil {
		return errors.WithStack(err)
	}
	readerAt := bytes.NewReader(tampon.Bytes())
	// Ouverture du fichier comme fichier et clés de registres
	registre, err := regparser.NewRegistry(readerAt)
	if registre == nil {
		return err
	}
	if err != nil {
		return errors.WithStack(err)
	}
	// On récupère les colonnes de la table
	var nomColonnesTable []string = []string{}
	for _, colonne := range configExtraction.Table[0].Colonnes {
		nomColonnesTable = append(nomColonnesTable, colonne.Nom)
	}
	// On crée une requête d'insertion dans la BDD
	var abase aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	var requeteInsertion aquabase.RequeteInsertion = abase.InitRequeteInsertionExtraction(configExtraction.Table[0].Nom, nomColonnesTable)
	// Ouverture de la clé de registre contenant les comptes personnels
	if configExtraction.Complement["registre"] == "" {
		return errors.WithStack(errors.New("[AQUA_ERR] - Problème de configuration de l’extracteur" + configExtraction.Nom + " : valeur conplémentaire « registe non définie »."))
	}
	cleDeBase := registre.OpenKey(configExtraction.Complement["registre"])
	var probleme error
	if configExtraction.Complement["parcourir_enfants"] == "oui" {
		var enfants []*regparser.CM_KEY_NODE = cleDeBase.Subkeys()
		var listeExclusions []string = strings.Split(configExtraction.Complement["exclusions"], ";")
		for _, cleEnfant := range enfants {
			// On vérifie qu'on n'est pas dans le cas d'une exclusion
			var pasCetteCle = false
			for _, exclusion := range listeExclusions {
				if cleEnfant.Name() == exclusion {
					pasCetteCle = true
				}
			}
			if pasCetteCle {
				continue
			}
			err = traiterCle(cleEnfant, source, &requeteInsertion, configExtraction, idMachine)
			if err != nil {
				probleme = errors.WithStack(err)
			}
		}
	} else {
		err = traiterCle(cleDeBase, source, &requeteInsertion, configExtraction, idMachine)
		if err != nil {
			probleme = errors.WithStack(err)
		}
	}
	err = requeteInsertion.Executer()
	if err != nil {
		return errors.WithStack(err)
	}
	return probleme
}

func (s Registre) CreationTable(cheminProjet string) error {
	aqua := aquabase.InitDB_Extraction(cheminProjet)
	err := aqua.CreateTableIfNotExist1("sam", colonnesTableSam, true)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func (s Registre) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if pourcentageChargement == -1 {
		bdd := aquabase.InitDB_Extraction(cheminProjet)
		if !bdd.EstTableVide("sam") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (s Registre) SQLChronologie() string {
	return "SELECT id, \"sam\", \"sam\", source, horodatage, \"opération \" || operation || \" effecutée sur le compte \" || nomCompte || \" (idCompte : \" || idCompte || \")\" FROM sam"
}
