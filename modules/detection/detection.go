package detection

import (
	"aquarium/modules/aquabase"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/* VARIABLES LOCALES */
var cheminRegles string = ""

type Regle struct {
	Nom         string    `json:"nom"`
	Auteur      string    `json:"auteur"`
	Description string    `json:"description"`
	Criticite   int       `json:"criticite"`
	Date        time.Time `json:"date"`
	SQL         string    `json:"sql"`
}

type regleSQL struct {
	Nom string `json:"nom"`
	SQL string `json:"sql"`
}

/* FONCTIONS LOCALES */

func lancerRegle(cheminProjet string, cheminRegle string) (int, error) {
	// On charge la requete SQL associée à la règle
	var detailsRegle regleSQL
	donneesFichier, err := os.ReadFile(cheminRegle)
	if err != nil {
		log.Println("WARN | Le fichier de règle n'existe pas ou n'a pas pu être ouvert : ", err.Error())
		return 0, err
	}
	err = json.Unmarshal(donneesFichier, &detailsRegle)
	if err != nil {
		return 0, err
	}
	log.Println(detailsRegle.SQL)
	// On exécute la requête SQL
	bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
	if err != nil {
		return 0, err
	}
	defer bd.Close()
	resultat, err := bd.Query(detailsRegle.SQL)
	if err != nil {
		return 0, nil
	}
	defer resultat.Close()
	if resultat.Next() {
		return 2, nil
	} else {
		return 1, nil
	}
}

func emplacementRegles() string {
	if cheminRegles != "" {
		return cheminRegles
	} else {
		emplacementExecutable, err := os.Executable()
		if err != nil {
			return ""
		}
		emplacementExecutable, err = filepath.EvalSymlinks(emplacementExecutable)
		if err != nil {
			return ""
		}
		// On cherche la liste des règles
		cheminRegles = filepath.Join(filepath.Dir(emplacementExecutable), "ressources", "regles_detection")
		return cheminRegles
	}
}

/* FONCTIONS GLOBALES */

/* Fonction qui renvoie une liste de règles associées à leur état (0:non lancé, 1:négatif, 2:positif) */
func ListeReglesDetection(cheminProjet string, lancerRegles bool) (map[string]int, []string, error) {
	var listeRegles map[string]int = map[string]int{}
	// On commence par chercher l'emplacement du logiciel aquarium
	fichiersRegles, err := os.ReadDir(emplacementRegles())
	if err != nil {
		return nil, nil, err
	}
	var probleme error = nil
	var reglesEnErreur []string = []string{}
	for _, fichierRegle := range fichiersRegles {
		var nomRegle string = strings.Replace(fichierRegle.Name(), ".json", "", 1)
		if lancerRegles {
			var cheminRegle string = filepath.Join(emplacementRegles(), fichierRegle.Name())
			listeRegles[nomRegle], err = lancerRegle(cheminProjet, cheminRegle)
			if err != nil {
				probleme = err
				listeRegles[nomRegle] = 0
				log.Println("detection.go => lancerRegle(", nomRegle, ") : ", err)
			} else if listeRegles[nomRegle] == 0 {
				reglesEnErreur = append(reglesEnErreur, nomRegle)
			}
		} else {
			listeRegles[nomRegle] = 0
		}
	}
	return listeRegles, reglesEnErreur, probleme
}

func ResultatRegleDetection(cheminProjet string, nomRegle string) (int, error) {
	return lancerRegle(cheminProjet, filepath.Join(emplacementRegles(), nomRegle+".json"))
}

func DetailsRegleDetection(cheminProjet string, nomRegle string) (Regle, error) {
	var donneesRegle Regle
	donneesFichier, err := os.ReadFile(filepath.Join(emplacementRegles(), nomRegle+".json"))
	if err != nil {
		log.Println("WARN | Le fichier de règle n'existe pas ou n'a pas pu être ouvert : ", err.Error())
		return Regle{}, err
	}
	err = json.Unmarshal(donneesFichier, &donneesRegle)
	return donneesRegle, err
}

func ResultatSQL(cheminProjet string, cheminRegle string, nomRegle string) ([]map[string]interface{}, error) {
	// On charge la requete SQL associée à la règle
	var detailsRegle regleSQL
	donneesFichier, err := os.ReadFile(cheminRegle)
	if err != nil {
		log.Println("WARN | Le fichier de règle "+cheminRegle+" n'existe pas ou n'a pas pu être ouvert : ", err.Error())
		return nil, err
	}
	err = json.Unmarshal(donneesFichier, &detailsRegle)
	if err != nil {
		return nil, err
	}
	log.Println(detailsRegle.SQL)

	// On exécute la requête SQL
	var adb = aquabase.InitBDDExtraction(cheminProjet)
	result := adb.SelectFrom(detailsRegle.SQL)
	fmt.Println(result)
	return result, nil
}
