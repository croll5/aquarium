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

package journaux

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"aquarium/modules/extraction/utilitaires"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	SEPARATEUR_ENCODAGE       = "|aqua_encodage:"
	INDICATEUR_CONCATENATION  = "[aqua_concat]"
	PARAM_SEPARATEUR          = "separateur"
	PARAM_ENCODAGE            = "encodage"
	PARAM_SEPARATEUR_CHAMPS   = "separateur_champs"
	PARAM_SYMBOLE_ASSOCIATION = "symbole_association"
	PARAM_GUILLEMENTS         = "guillemets"
	PARAM_REGEX               = "regex"
	PARAM_EXCLUSION           = "exclusion"
)

var fonctionsTraitementContenuColonne map[string]func(dicChamps map[string]string, cheminFichier string, idMachine string) interface{} = map[string]func(map[string]string, string, string) interface{}{}

type Journaux struct{}

func (jr Journaux) Extraction(cheminProjet string, fichier io.Reader, cheminFichierAExtraire string, configExtraction config.ConfigExtraction, idMachine string) error {
	var listeEvenements []string = []string{decoderFichier(fichier, configExtraction.Complement[PARAM_ENCODAGE])}
	if configExtraction.Complement[PARAM_SEPARATEUR] != "" {
		listeEvenements = getListeDesEvenements(listeEvenements[0], configExtraction.Complement[PARAM_SEPARATEUR])
	}
	var probleme error
	for _, table := range configExtraction.Table {
		var requeteInsertion aquabase.RequeteInsertion = *utilitaires.CreerRequeteInstertionDepuisConfig(cheminProjet, idMachine, &table)
		for _, evenement := range listeEvenements {
			ajouterEvenementDansRequete(&requeteInsertion, cheminFichierAExtraire, evenement, table, configExtraction, idMachine)
		}
		err := requeteInsertion.Executer()
		if err != nil {
			probleme = errors.WithStack(err)
		}
	}
	return probleme
}

/* -------------------------------- FONCTIONS LOCALES --------------------------------*/

func valeurDecodee(contenuColonne string, dicChamps map[string]string, cheminFichier string, idMachine string) interface{} {
	// Si la fonction existe déjà, on l’utilise
	fonctionTraitement, existe := fonctionsTraitementContenuColonne[contenuColonne]
	if existe {
		return fonctionTraitement(dicChamps, cheminFichier, idMachine)
	}
	log.Println("Extraction du contenu ", contenuColonne)
	// On commence par séparer la clé de l’encodage
	parametresColonne := strings.Split(contenuColonne, SEPARATEUR_ENCODAGE)
	// S’il n’y a pas d’encodage, on met « string »
	if len(parametresColonne) < 2 {
		parametresColonne = []string{parametresColonne[0], "string"}
	}
	// Et on regarde à quoi correspond la clé
	switch contenuColonne {
	case "aqua_source":
		fonctionTraitement = func(champs map[string]string, cheminFichier, idMachine string) interface{} {
			return cheminFichier
		}
	case config.AQUA_MACHINE:
		fonctionTraitement = func(champs map[string]string, cheminFichier, idMachine string) interface{} {
			return idMachine
		}
	default:
		fonctionDecodage := utilitaires.GetFonctionDecodageString(parametresColonne[1])
		if strings.Contains(parametresColonne[0], INDICATEUR_CONCATENATION) {
			cles := strings.Split(parametresColonne[0], INDICATEUR_CONCATENATION)
			fonctionTraitement = func(dicChamps map[string]string, cheminFichier, idMachine string) interface{} {
				resultat := ""
				for i := 0; i < len(cles)-1; i++ {
					resultat = resultat + dicChamps[cles[i]]
					if i < len(cles)-2 {
						resultat += cles[len(cles)-1]
					}
				}
				return fonctionDecodage(resultat)
			}
		} else {
			fonctionTraitement = func(dicChamps map[string]string, cheminFichier, idMachine string) interface{} {
				val, ok := dicChamps[parametresColonne[0]]
				if !ok {
					return nil
				}
				return val
			}
		}
	}
	// On choisit la fonction utilisée en fonction de l’encodage
	fonctionsTraitementContenuColonne[contenuColonne] = fonctionTraitement
	return fonctionTraitement(dicChamps, cheminFichier, idMachine)
}

func getListeDesEvenements(fichier string, separateur string) []string {
	separateur = reecritureSeparateur(separateur)
	return strings.Split(fichier, separateur)
}

func reecritureSeparateur(separateur string) string {
	separateur = strings.ReplaceAll(separateur, "\\n", "\n")
	separateur = strings.ReplaceAll(separateur, "\\r", "\r")
	separateur = strings.ReplaceAll(separateur, "&lt;", "<")
	return separateur
}

func extraireChampsEvenement(evenement string, configExtraction config.ConfigExtraction) map[string]string {
	var resultat map[string]string = map[string]string{}
	// Traiter le cas de l’utilisation d’une expression régulière
	if configExtraction.Complement[PARAM_REGEX] != "" {
		return extraireValeursRegex(evenement, reecritureSeparateur(configExtraction.Complement[PARAM_REGEX]))
	}
	separateur := reecritureSeparateur(configExtraction.Complement[PARAM_SEPARATEUR_CHAMPS])
	symbAssociation := reecritureSeparateur(configExtraction.Complement[PARAM_SYMBOLE_ASSOCIATION])
	guillemet := reecritureSeparateur(configExtraction.Complement[PARAM_GUILLEMENTS])
	tailleSymAssociation := len(symbAssociation)
	tailleSeparateur := len(separateur)
	tailleGuillements := len(guillemet)
	indexDebutCle := 0
	cle := "0"
	compteurCles := 0
	indexDebutValeur := 0
	estDansCle := tailleSymAssociation > 0
	finitParGuillemets := false
	tailleEvenement := len(evenement)
	for i := 0; i < tailleEvenement; i++ {
		if estDansCle {
			if i < (tailleEvenement-tailleSymAssociation) && evenement[i:i+tailleSymAssociation] == symbAssociation {
				cle = evenement[indexDebutCle:i]
				i = i + tailleSymAssociation
				estDansCle = false
				// On regarde s'il y a des guillemets
				if tailleGuillements > 0 && i < (tailleEvenement-tailleGuillements) && evenement[i:i+tailleGuillements] == guillemet {
					i = i + tailleGuillements
					finitParGuillemets = true
				}
				// On définit l’index du début de la valeur
				indexDebutValeur = i
			}
		} else {
			if finitParGuillemets {
				// On regarde si les caractères suivants sont le guillemet
				if i < (tailleEvenement-tailleGuillements) && evenement[i:i+tailleGuillements] == guillemet {
					resultat[cle] = evenement[indexDebutValeur:i]
					i = i + tailleGuillements
					for i < (tailleEvenement-tailleSeparateur) && evenement[i:i+tailleSeparateur] == separateur {
						i = i + tailleSeparateur
					}
					if tailleSymAssociation > 0 {
						estDansCle = true
						indexDebutCle = i
					} else {
						compteurCles++
						cle = fmt.Sprint(compteurCles)
						indexDebutValeur = i
					}
					finitParGuillemets = false
				} else if i == (tailleEvenement - tailleGuillements) {
					resultat[cle] = evenement[indexDebutValeur : i+1]
				}
			} else {
				if i < (tailleEvenement-tailleSeparateur) && evenement[i:i+tailleSeparateur] == separateur {
					resultat[cle] = evenement[indexDebutValeur:i]
					i = i + tailleSeparateur
					for i < (tailleEvenement-tailleSeparateur) && evenement[i:i+tailleSeparateur] == separateur {
						i = i + tailleSeparateur
					}
					if tailleSymAssociation > 0 {
						estDansCle = true
						indexDebutCle = i
					} else {
						compteurCles++
						cle = fmt.Sprint(compteurCles)
						indexDebutValeur = i
					}
				} else if i == tailleEvenement-1 {
					resultat[cle] = evenement[indexDebutValeur : i+1]
				}
			}
		}
	}
	return resultat
}

func extraireValeursRegex(donnees string, regex string) map[string]string {
	var resultat map[string]string = map[string]string{}
	rex, err := regexp.Compile(regex)
	if err != nil {
		resultat["0"] = "[AQUA] Erreur dans l’extraction de l’expression régulière : " + err.Error()
		resultat["1"] = "[AQUA] Erreur dans l’extraction de l’expression régulière : " + err.Error()
		return resultat
	}
	valeurs := rex.FindStringSubmatch(donnees)
	for i, valeur := range valeurs {
		resultat[fmt.Sprint(i)] = valeur
	}
	return resultat
}

func ajouterEvenementDansRequete(requeteInstertion *aquabase.RequeteInsertion, cheminFichierAExtraire string, evenement string, configTable config.ConfigTableBDD, configExtraction config.ConfigExtraction, idMachine string) error {
	if evenement == "" {
		return nil
	}
	if configExtraction.Complement[PARAM_EXCLUSION] != "" && strings.HasPrefix(evenement, configExtraction.Complement[PARAM_EXCLUSION]) {
		return nil
	}
	dictChamps := extraireChampsEvenement(evenement, configExtraction)
	var valeursAAjouter []interface{} = make([]interface{}, 0)
	for _, colonne := range configTable.Colonnes {
		valeursAAjouter = append(valeursAAjouter, valeurDecodee(colonne.Contenu, dictChamps, cheminFichierAExtraire, idMachine))
	}
	err := requeteInstertion.AjouterDansRequete(valeursAAjouter...)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func decoderFichier(fichier io.Reader, encodage string) string {
	// Copie du contenu du fichier dans un tampon, pour pouvoir l'ouvrir avec l'extracteur de registres
	var tampon bytes.Buffer
	if _, err := io.Copy(&tampon, fichier); err != nil {
		log.Println("Format de fichier non supporté : ", err.Error())
	}
	switch encodage {
	case "utf16":
		return utilitaires.Utf16LEToUtf8(tampon.String())
	}
	return tampon.String()
}

func ajouterConcatenation(cle string, dictChamps map[string]string) string {
	cle = strings.Split(cle, SEPARATEUR_ENCODAGE)[0]
	valeurs := strings.Split(cle, INDICATEUR_CONCATENATION)
	resultat := ""
	for i := 0; i < len(valeurs)-1; i++ {
		resultat = resultat + dictChamps[valeurs[i]]
		if i < len(valeurs)-2 {
			resultat += valeurs[len(valeurs)-1]
		}
	}
	return resultat
}
