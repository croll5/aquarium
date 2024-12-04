package gestionprojet

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bodgit/sevenzip"

	"database/sql"

	_ "modernc.org/sqlite"
)

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

/*
  - Fonction qui crée l'arborescence de base de l'analyse
    @chemin : le chemin dans lequel on veut créer le projet
*/
func CreationArborescence(chemin string) bool {
	// Création d'un fichier .aqua contenant les infos essentielles du projet
	estvide, err := IsDirEmpty(chemin)
	if err != nil || !estvide {
		return false
	}
	os.MkdirAll(filepath.Join(chemin, "analyse"), os.ModeDir)
	fichier, err := os.Create(filepath.Join(chemin, "analyse.aqua"))
	if err != nil {
		log.Println(err)
		return false
	}
	defer fichier.Close()
	fichier.WriteString("coucou")
	// Création de la base de données qui contiendra la chronologie des évènements
	bd, err := sql.Open("sqlite", filepath.Join(chemin, "analyse", "extractions.db"))
	if err != nil {
		log.Println(err)
		return false
	}
	defer bd.Close()
	var requete string = "CREATE TABLE chronologie(id INT PRIMARY KEY, extracteur VARCHAR(25), horodatage DATETIME, source VARCHAR(25), message TEXT)"
	bd.Exec(requete)
	requete = "CREATE TABLE indicateurs(id INT PRIMARY KEY, type VARCHAR(32), valeur VARCHAR(50), tlp VARCHAR(10), pap VARCHAR(10), commentaire TEXT)"
	bd.Exec(requete)
	requete = "CREATE TABLE indicateurs_evenements(id_indicateur INT, id_evenement INT, FOREIGN KEY(id_indicateur) REFERENCES indicateurs(id), FOREIGN KEY(id_evenement) REFERENCES chronologie(id))"
	bd.Exec(requete)
	requete = "CREATE TABLE navigateurs(id INTEGER PRIMARY KEY AUTOINCREMENT, horodatage DATETIME, url VARCHAR(50), title VARCHAR(50), domain_name VARCHAR(25), visit_count INT)"
	bd.Exec(requete)
	return true
}

func CreationDossierModele(chemin string) error {
	estvide, err := IsDirEmpty(chemin)
	if err != nil {
		log.Println(err)
		return err
	}
	if !estvide {
		log.Println(err)
		return errors.New("Le dossier " + chemin + " n'est pas vide.")
	}
	os.MkdirAll(filepath.Join(chemin, "analyse"), os.ModeDir)
	fichier, err := os.Create(filepath.Join(chemin, "modele.aqua"))
	if err != nil {
		log.Println(err)
		return err
	}
	defer fichier.Close()
	return nil
}

/* Fonction permettant l'ouverture des ORCs et leur copie dans le répertoire de l'analyse
 *
 */
func RecuperationOrcs(listeOrcs []string, cheminAnalyse string) bool {
	if len(listeOrcs) == 0 {
		return false
	}
	// Dézippage des fichiers ORC donnés par l'utilisateur
	for i := 0; i < len(listeOrcs); i++ {
		log.Println("INFO | Dézippage du fichier ", listeOrcs[i])
		// Vérification que le fichier donné est du bon format
		var nomFichierOrc = filepath.Base(listeOrcs[i])
		var caracteristiques []string = strings.Split(nomFichierOrc, "_")
		if caracteristiques[0] != "DFIR-ORC" || len(caracteristiques) != 4 {
			log.Println("ERROR | Le nom du fichier ORC donné en argument doit commencer par \"DFIR-ORC\"")
			return false
		}
		// Ectraction à proprement parler
		extractArchive(listeOrcs[i], filepath.Join(cheminAnalyse, "collecteORC", strings.Replace(caracteristiques[3], ".7z", "", 1)))
	}
	return true
}

/* Fonction permettant d'extraire un fichier d'un dossier compressé en 7z
 */
func extractFile(file *sevenzip.File, destination string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	log.Println("INFO | Dezippage du fichier ", file.Name)
	os.MkdirAll(filepath.Join(destination, filepath.Dir(file.Name)), 0755)
	fichierExtrait, err := os.Create(filepath.Join(destination, file.Name))
	if err != nil {
		log.Println("ERROR | Problème dans la création du fichier de copie : ", err.Error())
	}
	defer fichierExtrait.Close()

	_, err = io.Copy(fichierExtrait, rc)
	if err != nil {
		log.Println("ERROR | Problème dans l'extraction de l'ORC : ", err.Error())
	}

	return nil
}

/* Fonction permettant de décompresser un dossier compressé en 7z
 * Cette fonction utilise la bibliothèque sevenzip,
 * dont la documentation est présente ici : https://pkg.go.dev/github.com/bodgit/sevenzip
 */
func extractArchive(archive string, destination string) error {
	r, err := sevenzip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err = extractFile(f, destination); err != nil {
			return err
		}
	}

	return nil
}

func EcritureFichierAqua(nomAnalyste string, description string, debutAnalyse time.Time, finAnalyse time.Time, cheminProjet string) error {
	var debut string = debutAnalyse.Format("02/01/2006 15 h 04")
	var fin string = finAnalyse.Format("02/01/2006 15 h 04")
	caracteristiques := map[string]string{"nom_auteur": nomAnalyste, "debut_analyse": debut, "fin_analyse": fin, "description": description}
	caracteristiques_json, err := json.Marshal(caracteristiques)
	if err != nil {
		log.Println("Problème dans la conversion de map en json : ", err.Error())
		return err
	}
	fichier, err := os.OpenFile(filepath.Join(cheminProjet, "analyse.aqua"), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("Problème dans l'ouverture du fichier analyse.aqua : ", err.Error())
		return err
	}
	defer fichier.Close()
	_, err = fichier.Write(caracteristiques_json)
	if err != nil {
		log.Println("Problème dans l'écriture des données aqua : ", err.Error())
		return err
	}
	return nil
}

func EcritureFichierModeleAqua(nomModele string, description string, dateCreation time.Time, cheminProjet string) error {
	var creation string = dateCreation.Format("02/01/2006 15 h 04")
	caracteristiques := map[string]string{"nom_modele": nomModele, "date_creation": creation, "description": description}
	caracteristiques_json, err := json.Marshal(caracteristiques)
	if err != nil {
		log.Println("Problème dans la conversion de map en json : ", err.Error())
		return err
	}
	fichier, err := os.OpenFile(filepath.Join(cheminProjet, "modele.aqua"), os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("Problème dans l'ouverture du fichier modele.aqua : ", err.Error())
		return err
	}
	defer fichier.Close()
	_, err = fichier.Write(caracteristiques_json)
	if err != nil {
		log.Println("Problème dans l'écriture des données aqua : ", err.Error())
		return err
	}
	return nil
}
