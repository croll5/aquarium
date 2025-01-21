/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025)

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à [rappeler les
caractéristiques techniques de votre logiciel].

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

package gestionprojet

import (
	"aquarium/modules/extraction"
	"aquarium/modules/extraction/utilitaires"
	"aquarium/modules/rapport"
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
func CreationArborescence(chemin *string) bool {
	// Création d'un fichier .aqua contenant les infos essentielles du projet
	estvide, err := IsDirEmpty(*chemin)
	if err != nil || !estvide {
		*chemin = filepath.Join(*chemin, time.Now().Format("20060102150405")+" - analyse aquarium")
		log.Println(*chemin)
	}
	os.MkdirAll(filepath.Join(*chemin, "analyse"), 0766)
	fichier, err := os.Create(filepath.Join(*chemin, "analyse.aqua"))
	if err != nil {
		log.Println(err)
		return false
	}
	defer fichier.Close()
	// Création de la base de données qui contiendra la chronologie des évènements
	extraction.CreationBaseAnalyse(*chemin)
	// Création de la table des informations spécifiques au rapport
	var rprt rapport.Rapport = *rapport.InitRapport(*chemin)
	rprt.CreerTables()
	// Creation d'un dossier contenant les règles de detection de l'utilisateur
	os.MkdirAll(filepath.Join(*chemin, "regles_detection"), 0766)
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
	os.MkdirAll(filepath.Join(chemin, "analyse"), 0766)
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
