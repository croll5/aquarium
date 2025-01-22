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
	"aquarium/modules/extraction/avlogs"
	"aquarium/modules/extraction/divers"
	"aquarium/modules/extraction/evtx"
	"aquarium/modules/extraction/getthis"
	"aquarium/modules/extraction/navigateur"
	"aquarium/modules/extraction/prefetch"
	"aquarium/modules/extraction/sam"
	"aquarium/modules/extraction/werr"
	"errors"
	"log"
	"path/filepath"
	"time"
)

type Extracteur interface {
	Extraction(string) error
	Description() string
	PrerequisOK(string) bool
	CreationTable(string) error
	PourcentageChargement(string, bool) float32
	Annuler() bool
	DetailsEvenement(int) string
	SQLChronologie() string
}

type InfosExtracteur struct {
	Description string
	Progression float32
}

var liste_extracteurs map[string]Extracteur = map[string]Extracteur{
	"avs":        avlogs.AvLog{},
	"evtx":       evtx.Evtx{},
	"navigateur": navigateur.Navigateur{},
	"werr":       werr.Werr{},
	"sam":        sam.Sam{},
	"getthis":    getthis.Getthis{},
	"divers":     divers.Divers{},
	"prefetch":   prefetch.Prefetch{},
}

var colonnesTableChronologie map[string]string = map[string]string{"idEvt": "INT", "extracteur": "TEXT", "nomTable": "TEXT", "source": "TEXT", "horodatage": "DATETIME", "message": "TEXT"}
var colonnesSimmplesChronologie []string = []string{"idEvt", "extracteur", "nomTable", "source", "horodatage", "message"}

func ListeExtracteursHtml(cheminProjet string) (map[string]InfosExtracteur, error) {
	// On itère sur tous les extracteurs
	var resultat = map[string]InfosExtracteur{}
	for k, v := range liste_extracteurs {
		//log.Println(filepath.Join(cheminProjet, "collecteORC"))
		if v.PrerequisOK(filepath.Join(cheminProjet, "collecteORC")) {
			resultat[k] = InfosExtracteur{Description: v.Description(), Progression: v.PourcentageChargement(cheminProjet, true)}
		}
	}
	return resultat, nil
}

func Extraction(module string, cheminProjet string) error {
	if liste_extracteurs[module] == nil {
		return errors.New("Erreur : module " + module + " non reconnu")
	}

	err := liste_extracteurs[module].Extraction(cheminProjet)
	return err
}

func CreationBaseAnalyse(cheminProjet string) {
	for _, extracteur := range liste_extracteurs {
		extracteur.CreationTable(cheminProjet)
	}
	var base aquabase.Aquabase = *aquabase.InitDB_Extraction(cheminProjet)
	base.CreateTableIfNotExist2("chronologie", colonnesTableChronologie, true)
}

func ProgressionExtraction(cheminProjet string, idExtracteur string) float32 {
	return liste_extracteurs[idExtracteur].PourcentageChargement(cheminProjet, false)
}

func AnnulerExtraction(idExtracteur string) bool {
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		if liste_extracteurs[idExtracteur].Annuler() {
			ticker.Stop()
			return true
		}
	}
	time.Sleep(30 * time.Second)
	return false
}

func DetailsEvenement(idExtracteur string, idEvenement int) string {
	return liste_extracteurs[idExtracteur].DetailsEvenement(idEvenement)
}

func ExtraireTableChronologie(cheminProjet string) error {
	var listeRequetesChronologie []string = []string{}
	for _, extracteur := range liste_extracteurs {
		if extracteur.SQLChronologie() != "" {
			listeRequetesChronologie = append(listeRequetesChronologie, extracteur.SQLChronologie())
		}
	}
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	err := base.RemplirTableDepuisRequetes("chronologie", colonnesSimmplesChronologie, listeRequetesChronologie, true, "horodatage")
	log.Println(err)
	return nil
}

func ValeursTableChronologie(cheminProjet string, debut int, taille int) []map[string]interface{} {
	var abase *aquabase.Aquabase = aquabase.InitDB_Extraction(cheminProjet)
	return abase.RecupererValeursTable("chronologie", colonnesSimmplesChronologie, debut, taille)
}
