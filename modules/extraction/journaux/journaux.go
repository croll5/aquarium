package journaux

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"aquarium/modules/extraction/utilitaires"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const SEPARATEUR_ENCODAGE = "|aqua_encodage:"
const INDICATEUR_CONCATENATION = "[aqua_concat]"

type Journaux struct{}

func (jr Journaux) Extraction(cheminProjet string, fichier bytes.Buffer, cheminFichierAExtraire string, configExtraction config.ConfigExtraction) error {
	var listeEvenements []string = []string{decoderFichier(fichier, configExtraction.Complement["encodage"])}
	if configExtraction.Complement["separateur"] != "" {
		listeEvenements = getListeDesEvenements(listeEvenements[0], configExtraction.Complement["separateur"])

	}
	var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	for _, table := range configExtraction.Table {
		var requeteInsertion aquabase.RequeteInsertion = abase.InitRequeteInsertionExtraction(table.Nom, table.GetNomsColonnes())
		for _, evenement := range listeEvenements {
			ajouterEvenementDansRequete(&requeteInsertion, cheminFichierAExtraire, evenement, configExtraction.Complement["symbole_association"], configExtraction.Complement["separateur_champs"], table, configExtraction)
		}
		requeteInsertion.Executer()
	}
	return nil
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
	if configExtraction.Complement["regex"] != "" {
		return extraireValeursRegex(evenement, reecritureSeparateur(configExtraction.Complement["regex"]))
	}
	separateur := reecritureSeparateur(configExtraction.Complement["separateur_champs"])
	symbAssociation := reecritureSeparateur(configExtraction.Complement["symbole_association"])
	guillemet := reecritureSeparateur(configExtraction.Complement["guillemets"])
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
	log.Println(regex)
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

func ajouterEvenementDansRequete(requeteInstertion *aquabase.RequeteInsertion, cheminFichierAExtraire string, evenement string, symboleAssociation string, separateurChamps string, configTable config.ConfigTableBDD, configExtraction config.ConfigExtraction) error {
	if evenement == "" {
		return nil
	}
	if configExtraction.Complement["exclusion"] != "" && strings.HasPrefix(evenement, configExtraction.Complement["exclusion"]) {
		return nil
	}
	dictChamps := extraireChampsEvenement(evenement, configExtraction)
	var valeursAAjouter []interface{} = make([]interface{}, 0)
	for _, colonne := range configTable.Colonnes {
		nomColonne := strings.Split(colonne.Contenu, SEPARATEUR_ENCODAGE)[0]
		if dictChamps[nomColonne] == "" {
			if nomColonne == "aqua_source" {
				valeursAAjouter = append(valeursAAjouter, cheminFichierAExtraire)
			} else if strings.Contains(colonne.Contenu, INDICATEUR_CONCATENATION) {
				valeur := ajouterConcatenation(colonne.Contenu, dictChamps)
				valeursAAjouter = append(valeursAAjouter, valeurDecodee(valeur, colonne.Contenu))
			} else {
				valeursAAjouter = append(valeursAAjouter, "[AQUA] Erreur : Clé introuvable")
			}
		} else {
			valeursAAjouter = append(valeursAAjouter, valeurDecodee(dictChamps[nomColonne], colonne.Contenu))
		}
	}
	requeteInstertion.AjouterDansRequete(valeursAAjouter...)
	return nil
}

func decoderFichier(fichier bytes.Buffer, encodage string) string {
	switch encodage {
	case "utf16":
		return utilitaires.Utf16LEToUtf8(fichier.String())
	}
	return fichier.String()
}

func valeurDecodee(valeur string, cle string) interface{} {
	if strings.Contains(cle, SEPARATEUR_ENCODAGE) {
		return utilitaires.DecoderString(valeur, strings.Split(cle, SEPARATEUR_ENCODAGE)[1])
	}
	return valeur
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
