package detection

import (
	"aquarium/modules/aquabase"
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
	IsGlobal    bool
}

type regleSQL struct {
	Nom string `json:"nom"`
	SQL string `json:"sql"`
}

/* FONCTIONS LOCALES */

func lancerRegle(cheminProjet string, cheminRegle string) (int, error) {
	// Recuperation de la requete SQL associée à la règle
	var detailsRegle regleSQL
	donneesFichier, err := os.ReadFile(cheminRegle)
	if err != nil {
		log.Println("WARN | Le fichier de règle "+cheminRegle+" n'existe pas ou n'a pas pu être ouvert : ", err.Error())
		return 0, err
	}
	err = json.Unmarshal(donneesFichier, &detailsRegle)
	if err != nil {
		return 0, err
	}
	ruleName := filepath.Base(cheminRegle)

	// Execution la requête SQL
	var adb = aquabase.InitDB_Extraction(cheminProjet)
	df := adb.SelectFrom0(detailsRegle.SQL)
	isError := df.Table.Nrow() > 0

	// Renseignement de la table sql des regles
	adb_rules := aquabase.InitDB_Rules(cheminProjet)
	tableName := "regles"
	colName := []string{"name", "isError"}
	tableColumns := map[string]string{
		// "id" en autoincrement par defaut
		"name":    "TEXT UNIQUE",
		"isError": "INTEGER",
	}
	err = adb_rules.CreateTableIfNotExist2(tableName, tableColumns, true)
	if err != nil {
		return 0, err
	}
	err = adb_rules.InsertOrReplace(tableName, colName, []interface{}{ruleName, isError})
	if err != nil {
		return 0, err
	}
	//fmt.Println(adb_rules.SelectFrom0("SELECT * FROM regles"))

	// Renvoi 2 si le dataframe n'est pas vide sinon 1
	if isError {
		id_frame := adb_rules.SelectFrom0("SELECT id FROM regles WHERE name='" + ruleName + "'")
		id_value := id_frame.Iloc(0, 0)
		fmt.Println(id_value)

		table_name := "error_" + id_value
		err := adb_rules.DropTable(table_name)
		if err != nil {
			return 0, err
		}
		err = adb_rules.CreateTableIfNotExist1(table_name, df.Table.Names(), false)
		if err != nil {
			return 0, err
		}
		err = adb_rules.SaveDf(df.Table, table_name)
		if err != nil {
			return 0, err
		}
		return 2, nil
	}
	return 1, nil
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

/** Extract all rule names from a folder
 * @rulesPath : the path folder where rules.json exists
 * @return : list of name rules
 */
func searchDetectionRules(rulesPath string, args ...string) ([]string, error) {
	var listeRegles []string
	// Vérifier si une clé est fournie
	var key string
	if len(args) > 0 {
		key = args[0]
	}
	// Open in kernel the folder
	fichiersRegles, err := os.ReadDir(rulesPath)
	if err != nil {
		return nil, err
	}
	// Catch all json file in a list
	for _, fichierRegle := range fichiersRegles {
		if strings.HasSuffix(fichierRegle.Name(), ".json") {
			nomRegle := strings.TrimSuffix(fichierRegle.Name(), ".json")
			if key == "" || strings.EqualFold(nomRegle, key) {
				listeRegles = append(listeRegles, nomRegle)
			}
		}
	}
	return listeRegles, nil
}

/********************************************************************************/
/****************************** FONCTIONS GLOBALES ******************************/
/********************************************************************************/

/* Fonction qui renvoie une liste de règles associées à leur état (0:non lancé, 1:négatif, 2:positif) */
func ListeReglesDetection(cheminProjet string, lancerRegles bool) (map[string]map[string]int, []string, error) {
	listeRegles := make(map[string]map[string]int)
	// Search all rule files
	path_local := filepath.Join(cheminProjet, "regles_detection")
	path_global := emplacementRegles()
	listeRegles_local, error_local := searchDetectionRules(path_local)
	listeRegles_global, error_global := searchDetectionRules(path_global)
	if error_local != nil {
		return nil, nil, error_local
	}
	if error_global != nil {
		return nil, nil, error_global
	}
	// Helper function to handle the rule logic
	var probleme error = nil
	var reglesEnErreur []string = []string{}
	handleRule := func(rule string, isGlobal int, path string) {
		state := 0
		var err error
		if lancerRegles {
			path_rule := filepath.Join(path, rule+".json")
			state, err = lancerRegle(cheminProjet, path_rule)
			if err != nil {
				state = 0
				probleme = err
				log.Println("detection.go => lancerRegle(", rule, ") : ", err)
			} else if state == 0 {
				reglesEnErreur = append(reglesEnErreur, rule)
			}
		}
		listeRegles[rule] = map[string]int{
			"isGlobal": isGlobal,
			"state":    state,
		}
	}
	// Merge both list in a list of dict with parameters of each rule
	// Une regle créé par l'user et prioritaire par rapport à une regle de base
	for _, rule := range listeRegles_global {
		handleRule(rule, 1, path_global)
	}
	for _, rule := range listeRegles_local {
		handleRule(rule, 0, path_local)
	}
	return listeRegles, reglesEnErreur, probleme
}

func DetailsRegleDetection(cheminProjet string, nomRegle string) (Regle, error) {
	var donneesRegle Regle
	// Search where the rule is saved
	path_local := filepath.Join(cheminProjet, "regles_detection")
	path_global := emplacementRegles()
	var path string
	exist, _ := searchDetectionRules(path_local, nomRegle)
	if len(exist) == 1 {
		path = path_local
	} else {
		path = path_global
		donneesRegle.IsGlobal = true
	}
	// Read and extract the rule.json data
	donneesFichier, err := os.ReadFile(filepath.Join(path, nomRegle+".json"))
	if err != nil {
		log.Println("WARN DetailsRegleDetection() | Le fichier de règle n'existe pas ou n'a pas pu être ouvert : ", err.Error())
		return Regle{}, err
	}
	err = json.Unmarshal(donneesFichier, &donneesRegle)
	return donneesRegle, err
}

func ResultatRegleDetection(cheminProjet string, nomRegle string) (int, error) {
	// Search where the rule is saved
	path_local := filepath.Join(cheminProjet, "regles_detection")
	path_global := emplacementRegles()
	var path string
	exist, _ := searchDetectionRules(path_local, nomRegle)
	if len(exist) == 1 {
		path = path_local
	} else {
		path = path_global
	}
	// Execute the SQL request
	return lancerRegle(cheminProjet, filepath.Join(path, nomRegle+".json"))
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
	var adb = aquabase.InitDB_Extraction(cheminProjet)
	result := adb.SelectFrom(detailsRegle.SQL)
	fmt.Println(result)
	return result, nil
}

func SuppressionRegleDetection(cheminProjet string, nomRegle string) error {
	// Search where the rule is saved
	path_local := filepath.Join(cheminProjet, "regles_detection")
	exist, _ := searchDetectionRules(path_local, nomRegle)
	if len(exist) != 1 {
		fmt.Println("Annulation de suppression de la regle: " + nomRegle)
		return nil
	}
	// Delete the rule.json data
	err := os.Remove(filepath.Join(path_local, nomRegle+".json"))
	if err != nil {
		log.Println("WARN | Le fichier de règle n'a pas pu être supprimé : ", err.Error())
		return err
	}
	fmt.Println("Suppression de la regle: " + nomRegle)
	return nil
}
