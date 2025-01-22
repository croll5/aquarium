/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025)

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

package utilitaires

import (
	"database/sql"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bodgit/sevenzip"
)

/*
* Fonction permettant l'insersion d'un évènement dans la table “chronologie“ de la base de données
@cheminProjet : la racine du projet aqua
@extracteur : l'identifiant de l'extracteur
@horodatage : la date à laquelle l'évènement a eu lieu
@source : le fichier duquel a été extrait l'évènement
@message : la destription de l'évènement
@return : une erreur s'il y en a eu une
*
*/
func AjoutEvenementDansBDD(cheminProjet string, extracteur string, horodatage time.Time, source string, message string) error {
	bd, err := sql.Open("sqlite", filepath.Join(cheminProjet, "analyse", "extractions.db"))
	//log.Println(filepath.Join(cheminProjet, "analyse", "extractions.db"))
	if err != nil {
		return err
	}
	defer bd.Close()
	requete, err := bd.Prepare("INSERT INTO chronologie(extracteur, horodatage, source, message) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = requete.Exec(extracteur, horodatage, source, message)
	return err
}

func FileTimeVersGo(date []byte) time.Time {
	var dateInt = int64(binary.LittleEndian.Uint64(date))
	var difference = dateInt / 10000000
	var complement = dateInt % 10000000
	var referentiel = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(referentiel+difference, complement)
}

/* Fonction permettant d'extraire un fichier d'un dossier compressé en 7z
 */
func ExtraireFichierDepuis7z(file *sevenzip.File, destination string) error {
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
