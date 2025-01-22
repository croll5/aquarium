/*
Copyright ou © ou Copr. Cécile Rolland et Charles Mailley, (21 janvier 2025)

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
	ruleName := strings.Replace(filepath.Base(cheminRegle), ".json", "", -1)

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
		id_value := id_frame.Strloc(0, 0)

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

/*
*
return the path of the rule if it exist and if its a global rule
*/
func getRulePathFolder(cheminProjet string, nomRegle string) (string, bool) {
	// Search where the rule is saved in local
	path_local := filepath.Join(cheminProjet, "regles_detection")
	existInLocal, _ := searchDetectionRules(path_local, nomRegle)
	if len(existInLocal) == 1 {
		return path_local, false
	}
	// Search where the rule is saved in global
	path_global := emplacementRegles()
	existInGlobal, _ := searchDetectionRules(path_global, nomRegle)
	if len(existInGlobal) == 1 {
		return path_global, true
	}
	// The rule doesnt exist
	return "", false
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
		} else {
			// Cherche si la regle a deja été executé avant
			adb_rules := aquabase.InitDB_Rules(cheminProjet)
			query := fmt.Sprintf("SELECT isError FROM regles WHERE name=\"%s\"", rule)
			df := adb_rules.SelectFrom0(query)
			if df.Table.Nrow() > 0 {
				value, _ := df.Intloc(0, 0)
				if value == 1 {
					state = 2
				} else {
					state = 1
				}
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
	path, isGlobal := getRulePathFolder(cheminProjet, nomRegle)
	donneesRegle.IsGlobal = isGlobal
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
	path, _ := getRulePathFolder(cheminProjet, nomRegle)
	// Execute the SQL request
	return lancerRegle(cheminProjet, filepath.Join(path, nomRegle+".json"))
}

func ResultatSQL(cheminProjet string, ruleName string) ([]map[string]interface{}, error) {
	adb_rules := aquabase.InitDB_Rules(cheminProjet)
	id_frame := adb_rules.SelectFrom0("SELECT id FROM regles WHERE name='" + ruleName + "'")
	id_value := id_frame.Strloc(0, 0)
	df := adb_rules.SelectFrom0("SELECT * FROM error_" + id_value)
	return df.ToMap(), nil
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
	// Delete all data about this rule from regles.db
	adb_rules := aquabase.InitDB_Rules(cheminProjet)
	err = adb_rules.DropTable(nomRegle)
	if err != nil {
		return err
	}
	err = adb_rules.RemoveFromWhere("regles", "name='"+nomRegle+"'")
	if err != nil {
		return err
	}
	return nil
}

func StatutReglesDetection(cheminProjet string) ([]map[string]interface{}, error) {
	adb_rules := aquabase.InitDB_Rules(cheminProjet)
	df := adb_rules.SelectFrom0("SELECT * FROM regles")
	if df.Error != nil {
		fmt.Println("Table 'regle' inexistante ou erreur")
		return []map[string]interface{}{}, nil
	}
	return df.ToMap(), nil
}

func NewDetectionRule(chemin_projet string, json_rule string) error {
	chemin_regles := filepath.Join(chemin_projet, "regles_detection")
	// Conversion de la chaîne JSON en une structure Go
	var regle map[string]interface{}
	if err := json.Unmarshal([]byte(json_rule), &regle); err != nil {
		return err
	}
	// Récupération du nom à partir du JSON
	nom, ok := regle["nom"].(string)
	if !ok {
		return fmt.Errorf("Json without the variable: nom")
	}
	nameBeforeModification, ok := regle["nameBeforeModification"].(string)
	if !ok {
		return fmt.Errorf("Json without the variable: nameBeforeModification")
	}
	//Verification que la regle n'existe pas déjà
	rulePathFolder, _ := getRulePathFolder(chemin_projet, nom)
	if len(rulePathFolder) != 0 && nom != nameBeforeModification {
		return fmt.Errorf("The name '" + nom + "' is already used")
	}
	// Suppression du champ json nameBeforeModification et du json avec l'ancien nom
	if nameBeforeModification != "" {
		err := SuppressionRegleDetection(chemin_projet, nameBeforeModification)
		if err != nil {
			return err
		}
	}
	delete(regle, "nameBeforeModification")
	// Conversion de la structure Go en JSON formaté
	data, err := json.MarshalIndent(regle, "", "  ")
	if err != nil {
		return err
	}
	// Création du chemin complet du fichier avec le nom du JSON
	chemin_complet := filepath.Join(chemin_regles, nom+".json")
	// Écriture des données JSON dans un fichier
	if err := os.WriteFile(chemin_complet, data, 0644); err != nil {
		return err
	}
	return nil
}
