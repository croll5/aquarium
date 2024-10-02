package sam

import (
	"aquarium/modules/extraction/utilitaires"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"

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
		enfants := cleDeBase.Subkeys()
		for _, compte := range enfants {
			// Ajout de l'évènement de dernière modification du compte à la BDD
			err := utilitaires.AjoutEvenementDansBDD(cheminProjet, "sam", compte.LastWriteTime().Time, "SAM/SAM/"+fichierSAM.Name, "Modification du compte "+compte.Name())
			if err != nil {
				return err
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
