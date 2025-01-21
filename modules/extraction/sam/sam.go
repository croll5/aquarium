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

package sam

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/extraction/utilitaires"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bodgit/sevenzip"
	"www.velocidex.com/golang/regparser"
)

/* VARIABLES GLOBALES */
var pourcentageChargement float32 = -1

var colonnesTableSam []string = []string{"horodatage", "idCompte", "nomCompte", "operation", "source"}

type Sam struct{}

func traiterInfosCompte(compte *regparser.CM_KEY_NODE, dejaFait *map[string][]bool, nomCompte string, source string, requete *aquabase.RequeteInsertion) {
	if compte.Name() != "Names" {
		fmt.Println()
		log.Println(compte.Name(), ":", compte.Values()[0].ValueName())
		var DonneesF []byte = compte.Values()[0].ValueData().Data
		var minimum time.Time = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		if (*dejaFait)[compte.Name()] == nil {
			(*dejaFait)[compte.Name()] = []bool{false, false}
		}
		// Eventuel ajout de la date de dernière connexion
		var derniereConnexion time.Time = utilitaires.FileTimeVersGo(DonneesF[8:16])
		if derniereConnexion.After(minimum) && derniereConnexion.Before(time.Now()) && (!(*dejaFait)[compte.Name()][0]) {
			requete.AjouterDansRequete(derniereConnexion.String(), compte.Name(), nomCompte, "derniereConnexion", source)
			(*dejaFait)[compte.Name()] = []bool{true, (*dejaFait)[compte.Name()][1]}
		}
		var creationCompte = utilitaires.FileTimeVersGo(DonneesF[24:32])
		if creationCompte.After(minimum) && creationCompte.Before(time.Now()) && (!(*dejaFait)[compte.Name()][1]) {
			requete.AjouterDansRequete(creationCompte.String(), compte.Name(), nomCompte, "creation", source)
			(*dejaFait)[compte.Name()] = []bool{(*dejaFait)[compte.Name()][0], true}
		}
	}
}

func (s Sam) Extraction(cheminProjet string) error {
	// Ouverture du fichier SAM.7z, qui contient les fichiers SAM
	pourcentageChargement = 0
	r, err := sevenzip.OpenReaderWithPassword(filepath.Join(cheminProjet, "collecteORC", "SAM", "SAM.7z"), "avproof")
	if err != nil {
		return err
	}
	defer r.Close()
	var dejaFait map[string][]bool = map[string][]bool{}
	// Parcourt des fichiers contenus dans SAM.7z
	for numFichier, fichierSAM := range r.File {
		rc, err := fichierSAM.Open()
		if err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		defer rc.Close()
		// Copie du contenu du fichier dans un tampon, pour pouvoir l'ouvrir avec l'extracteur de registres
		var tampon bytes.Buffer
		if _, err := io.Copy(&tampon, rc); err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		readerAt := bytes.NewReader(tampon.Bytes())
		// Ouverture du fichier comme fichier de registre
		registre, err := regparser.NewRegistry(readerAt)
		if registre == nil {
			break
		}
		if err != nil {
			log.Println("Format de fichier non supporté  : ", err.Error())
		}
		// Ouverture de la clé de registre contenant les comptes personnels
		cleDeBase := registre.OpenKey("SAM/Domains/Account/Users/Names")
		// Récupération de toutes les clés enfants, donc des clés de comptes personnels
		enfants := cleDeBase.Subkeys()
		var nomsDesComptes map[string]string = map[string]string{}
		for _, compte := range enfants {
			// Ajout de l'évènement de dernière modification du compte à la BDD
			var idCompte string = strings.ToUpper(fmt.Sprintf("00000%x", compte.Values()[0].Type()))
			nomsDesComptes[idCompte] = compte.Name()
			if err != nil {
				return err
			}
		}
		deuxiemeEssai := registre.OpenKey("SAM/Domains/Account/Users")
		enfants = deuxiemeEssai.Subkeys()
		var abase aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
		var requete aquabase.RequeteInsertion = abase.InitRequeteInsertionExtraction("sam", colonnesTableSam)
		for _, compte := range enfants {
			traiterInfosCompte(compte, &dejaFait, nomsDesComptes[compte.Name()], fichierSAM.Name, &requete)
		}
		requete.Executer()
		pourcentageChargement = float32(numFichier) * 100 / float32(len(r.File))
	}
	pourcentageChargement = 101
	return nil
}

func (s Sam) CreationTable(cheminProjet string) error {
	aqua := aquabase.InitDB_Extraction(cheminProjet)
	aqua.CreateTableIfNotExist1("sam", colonnesTableSam, true)
	return nil
}

func (s Sam) PourcentageChargement(cheminProjet string, verifierTableVide bool) float32 {
	if pourcentageChargement == -1 {
		bdd := aquabase.InitDB_Extraction(cheminProjet)
		if !bdd.EstTableVide("sam") {
			pourcentageChargement = 100
		}
	}
	return pourcentageChargement
}

func (s Sam) Annuler() bool {
	// Il n'y a pas suffisement de fichiers pour que cela ne soit pertinent.
	// On attend donc simplement que l'extraction soit terminée.
	return pourcentageChargement >= 100
}

func (s Sam) Description() string {
	return "Extraction des fichiers SAM (contenant la base d'utilisateurs)"
}

func (s Sam) PrerequisOK(cheminCollecte string) bool {
	dossierSAM, err := os.ReadDir(filepath.Join(cheminCollecte, "SAM"))
	if err != nil {
		return false
	}
	for _, fichier := range dossierSAM {
		if fichier.Name() == "SAM.7z" {
			return true
		}
	}
	return false
}

func (s Sam) DetailsEvenement(idEvt int) string {
	return "Pas d'informations supplémentaires"
}

func (s Sam) SQLChronologie() string {
	return "SELECT id, \"sam\", \"sam\", source, horodatage, \"opération \" || operation || \" effecutée sur le compte \" || nomCompte || \" (idCompte : \" || idCompte || \")\" FROM sam"
}
