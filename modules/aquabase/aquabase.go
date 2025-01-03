package aquabase

import (
	"aquarium/modules/aquaticket"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gota/gota/dataframe"
)

/** Structure jouant le role d'interface avec la database
 */
type Aquabase struct {
	dbPath string
	dbName string
}

type RequeteInsertion struct {
	nomTable      string
	colonnesTable []string
	valeurs       [][]interface{}
}

type infosBDD struct {
	bdd     *sql.DB
	tickets *aquaticket.Distributeur
}

var basesDeDonnees map[string]infosBDD = map[string]infosBDD{}

/* -------------------------- GESTION DE LA BASE DE DONNÉES -------------------------- */

/*
	Fonction qui ouvre la base de données renseignée
	 - @param chemin : le chemin vers la base de données
	 - @return : un pointeur vers la base de données ouverte

Remarque : si la base est déjà ouverte, le programme revoie juste un pointeur vers celle-ci
*/
func GetInfosBDD(chemin string) (infosBDD, error) {
	if _, ok := basesDeDonnees[chemin]; ok {
		return basesDeDonnees[chemin], nil
	}
	bdd, err := sql.Open("sqlite", chemin)
	distributeur := aquaticket.NouveauDistributeur()
	basesDeDonnees[chemin] = infosBDD{bdd: bdd, tickets: &distributeur}
	if err != nil {
		return infosBDD{}, err
	}
	return basesDeDonnees[chemin], nil
}

/*
		Fonction qui ouvre la base d’analyse

	  - @param cheminProjet : le chemin du dossier d’analyse
	  - @return : un pointeur vers la base ouverte, et s’il y a lieu une erreur
*/
func GetInfosBaseExtraction(cheminProjet string) (infosBDD, error) {
	return GetInfosBDD(filepath.Join(cheminProjet, "analyse", "extractions.db"))
}

/* Fonction permettant de fermer une base de données */
func FermerBDD(cheminBDD string) error {
	err := basesDeDonnees[cheminBDD].bdd.Close()
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
			log.Println("[ERROR] Erreur dans la fermeture de la table ", cle, " : ", err)
			probleme = err
		} else {
			log.Println("[INFO] Table ", cle, " fermée avec succès.")
		}
	}
	return probleme
}

/** Create a connexion to the database.
 * After this functino: use defer db.Close() to close the connexion
 * @return : the Aquabase database connexion & an error if exist
 * Exemple : db, err := adb.Login()
 */
func (adb Aquabase) Login() (infosBDD, error) {
	return GetInfosBDD(adb.dbPath)
}

/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ----------------------------------------   INITIALISATION   ---------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */

/*
constructeur de la class Aquabase
  - @dbPath : path of the database
  - @return : the Aquabase object pointers
  - Exemple : adb := aquabase.Init("C:\AquariumLab\analyse\extractions.db")
*/
func Init(dbPath string) *Aquabase {
	adb := Aquabase{}
	adb.dbPath = dbPath
	adb.dbName = filepath.Base(dbPath)
	if !adb.createDatabaseIfNotExist() {
		return nil
	}
	return &adb
}

/*
constructeur de la class Aquabase avec une base de données personnalisée
  - @projectPath : path of the project
  - @return : the Aquabase object pointers
  - Exemple : adb := aquabase.Init("C:\AquariumLab")
*/
func InitDB_Extraction(projectPath string) *Aquabase {
	adb := Aquabase{}
	adb.dbName = "extractions.db"
	adb.dbPath = filepath.Join(projectPath, "analyse", adb.dbName)
	if !adb.createDatabaseIfNotExist() {
		return nil
	}
	return &adb
}

/*
constructeur de la class Aquabase avec une base de données personnalisée
  - @projectPath : path of the project
  - @return : the Aquabase object pointers
  - Exemple : adb := aquabase.Init("C:\AquariumLab")
*/
func InitDB_Rules(projectPath string) *Aquabase {
	adb := Aquabase{}
	adb.dbName = "regles.db"
	adb.dbPath = filepath.Join(projectPath, "analyse", adb.dbName)
	if !adb.createDatabaseIfNotExist() {
		return nil
	}
	return &adb
}

/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ----------------------------------------       CREATE       ---------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */

/*
Creation du fichier de la base de données
  - @return default: true. false if the file.db cannot be create
*/
func (adb Aquabase) createDatabaseIfNotExist() bool {
	// Vérifie si le fichier existe déjà
	if _, err := os.Stat(adb.dbPath); err == nil {
		return true
	}
	abd_file, err := os.Create(adb.dbPath)
	if err != nil {
		return false
		log.Println("adb.WARNING - can't create the database file: " + adb.dbPath)
	}
	defer abd_file.Close()
	fmt.Println("Creation de la dbb: " + adb.dbPath)
	return true
}

/*
Create a table in the database
  - @tableName : name of the new table
  - @tableColumns : columns of the new table
  - @return : an error if exist
  - Exemple:  CreateTableIfNotExist( "getthis", ["colA", "ColB"])
*/
func (adb Aquabase) CreateTableIfNotExist(tableName string, tableColumns []string) error {
	// Open or create the sqliteDB
	infosBdd, err := adb.Login()
	if err != nil {
		return fmt.Errorf("CreateTableIfNotExist(): %w", err)
	}
	// Check the table existance
	var name string
	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", tableName)
	res := infosBdd.bdd.QueryRow(query).Scan(&name)
	if res == nil {
		return nil
	}
	// CREATE TABLE IF NOT EXISTS
	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, ", tableName)
	for i, col := range tableColumns {
		query += fmt.Sprintf("%s TEXT", col)
		if i < len(tableColumns)-1 {
			query += ", "
		}
	}
	query += ");"
	// SQL code execution
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		_, err := infosBdd.bdd.Exec(query)
		return err
	})
	fmt.Println("Create table '" + tableName + "' in " + adb.dbName)
	return err
}

/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ----------------------------------------       DELETE       ---------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */

/*
Delete a table from the database
  - @tableName : name of the new table
*/
func (adb Aquabase) DropTable(table string) error {
	// Open sqliteDB
	infosBdd, err := adb.Login()
	if err != nil {
		fmt.Println("adb.WARNING: DropTable failed: " + table)
		return err
	}
	// Drop the table
	queryDrop := fmt.Sprintf(`DROP TABLE IF EXISTS '%s'`, table)
	infosBdd.bdd.Exec(queryDrop)
	return nil
}

/*
Delete values from a table with condition(s)
  - @tableName : name of the new table
  - @tableColumns : columns of the new table
  - @return : an error if exist
  - Exemple:  RemoveFromWhere( "getthis", ["colA", "ColB"])
*/
func (adb Aquabase) RemoveFromWhere(table string, where string) error {
	// Open sqliteDB
	infosBdd, err := adb.Login()
	if err != nil {
		return fmt.Errorf("SaveDf(): %w", err)
	}
	// SLQ query
	queryDelete := fmt.Sprintf(`DELETE FROM '%s' WHERE %s`, table, where)
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		_, err = infosBdd.bdd.Exec(queryDelete)
		return err
	})
	if err != nil {
		return fmt.Errorf("Delete values failed: %w", err)
	}
	return nil
}

/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ----------------------------------------       INSERT       ---------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */

/*
Insert into a table a dataframe
  - @df : the dataframe to save in the db
  - @tableName : the table were we save the data
  - @return : an error if exist
  - Exemple:  SaveDf(df, "getthis")
*/
func (adb Aquabase) SaveDf(df dataframe.DataFrame, tableName string) error {
	// Open sqliteDB
	infosBdd, err := adb.Login()
	if err != nil {
		return fmt.Errorf("SaveDf(): %w", err)
	}
	// Start a transaction
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		tx, err := infosBdd.bdd.Begin()
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
	})
	return err
}

/*
Fonction d’initialisation d’une requête d’insertion dans la base extraction
  - @param nomTable : le nom de la base dans laquelle il faut insérer les valeurs
  - @return : un objet de type RequeteInsertion
*/
func InitRequeteInsertionExtraction(nomTable string, colonnesTable []string) RequeteInsertion {
	var requete RequeteInsertion = RequeteInsertion{}
	requete.nomTable = nomTable
	requete.colonnesTable = colonnesTable
	requete.valeurs = make([][]interface{}, 0)
	return requete
}

func (requete *RequeteInsertion) AjouterDansRequete(valeurs ...any) error {
	if len(valeurs) != len(requete.colonnesTable) {
		return errors.New("Mauvais nombre de colonnes")
	}
	// On en fait une unique chaîne de caractères
	requete.valeurs = append(requete.valeurs, valeurs)
	return nil
}

func (requete *RequeteInsertion) Executer(cheminProjet string) error {
	if len(requete.valeurs) == 0 {
		log.Println("Il n’y avait aucun évènement !")
		return nil
	}
	infosBdd, err := GetInfosBaseExtraction(cheminProjet)
	if err != nil {
		return err
	}
	// Préparation des instesions
	var texteRequete string = "INSERT INTO " + requete.nomTable + "("
	texteRequete += strings.Join(requete.colonnesTable, ",")
	texteRequete += ") VALUES (" + strings.Repeat("?,", len(requete.colonnesTable)-1) + "?)"
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		// Création de la transaction
		tx, err := infosBdd.bdd.Begin()
		if err != nil {
			return fmt.Errorf("ERROR: requete.Executer() impossible de creer la transaction : %w", err)
		}
		// Prepare the query insertion

		stmt, err := tx.Prepare(texteRequete)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("ERROR: requete.Executer() [Can't Prepare query]: %w", err)
		}
		defer stmt.Close()
		// Add rows in the table
		for _, ligne := range requete.valeurs {
			_, err = stmt.Exec(ligne...)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("ERROR: exportDfToDB() [Can't add data]: %w", err)
			}
		}
		// Commit the transaction
		err = tx.Commit()
		return err
	})
	return err
}

/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ----------------------------------------       SELECT       ---------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */
/* ---------------------------------------------------------------------------------------------------- */

/** Pragma request to obtains all the table name of the database
 * @return : dict of all table with the text "Columns: %d - Rows: %d"
 * Exemple: listTable := GetAllTableNames()
 */
func (adb Aquabase) GetAllTableNames() map[string]string {
	// Open sqliteDB
	infosBdd, err := adb.Login()
	if err != nil {
		return map[string]string{"Error": "Can't connect to database"}
	}
	// Request the list of tables in the DB
	tables := make(map[string]string)
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		rows, err := infosBdd.bdd.Query("SELECT name FROM sqlite_master WHERE type='table'")
		if err != nil {
			tables["Error"] = "GetAllTableNames(): Can't get tables list"
			return err
		}
		defer rows.Close()
		// For each table
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				tables["Error"] = "GetAllTableNames(): scanning table name"
				return err
			}
			// Get number of columns
			var columnCount int
			columnQuery := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
			columnRows, err := infosBdd.bdd.Query(columnQuery)
			if err != nil {
				tables["Error"] = "GetAllTableNames(): querying column info for table" + tableName
				return err
			}
			for columnRows.Next() {
				columnCount++
			}
			columnRows.Close()
			// Get number of rows
			var rowCount int
			rowQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
			err = infosBdd.bdd.QueryRow(rowQuery).Scan(&rowCount)
			if err != nil {
				tables["Error"] = "GetAllTableNames(): querying row count for table" + tableName
				return err
			}
			tables[tableName] = fmt.Sprintf("Columns: %d - Rows: %d", columnCount, rowCount)
		}
		if err := rows.Err(); err != nil {
			tables["Error"] = "GetAllTableNames(): during rows iteration"
		}
		return err
	})
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
	infosBdd, err := adb.Login()
	if err != nil {
		return []map[string]interface{}{{"Error": "SelectAllFrom(): Can't connect to database"}}
	}
	// SQL Request
	var results []map[string]interface{}
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		query := fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, limit)
		rows, err := infosBdd.bdd.Query(query)
		if err != nil {
			return errors.New("SelectAllFrom(): querying table data")
		}
		defer rows.Close()
		// Take columns
		columns, err := rows.Columns()
		if err != nil {
			return errors.New("SelectAllFrom(): getting columns")
		}
		// Create the dataframe
		for rows.Next() {
			columnPointers := make([]interface{}, len(columns))
			columnValues := make([]interface{}, len(columns))
			for i := range columnValues {
				columnPointers[i] = &columnValues[i]
			}
			if err := rows.Scan(columnPointers...); err != nil {
				return errors.New("SelectAllFrom(): scanning row: " + table)
			}

			rowMap := make(map[string]interface{})
			for i, colName := range columns {
				rowMap[colName] = columnValues[i]
			}
			results = append(results, rowMap)
		}
		if err := rows.Err(); err != nil {
			return errors.New("SelectAllFrom(): during rows iteration: " + table)
		}
		if len(results) == 0 {
			return errors.New("no value in: " + table)
		}
		return nil
	})
	if err != nil {
		return []map[string]interface{}{{"Error": err.Error()}}
	}
	return results
}

func (adb Aquabase) SelectFrom(sqlQuery string) []map[string]interface{} {
	// Open sqliteDB
	infosBdd, err := adb.Login()
	if err != nil {
		return []map[string]interface{}{{"Error": "SelectFrom(): Can't connect to database"}}
	}
	// SQL Request
	var results []map[string]interface{}
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		rows, err := infosBdd.bdd.Query(sqlQuery)
		if err != nil {
			return errors.New("SelectFrom(): querying table data")
		}
		defer rows.Close()
		// Take columns
		columns, err := rows.Columns()
		if err != nil {
			return errors.New("SelectFrom(): getting columns")
		}
		// Create the dataframe
		for rows.Next() {
			columnPointers := make([]interface{}, len(columns))
			columnValues := make([]interface{}, len(columns))
			for i := range columnValues {
				columnPointers[i] = &columnValues[i]
			}
			if err := rows.Scan(columnPointers...); err != nil {
				return errors.New("SelectFrom(): scanning row: " + sqlQuery)
			}
			rowMap := make(map[string]interface{})
			for i, colName := range columns {
				rowMap[colName] = columnValues[i]
			}
			results = append(results, rowMap)
		}
		if err := rows.Err(); err != nil {
			return errors.New("SelectFrom(): during rows iteration: " + sqlQuery)
		}
		return nil
	})
	if err != nil {
		return []map[string]interface{}{{}}
	}
	return results
}

func (adb Aquabase) EstTableVide(table string) bool {
	infosBdd, err := GetInfosBDD(adb.dbPath)
	if err != nil {
		log.Println("[ERROR] Problème dans l'ouverture de la base : ", err)
		return true
	}
	var contientDonnees bool
	err = infosBdd.tickets.ExecutionQuandTicketPret(func() error {
		resultat, err := infosBdd.bdd.Query("SELECT * FROM " + table + " LIMIT 1")
		if err != nil {
			log.Println("[ERROR] Problème dans la récupération des informations de la table : ", err)
			return err
		}
		defer resultat.Close()
		contientDonnees = resultat.Next()
		return nil
	})
	if err != nil {
		return true
	}
	return !contientDonnees
}

/* -------------------------- FONCTIONS ANNEXES -------------------------- */

func nettoyage(entree string) string {
	entree = strings.ReplaceAll(entree, "\"", "&quot;")
	entree = strings.ReplaceAll(entree, "<", "&lt;")
	entree = strings.ReplaceAll(entree, ">", "&gt;")
	entree = strings.ReplaceAll(entree, "&", "&amp")
	return entree
}
