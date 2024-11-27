package sam

import (
	"github.com/croll5/aquarium/modules/extraction/utilitaires"
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

type Sam struct{}

func (s Sam) Extraction(cheminProjet string) error {
	// Ouverture du fichier SAM.7z, qui contient les fichiers SAM
	r, err := sevenzip.OpenReaderWithPassword(filepath.Join(cheminProjet, "collecteORC", "SAM", "SAM.7z"), "avproof")
	defer r.Close()
	if err != nil {
		return err
	}
	var dejaFait map[string][]bool = map[string][]bool{}
	// Parcourt des fichiers contenus dans SAM.7z
	for _, fichierSAM := range r.File {
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
			//err := utilitaires.AjoutEvenementDansBDD(cheminProjet, "sam", compte.LastWriteTime().Time, "SAM/SAM/"+fichierSAM.Name, "Modification du compte "+compte.Name())
			if err != nil {
				return err
			}
		}
		deuxiemeEssai := registre.OpenKey("SAM/Domains/Account/Users")
		enfants = deuxiemeEssai.Subkeys()
		for _, compte := range enfants {
			if compte.Name() != "Names" {
				fmt.Println()
				log.Println(compte.Name(), ":", compte.Values()[0].ValueName())
				var DonneesF []byte = compte.Values()[0].ValueData().Data
				var minimum time.Time = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
				if dejaFait[compte.Name()] == nil {
					dejaFait[compte.Name()] = []bool{false, false}
				}
				// Eventuel ajout de la date de dernière connexion
				var derniereConnexion time.Time = utilitaires.FileTimeVersGo(DonneesF[8:16])
				if derniereConnexion.After(minimum) && derniereConnexion.Before(time.Now()) && (!dejaFait[compte.Name()][0]) {
					var message string = "Dernière connexion au compte " + compte.Name() + " (" + nomsDesComptes[compte.Name()] + ")"
					utilitaires.AjoutEvenementDansBDD(cheminProjet, "sam", derniereConnexion, "SAM/SAM/"+fichierSAM.Name, message)
					dejaFait[compte.Name()] = []bool{true, dejaFait[compte.Name()][1]}
				}
				var creationCompte = utilitaires.FileTimeVersGo(DonneesF[24:32])
				if creationCompte.After(minimum) && creationCompte.Before(time.Now()) && (!dejaFait[compte.Name()][1]) {
					var message string = "Création du compte " + compte.Name() + " (" + nomsDesComptes[compte.Name()] + ")"
					utilitaires.AjoutEvenementDansBDD(cheminProjet, "sam", creationCompte, "SAM/SAM/"+fichierSAM.Name, message)
					dejaFait[compte.Name()] = []bool{dejaFait[compte.Name()][0], true}
				}
			}
		}
	}
	return nil
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
