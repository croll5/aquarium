package sam

import (
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
	r, err := sevenzip.OpenReaderWithPassword(filepath.Join(cheminProjet, "collecteORC", "SAM", "SAM.7z"), "avproof")
	defer r.Close()
	if err != nil {
		return err
	}
	for _, fichierSAM := range r.File {
		rc, err := fichierSAM.Open()
		if err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		defer rc.Close()

		var tampon bytes.Buffer
		if _, err := io.Copy(&tampon, rc); err != nil {
			log.Println("Format de fichier non supporté : ", err.Error())
		}
		readerAt := bytes.NewReader(tampon.Bytes())
		registre, err := regparser.NewRegistry(readerAt)
		if registre == nil {
			break
		}
		if err != nil {
			log.Println("Format de fichier non supporté  : ", err.Error())
		}
		cleDeBase := registre.OpenKey("SAM/Domains/Account/Users/Names")
		enfants := cleDeBase.Subkeys()

		log.Println(enfants[0].Name(), " : ", enfants[0].LastWriteTime().Time)
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
