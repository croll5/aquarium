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
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"database/sql"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bodgit/sevenzip"
	"github.com/pkg/errors"
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
		return errors.WithStack(err)
	}
	defer bd.Close()
	requete, err := bd.Prepare("INSERT INTO chronologie(extracteur, horodatage, source, message) VALUES (?, ?, ?, ?)")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = requete.Exec(extracteur, horodatage, source, message)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func GetFonctionDecodageBytes(encodage string) func([]byte) interface{} {
	switch encodage {
	case "littleEndian64":
		return func(donnees []byte) interface{} {
			return binary.LittleEndian.Uint64(donnees)
		}
	case "filetime":
		return func(donnees []byte) interface{} {
			if binary.LittleEndian.Uint64(donnees) == 0 {
				return nil
			}
			return FileTimeVersGo(donnees)
		}

	default:
		fonctionDecodage := GetFonctionDecodageString(encodage)
		return func(donnees []byte) interface{} {
			return fonctionDecodage(string(donnees))
		}
	}
}

func GetFonctionDecodageString(encodage string) func(string) interface{} {
	detailsEncodage := strings.Split(encodage, "[aqua_sep]")
	switch detailsEncodage[0] {
	case "utf16":
		return func(donnees string) interface{} {
			return Utf16LEToUtf8(donnees)
		}
	case "filetime":
		return func(donnees string) interface{} {
			date, err := strconv.ParseInt(donnees, 10, 64)
			if err != nil {
				log.Printf("Erreur de conversion de l'horodatage : %s, erreur : %v\n", donnees, err)
				return "[AQUA] Erreur dans l’extraction de la date au format filetime suivante :" + donnees
			}
			return FiletimeFromIntVersGo(date)
		}
	case "date":
		if len(detailsEncodage) < 2 {
			return func(donnees string) interface{} {
				return "[AQUA] Date non extraite : " + donnees
			}
		} else {
			return func(donnees string) interface{} {
				date, err := time.Parse(detailsEncodage[1], donnees)
				if err != nil {
					return "[AQUA] Erreur dans l’extraction de la date " + donnees + " : " + err.Error()
				}
				return date
			}
		}
	case "string":
		return func(donnees string) interface{} {
			return donnees
		}
	default:
		return func(donnees string) interface{} {
			return "[AQUA] Impossible de décoder « " + donnees + " ». Encodage non reconnu"
		}
	}
}

func FileTimeVersGo(date []byte) time.Time {
	var dateInt = int64(binary.LittleEndian.Uint64(date))
	return FiletimeFromIntVersGo(dateInt)
}

func FiletimeFromIntVersGo(date int64) time.Time {
	var difference = date / 10000000
	var complement = date % 10000000
	var referentiel = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(referentiel+difference, complement)
}

func Utf16LEToUtf8(s string) string {
	initial := []byte(s)
	runes := []rune{}
	for i := 0; i < len(initial)-1; i += 2 {
		char := uint16(initial[i]) | uint16(initial[i+1])<<8
		runes = append(runes, rune(char))
	}
	return string(runes)
}

/* Fonction permettant d'extraire un fichier d'un dossier compressé en 7z
 */
func ExtraireFichierDepuis7z(file *sevenzip.File, destination string) error {
	rc, err := file.Open()
	if err != nil {
		return errors.WithStack(err)
	}
	defer rc.Close()

	log.Println("INFO | Dezippage du fichier ", file.Name)
	os.MkdirAll(filepath.Join(destination, filepath.Dir(file.Name)), 0755)
	fichierExtrait, err := os.Create(filepath.Join(destination, file.Name))
	if err != nil {
		return errors.WithStack(err)
	}
	defer fichierExtrait.Close()

	_, err = io.Copy(fichierExtrait, rc)
	if err != nil {
		return errors.WithStack(err)
	}

	return err
}

/*
Fonction permettant de créer une requête d’insertion en base de données
à partir d’une configuration d’extraction
*/
func CreerRequeteInstertionDepuisConfig(cheminProjet string, idMachine string, configTable *config.ConfigTableBDD) *aquabase.RequeteInsertion {
	adb := aquabase.InitDB_Extraction(cheminProjet)
	var nomsColonnesTables []string = make([]string, len(configTable.Colonnes))
	var descriptifColonnes map[string]string = make(map[string]string, 0)
	var colonnesAIndexer []int = make([]int, 0)
	for i := range configTable.Colonnes {
		nomsColonnesTables[i] = configTable.Colonnes[i].Nom
		descriptifColonnes[configTable.Colonnes[i].Nom] = configTable.Colonnes[i].Contenu
		if configTable.Colonnes[i].Indexable {
			colonnesAIndexer = append(colonnesAIndexer, i)
		}
	}
	adb.CreateTableIfNotExist2(configTable.Nom, descriptifColonnes, true)
	adb.CreerIndex(configTable.Nom, []string{"id"})
	adb.CreateTableIfNotExist2(aquabase.TABLE_CHRONOLOGIE_GLOBALE, aquabase.ColonnesTableChronologieGlobale, false)
	adb.CreerIndex(aquabase.TABLE_CHRONOLOGIE_GLOBALE, []string{"id_machine", "horodatage"})
	requeteInsertion := adb.InitRequeteInsertionExtractionAvecIndex(configTable.Nom, idMachine, nomsColonnesTables, colonnesAIndexer)
	return &requeteInsertion
}
