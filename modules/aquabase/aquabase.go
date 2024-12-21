package aquabase

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-gota/gota/dataframe"
)

/** Structure jouant le role d'interface avec la database
 */
type Aquabase struct {
	dbPath string
}

type RequeteInsertion struct {
	nomTable string
	valeurs  []string
}

var basesDeDonnees map[string]*sql.DB = map[string]*sql.DB{}

/* -------------------------- GESTION DE LA BASE DE DONNÉES -------------------------- */

/*
	Fonction qui ouvre la base de données renseignée
	 - @param chemin : le chemin vers la base de données
	 - @return : un pointeur vers la base de données ouverte

Remarque : si la base est déjà ouverte, le programme revoie juste un pointeur vers celle-ci
*/
func GetBDD(chemin string) (*sql.DB, error) {
	if _, ok := basesDeDonnees[chemin]; ok {
		return basesDeDonnees[chemin], nil
	}
	bdd, err := sql.Open("sqlite", chemin)
	if err != nil {
		return nil, err
	}
	return bdd, nil
}

/*
		Fonction qui ouvre la base d’analyse

	  - @param cheminProjet : le chemin du dossier d’analyse
	  - @return : un pointeur vers la base ouverte, et s’il y a lieu une erreur
*/
func GetBaseExtraction(cheminProjet string) (*sql.DB, error) {
	return GetBDD(filepath.Join(cheminProjet, "analyse", "extractions.db"))
}

/* Fonction permettant de fermer une base de données */
func FermerBDD(cheminBDD string) error {
	err := basesDeDonnees[cheminBDD].Close()
	if err != nil {
		return err
	}
	delete(basesDeDonnees, cheminBDD)
	return nil
}

/* Fonction qui ferme toutes les bases de données ouvertes avec la fonction GetBDD */
func FermerToutesLesBDD() error {
	var probleme error = nil
	for cle := range basesDeDonnees {
		err := FermerBDD(cle)
		if err != nil {
			probleme = err
		}
	}
	return probleme
}

/*
* constructeur de la class Aquabase
  - @dbPath : path of the database
  - @return : the Aquabase object
  - Exemple : adb := aquabase.Init("C:\AquariumLab\analyse\extractions.db")
*/
func Init(dbPath string) Aquabase {
	a := Aquabase{}
	a.dbPath = dbPath
	return a
}

func InitBDDExtraction(cheminProjet string) Aquabase {
	a := Aquabase{}
	a.dbPath = filepath.Join(cheminProjet, "analyse", "extractions.db")
	return a
}

/** Create a connexion to the database.
 * After this functino: use defer db.Close() to close the connexion
 * @return : the Aquabase database connexion & an error if exist
 * Exemple : db, err := adb.Login()
 */
func (adb Aquabase) Login() (*sql.DB, error) {
	return GetBDD(adb.dbPath)
}

/* -------------------------- GESTION DES TABLES -------------------------- */

/** Create a table in the database
 * @tableName : name of the new table
 * @tableColumns : columns of the new table
 * @return : an error if exist
 * Exemple:  CreateTableIfNotExist( "getthis", ["colA", "ColB"])
 */
func (adb Aquabase) CreateTableIfNotExist(tableName string, tableColumns []string) error {
	// Open or create the sqliteDB
	db, err := adb.Login()
	if err != nil {
		return fmt.Errorf("CreateTableIfNotExist(): %w", err)
	}
	defer db.Close()
	// CREATE TABLE IF NOT EXISTS
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", tableName)
	for i, col := range tableColumns {
		query += fmt.Sprintf("%s TEXT", col)
		if i < len(tableColumns)-1 {
			query += ", "
		}
	}
	query += ");"
	// SQL code execution
	_, err = db.Exec(query)
	return err
}

/* -------------------------- SUPPRESSION DE VALEURS -------------------------- */

/** Delete values from a table with condition(s)
 * @tableName : name of the new table
 * @tableColumns : columns of the new table
 * @return : an error if exist
 * Exemple:  RemoveFromWhere( "getthis", ["colA", "ColB"])
 */
func (adb Aquabase) RemoveFromWhere(table string, where string) error {
	// Open sqliteDB
	db, err := adb.Login()
	if err != nil {
		return fmt.Errorf("SaveDf(): %w", err)
	}
	defer db.Close()
	// SLQ query
	queryDelete := fmt.Sprintf(`DELETE FROM '%s' WHERE %s`, table, where)
	_, err = db.Exec(queryDelete)
	if err != nil {
		return fmt.Errorf("Delete values failed: %w", err)
	}
	return nil
}

/* -------------------------- INSERTION DE VALEURS -------------------------- */

/** Insert into a table a dataframe
 * @df : the dataframe to save in the db
 * @tableName : the table were we save the data
 * @return : an error if exist
 * Exemple:  SaveDf(df, "getthis")
 */
func (adb Aquabase) SaveDf(df dataframe.DataFrame, tableName string) error {
	// Open sqliteDB
	db, err := adb.Login()
	if err != nil {
		return fmt.Errorf("SaveDf(): %w", err)
	}
	defer db.Close()
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ERROR: exportDfToDB() [Can't start a transaction]: %w", err)
	}
	// Prepare the query insertion
	queryAdd := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(df.Names(), ","),
		strings.Repeat("?,", len(df.Names())-1)+"?")
	stmt, err := tx.Prepare(queryAdd)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ERROR: exportDfToDB() [Can't Prepare query]: %w", err)
	}
	defer stmt.Close()
	// Add rows in the table
	for i := 0; i < df.Nrow(); i++ {
		values := make([]interface{}, df.Ncol())
		for j := 0; j < df.Ncol(); j++ {
			values[j] = df.Elem(i, j).String()
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ERROR: exportDfToDB() [Can't add data]: %w", err)
		}
	}
	// Commit the transaction
	err = tx.Commit()
	return err
}

/*
		Fonction d’initialisation d’une requête d’insertion dans la base extraction

	  - @param nomTable : le nom de la base dans laquelle il faut insérer les valeurs
	  - @return : un objet de type RequeteInsertion
*/
func InitRequeteInsertionExtraction(nomTable string) RequeteInsertion {
	var requete RequeteInsertion = RequeteInsertion{}
	requete.nomTable = nomTable
	requete.valeurs = []string{}
	return requete
}

func (requete *RequeteInsertion) AjouterDansRequete(valeurs []string) error {
	// On commence par nettoyer toutes les valeurs
	for i := range valeurs {
		valeurs[i] = "\"" + nettoyage(valeurs[i]) + "\""
	}
	// On en fait une unique chaîne de caractères
	requete.valeurs = append(requete.valeurs, "("+strings.Join(valeurs, ",")+")")
	return nil
}

func (requete *RequeteInsertion) Executer(cheminProjet string) error {
	if len(requete.valeurs) == 0 {
		log.Println("Il n’y avait aucun évènement !")
		return nil
	}
	var bdd *sql.DB
	bdd, err := GetBaseExtraction(cheminProjet)
	if err != nil {
		return err
	}
	var texteRequete string = "INSERT INTO " + requete.nomTable + " VALUES "
	texteRequete += strings.Join(requete.valeurs, ",")
	_, err = bdd.Exec(texteRequete)
	return err
}

/* -------------------------- LECTURE DE LA BASE -------------------------- */

/** Pragma request to obtains all the table name of the database
 * @return : dict of all table with the text "Columns: %d - Rows: %d"
 * Exemple: listTable := GetAllTableNames()
 */
func (adb Aquabase) GetAllTableNames() map[string]string {
	// Open sqliteDB
	db, err := adb.Login()
	if err != nil {
		return map[string]string{"Error": "Can't connect to database"}
	}
	defer db.Close()
	// Request the list of tables in the DB
	tables := make(map[string]string)
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		tables["Error"] = "GetAllTableNames(): Can't get tables list"
		return tables
	}
	defer rows.Close()
	// For each table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			tables["Error"] = "GetAllTableNames(): scanning table name"
			return tables
		}
		// Get number of columns
		var columnCount int
		columnQuery := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
		columnRows, err := db.Query(columnQuery)
		if err != nil {
			tables["Error"] = "GetAllTableNames(): querying column info for table" + tableName
			return tables
		}
		for columnRows.Next() {
			columnCount++
		}
		columnRows.Close()
		// Get number of rows
		var rowCount int
		rowQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
		err = db.QueryRow(rowQuery).Scan(&rowCount)
		if err != nil {
			tables["Error"] = "GetAllTableNames(): querying row count for table" + tableName
			return tables
		}
		tables[tableName] = fmt.Sprintf("Columns: %d - Rows: %d", columnCount, rowCount)
	}
	if err := rows.Err(); err != nil {
		tables["Error"] = "GetAllTableNames(): during rows iteration"
	}
	return tables
}

/** Simple SQL Selector with a limit size
 * @table : the table to select
 * @limit : the number of rows to select
 * @return : indexed dict contaning all rows data in a dict
 * Exemple: listTable := SelectAllFrom("getThis", 10)
 */
func (adb Aquabase) SelectAllFrom(table string, limit int) []map[string]interface{} {
	// Open sqliteDB
	db, err := adb.Login()
	if err != nil {
		return []map[string]interface{}{{"Error": "SelectAllFrom(): Can't connect to database"}}
	}
	defer db.Close()
	// SQL Request
	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, limit)
	rows, err := db.Query(query)
	if err != nil {
		return []map[string]interface{}{{"Error": "SelectAllFrom(): querying table data"}}
	}
	defer rows.Close()
	// Take columns
	columns, err := rows.Columns()
	if err != nil {
		return []map[string]interface{}{{"Error": "SelectAllFrom(): getting columns"}}
	}
	// Create the dataframe
	var results []map[string]interface{}
	for rows.Next() {
		columnPointers := make([]interface{}, len(columns))
		columnValues := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return []map[string]interface{}{{"Error": "SelectAllFrom(): scanning row: " + table}}
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			rowMap[colName] = columnValues[i]
		}
		results = append(results, rowMap)
	}
	if err := rows.Err(); err != nil {
		return []map[string]interface{}{{"Error": "SelectAllFrom(): during rows iteration: " + table}}
	}
	if len(results) == 0 {
		return []map[string]interface{}{{"Info": "no value in: " + table}}
	}
	return results
}

func (adb Aquabase) EstTableVide(table string) bool {
	bdd, err := GetBDD(adb.dbPath)
	if err != nil {
		return true
	}
	resultat, err := bdd.Query("SELECT * FROM " + table + " LIMIT 1")
	if err != nil {
		return false
	}
	return !resultat.Next()
}

/* -------------------------- FONCTIONS ANNEXES -------------------------- */

func nettoyage(entree string) string {
	entree = strings.ReplaceAll(entree, "\"", "&quot;")
	entree = strings.ReplaceAll(entree, "<", "&lt;")
	entree = strings.ReplaceAll(entree, ">", "&gt;")
	entree = strings.ReplaceAll(entree, "&", "&amp")
	return entree
}
