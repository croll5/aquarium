package navigateur

import (
    "io/ioutil"
    "log"
    "fmt"
)


type Navigateur struct {
	extrait bool
}

func (n Navigateur) Extraction(chemin_projet string) error {
    // Trouver le bon dossier
    // Trouver les fichiers à analyser là dedans :')
    /*
    * B030D72430D6F078_190000001A7F07_E0000001A9CF0_4_places.sqlite_{00000000-0000-0000-0000-000000000000}.data => Firefox (Faire la recherche sur les fichiers ayant "places.sqlite_")
    *
    */
    data, err := ioutil.ReadFile("B030D72430D6F078_190000001A7F07_E0000001A9CF0_4_places.sqlite_{00000000-0000-0000-0000-000000000000}.data");
    if err != nil{
        log.Panicf("File not found or not openable");
    }
    fmt.Printf("\nLength: %d bytes", len(data))
    fmt.Printf("\nData: %s", data)
    fmt.Printf("\nError: %v", err)
	return nil
}

func (n Navigateur) Description() string {
	return "Historique de navigation"
}

func (n Navigateur) PrerequisOK(cheminORC string) bool {
	return true
}
