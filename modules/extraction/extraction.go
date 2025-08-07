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

package extraction

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/config"
	"aquarium/modules/extraction/bdd_sqlite"
	"aquarium/modules/extraction/csv"
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/registre"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bodgit/sevenzip"
)

type Extracteur interface {
	Extraction(string, bytes.Buffer, string, config.ConfigExtraction) error
}

var liste_extracteurs map[string]Extracteur = map[string]Extracteur{
	// "avs":        avlogs.AvLog{},
	"evtx":   evtx.Evtx{},
	"sqlite": bdd_sqlite.SQLite{},
	// "navigateur": navigateur.Navigateur{},
	// "werr":       werr.Werr{},
	"registre": registre.Registre{},
	"csv":      csv.Csv{},
	// "divers":     divers.Divers{},
	// "prefetch":   prefetch.Prefetch{},
}

var liste_extractions map[string]config.ConfigExtraction = map[string]config.ConfigExtraction{}

/** Fonction qui renvoie des informations sur les extractions qui peuvent être effectuées
  * @param cheminProjet string : le chemin d'enregistrement de l'analyse aquarium
  * @return une liste de configurations d'extractions, et s'il y a lieu une erreur
**/
func ListeExtracteursHtml(cheminProjet string) (map[string]config.ConfigExtraction, error) {
	// On commence par récupérer la liste des extractions dans le fichier de config
	config, err := config.GetConfigurationProjet()
	if err != nil {
		return liste_extractions, err
	}
	// On itère sur toutes les extractions
	for _, extracteur := range config.Extractions {
		val, ok := liste_extractions[extracteur.Id]
		if !ok || (val.Progression == -1) {
			// TODO: Ajouter une vérification que le chemin existe
			extracteur.Progression = -1
			var adb *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
			if !adb.EstTableVide(extracteur.Table.Nom) {
				extracteur.Progression = 100
			}
			extracteur.AnnulationDemandee = false
			liste_extractions[extracteur.Id] = extracteur
		}
	}
	return liste_extractions, nil
}

/** Fonction qui exécute une extraction définie dans le fichier de configuration
  * @param idExtraction : l'identifiant de l'extraction
  * @cheminProjet string : le chemin d'enregistrement du projet
  * @return : une erreur s'il y a lieu
**/
func Extraction(idExtraction string, cheminProjet string) error {
	log.Println(liste_extractions[idExtraction].Complement)
	var probleme error
	// On liste les fichiers concernés par cette extraction
	listeFichiersAExtraire, probleme := config.ListeFichiersExtraction(liste_extractions[idExtraction], cheminProjet)
	log.Println(listeFichiersAExtraire)
	// On récupère la configuration de l'extraction
	var configExtraction config.ConfigExtraction = liste_extractions[idExtraction]
	// On met la progression à 0 (début de l'extraction)
	configExtraction.Progression = 0
	liste_extractions[idExtraction] = configExtraction
	// On crée la table qui sera utilisée par l'extracteur
	err := creerTableExtraction(cheminProjet, configExtraction)
	if err != nil {
		return err
	}
	// On compte le nombre de fichiers à extraire
	var nbFichiers int = 0
	for _, dossierAExtraire := range listeFichiersAExtraire {
		nbFichiers += len(dossierAExtraire.Elements)
	}
	log.Println(nbFichiers, " fichiers à extraire")
	var i int = 0
	// On boucle sur les fichiers à extraire
	for _, dossierAExtraire := range listeFichiersAExtraire {
		if liste_extractions[idExtraction].AnnulationDemandee {
			configExtraction.Progression = -1
			configExtraction.AnnulationDemandee = false
			liste_extractions[idExtraction] = configExtraction
			var adb *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
			adb.DropTable(configExtraction.Table.Nom)
			return probleme
		}
		if dossierAExtraire.Est7Z {
			extrationAchive7z(cheminProjet, dossierAExtraire, idExtraction, &i, nbFichiers)
		} else {
			extractionDossier(cheminProjet, dossierAExtraire, idExtraction, &i, nbFichiers)
		}
	}
	configExtraction.Progression = 101
	liste_extractions[idExtraction] = configExtraction
	return probleme
}

/** Fonction qui crée toutes les tables de l'analyse
  * @param cheminProjet string : le chemin d'enregistrement du projet
  * TODO: Essayer de se passer de cette fonction
**/
func CreationBaseAnalyse(cheminProjet string) {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	for _, extraction := range liste_extractions {
		// On récupère une liste des colonnes
		creerTableExtraction(cheminProjet, extraction)
	}
	configAnalyse, err := config.GetConfigurationProjet()
	if err != nil {
		return
	}
	var listeColonnesChronologie map[string]string = map[string]string{}
	for _, colonne := range configAnalyse.Chronologie.Colonnes {
		listeColonnesChronologie[colonne.Nom] = colonne.Type
	}
	base.CreateTableIfNotExist2(configAnalyse.Chronologie.Nom, listeColonnesChronologie, true)
}

/** Fonction qui renvoie le pourcentage de chargement de l'extraction
 ** @param cheminProjet string : le chemin de l'analyse aquarium
 ** @param idExtraction string : l'identifiant de l'extraction
**/
func ProgressionExtraction(cheminProjet string, idExtraction string) float32 {
	return liste_extractions[idExtraction].Progression
}

/** Fonction permettant d'annuler une extraction en cours
  * @param idExtraction string : l'identifiant de l'extraction à annuler
  * @return un booléen indiquant si l'extraction a bien pu être annulée
**/
func AnnulerExtraction(idExtraction string) bool {
	var configExtraction config.ConfigExtraction = liste_extractions[idExtraction]
	configExtraction.AnnulationDemandee = true
	liste_extractions[idExtraction] = configExtraction
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		if !liste_extractions[idExtraction].AnnulationDemandee {
			ticker.Stop()
			return true
		}
	}
	time.Sleep(30 * time.Second)
	return false
}

func DetailsEvenement(idExtraction string, idEvenement int) string {
	return liste_extractions[idExtraction].Description
}

/** Fonction qui lance l'extraction de la table chronologie, qui contient un
  * résumé de tous les évènements
  * @param cheminprojet : le chemin d'enregistrament de l'analyse aquarium
  * @return : une erreur s'il y a lieu
**/
func ExtraireTableChronologie(cheminProjet string) error {
	var listeRequetesChronologie []string = []string{}
	// On fait une liste des requêtes SQL à exécuter
	for _, extraction := range liste_extractions {
		if extraction.SQLChronologie != "" {
			listeRequetesChronologie = append(listeRequetesChronologie, extraction.SQLChronologie)
		}
	}
	// On liste les colonnes de la table chronologie
	configAnalyse, err := config.GetConfigurationProjet()
	if err != nil {
		return err
	}
	var listeColonnesChrolonogie []string = []string{}
	for _, colonne := range configAnalyse.Chronologie.Colonnes {
		listeColonnesChrolonogie = append(listeColonnesChrolonogie, colonne.Nom)
	}
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	err = base.RemplirTableDepuisRequetes(configAnalyse.Chronologie.Nom, listeColonnesChrolonogie, listeRequetesChronologie, true, "horodatage")
	return err
}

/** Fonction qui renvoie le contenu de la tables "chronologie", contenant un résumé de l'ensemble des
  * évènemets
  * @param cheminprojet string : le chemin d'enregistrement de l'analyse aquarium
  * @param debut int : l'index à partir duquel on veut récupérer les valeurs
  * @param taill int : le nombre de valeurs que l'on veut récupérer
  * @return une liste des lignes de la table
**/
func ValeursTableChronologie(cheminProjet string, debut int, taille int) []map[string]interface{} {
	var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	configAnalyse, err := config.GetConfigurationProjet()
	if err != nil {
		return []map[string]interface{}{}
	}
	var listeColonnesChrolonogie []string = []string{}
	for _, colonne := range configAnalyse.Chronologie.Colonnes {
		listeColonnesChrolonogie = append(listeColonnesChrolonogie, colonne.Nom)
	}
	return abase.RecupererValeursTable("chronologie", listeColonnesChrolonogie, debut, taille)
}

/** --------------------- FONCTIONS À USAGE INTERNE --------------------- **/

func creerTableExtraction(cheminProjet string, extraction config.ConfigExtraction) error {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	var listeColonnes map[string]string = map[string]string{}
	for _, colonne := range extraction.Table.Colonnes {
		listeColonnes[colonne.Nom] = colonne.Type
	}
	err := base.CreateTableIfNotExist2(extraction.Table.Nom, listeColonnes, true)
	return err
}

func extrationAchive7z(cheminProjet string, configArchive config.DossierAExtraire, idExtraction string, i *int, total int) error {
	// On commence par ouvrir l'archive
	archive, err := sevenzip.OpenReaderWithPassword(configArchive.Chemin, "avproof")
	if err != nil {
		return err
	}
	for _, numFichier := range configArchive.Elements {
		fichier, err := archive.File[numFichier].Open()
		if err != nil {
			continue
		}
		defer fichier.Close()
		// Copie du contenu du fichier dans un tampon, pour pouvoir l'ouvrir avec l'extracteur de registres
		var tampon bytes.Buffer
		if _, err := io.Copy(&tampon, fichier); err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		var source string = strings.Replace(filepath.Join(configArchive.Chemin, archive.File[numFichier].Name), cheminProjet, "", 1)
		liste_extracteurs[liste_extractions[idExtraction].Extracteur].Extraction(cheminProjet, tampon, source, liste_extractions[idExtraction])
		fichier.Close()
		// On change la progression du chargement
		*i++
		var confExtraction config.ConfigExtraction = liste_extractions[idExtraction]
		confExtraction.Progression = float32(*i) / float32(total) * 100
		liste_extractions[idExtraction] = confExtraction
	}
	return nil
}

func extractionDossier(cheminProjet string, configDossier config.DossierAExtraire, idExtraction string, i *int, total int) error {
	// On commence par lire le dossier
	contenuDossier, err := os.ReadDir(configDossier.Chemin)
	if err != nil {
		return err
	}
	for _, numFichier := range configDossier.Elements {
		var cheminFichier string = filepath.Join(configDossier.Chemin, contenuDossier[numFichier].Name())
		fichier, err := os.Open(cheminFichier)
		if err != nil {
			continue
		}
		defer fichier.Close()
		var tampon bytes.Buffer
		if _, err := io.Copy(&tampon, fichier); err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		liste_extracteurs[liste_extractions[idExtraction].Extracteur].Extraction(cheminProjet, tampon, strings.Replace(cheminFichier, cheminProjet, "", 1), liste_extractions[idExtraction])
		fichier.Close()
	}
	return nil
}
