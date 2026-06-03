/*
Copyright ou © ou Copr. Didier Hoizé, (20 avril 2026)

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

package params

import (
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"

	"aquarium/modules/utilitaires"
)

const (
	dossierConfig = "config"
	fichierParams = "params.xml"
)

type ParametresXML struct {
	XMLName        xml.Name `xml:"parametres" json:"-"`
	Contrastes     bool     `xml:"contrastes"`
	Dyslexie       bool     `xml:"dyslexie"`
	NonAuxBubulles bool     `xml:"non_aux_bubulles"`
	OuiAuDebug     bool     `xml:"oui_au_debug"`
}

func SauvegarderParametres(cheminBase string, contrastes bool, dyslexie bool, nonAuxBubulles bool, ouiAuDebug bool) error {
	parametres := ParametresXML{
		Contrastes:     contrastes,
		Dyslexie:       dyslexie,
		NonAuxBubulles: nonAuxBubulles,
		OuiAuDebug:     ouiAuDebug,
	}
	contenuXML, err := xml.MarshalIndent(parametres, "", "  ")
	if err != nil {
		return err
	}
	cheminConfig := filepath.Join(cheminBase, dossierConfig)
	if err = os.MkdirAll(cheminConfig, 0o755); err != nil {
		return err
	}
	cheminFichier := filepath.Join(cheminConfig, fichierParams)
	contenuFinal := append([]byte(xml.Header), contenuXML...)
	if err = os.WriteFile(cheminFichier, contenuFinal, 0o644); err != nil {
		return err
	}

	// Applique immediatement le nouvel etat de journalisation sans redemarrage.
	return utilitaires.InitLogger(cheminBase, ouiAuDebug)
}

func ChargerParametres(cheminBase string) (ParametresXML, error) {
	parametres := ParametresXML{}
	cheminFichier := filepath.Join(cheminBase, dossierConfig, fichierParams)
	contenu, err := os.ReadFile(cheminFichier)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return parametres, nil
		}
		return parametres, err
	}
	err = xml.Unmarshal(contenu, &parametres)
	if err != nil {
		return ParametresXML{}, err
	}
	return parametres, nil
}
