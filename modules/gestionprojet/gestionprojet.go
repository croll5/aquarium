package gestionprojet

import (
	"aquarium/modules/extraction"
	"aquarium/modules/extraction/utilitaires"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bodgit/sevenzip"

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
	// Création de la base de données qui contiendra la chronologie des évènements
	extraction.CreationBaseAnalyse(chemin)
	// Creation d'un dossier contenant les règles de detection de l'utilisateur
	os.MkdirAll(filepath.Join(chemin, "regles_detection"), os.ModeDir)
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
		ExtractArchive7z(listeOrcs[i], filepath.Join(cheminAnalyse, "collecteORC", strings.Replace(caracteristiques[3], ".7z", "", 1)))
	}
	return true
}

/* Fonction permettant de décompresser un dossier compressé en 7z
 * Cette fonction utilise la bibliothèque sevenzip,
 * dont la documentation est présente ici : https://pkg.go.dev/github.com/bodgit/sevenzip
 */
func ExtractArchive7z(archive string, destination string) error {
	r, err := sevenzip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err = utilitaires.ExtraireFichierDepuis7z(f, destination); err != nil {
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
