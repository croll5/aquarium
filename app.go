/*
Copyright ou © ou Copr. Cécile Rolland et Charles Mailley, (21 janvier 2025)

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

package main

// go mod init aquarium
// go mod tidy
// wails dev

import (
	"aquarium/modules/aquabase"
	"aquarium/modules/arborescence"
	"aquarium/modules/detection"
	"aquarium/modules/extraction"
	"aquarium/modules/gestionprojet"
	"aquarium/modules/rapport"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var chemin_projet string

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

/****************************************************************************/
/************************* APP FUNCTIONS **********************************/
/****************************************************************************/

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	aquabase.FermerToutesLesBDD()
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Call this function when a bug appear
func (a *App) signalerErreur(erreur error) {
	log.Println("ERR | Erreur non traitée : ", erreur)
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Erreur dans l'écriture du projet",
		Message: "Félicitation ! Vous venez de trouver un bogue dans le logiciel Aquarium !\n C'est cadeau : \n " + erreur.Error(),
	})
}

/***************************************************************************************/
/************************* INDEX FUNCTIONS **********************************/
/***************************************************************************************/

/*
	Cette fonction permet l'ouverture d'une analyse aquarium, à partir d'un fichier .aqua

@return : vrai si et seulement si l'analyse a été correctement ouverte
*/
func (a *App) OuvrirAnalyseExistante() bool {
	fichier, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Ouvrir une analyse existante",
		Filters: []runtime.FileFilter{{DisplayName: "Aquarium", Pattern: "analyse.aqua"}},
	})
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Problème dans la sélection du fichier",
			Message: "Veuillez sélectionner un fichier au format aquarium valide",
		})
		return false
	}
	if fichier == "" {
		return false
	}
	chemin_projet = filepath.Dir(fichier)
	return true
}

/*
	Fonction permettant de choisir le dossier dans lequel enregistrer un nouveau modèle

@return : le dossier d'enregistrement du modèle
*/
func (a *App) CreationDossierNouveauModele() string {
	// Partie création du squelette de l'analyse
	projet, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choix de l'emplacement du modèle"})
	if err != nil {
		return ""
	}
	chemin_projet = projet
	if gestionprojet.CreationDossierModele(chemin_projet) != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Problème dans la création du modèle",
			Message: "Les fichiers du modèle n'ont pas pu être créés. Vérifiez que le dossier sélectionné est vide et que vous avez les droits en écriture :/"})
		return ""
	}
	return chemin_projet
}

/***************************************************************************************/
/************************* NOUVELLE ANALYSE FUNCTIONS **********************************/
/***************************************************************************************/

/*
	Fonction permettant la création d'un nouveau projet

@return : le chemin vers le nouveau projet
*/
func (a *App) CreationNouveauProjet() string {
	// Partie création du squelette de l'analyse
	projet, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choix de l'emplacement de l'analyse"})
	if err != nil {
		return ""
	}
	chemin_projet = projet
	if !gestionprojet.CreationArborescence(&chemin_projet) {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Problème dans la création de l'analyse",
			Message: "Les fichiers d'analyse n'ont pas pu être créés. Vérifiez que le dossier sélectionné est vide et que vous avez les droits en écriture :/"})
		return ""
	}
	return chemin_projet
}

/*
	Fonction permettant l'ajout d'archives ORC dans un projet en cours de création

@return : le chemin des archives ORC choisies
*/
func (a *App) AjoutORCNouveauProjet() string {
	// Partie récupération et début de traitement de l'ORC
	orcs, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Ouverture des fichiers ORC à analyser",
		Filters: []runtime.FileFilter{{DisplayName: "7zip", Pattern: "*.7z"}}})
	if err != nil {
		return ""
	}
	if !gestionprojet.RecuperationOrcs(orcs, chemin_projet) {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Problème dans la récupération des archives ORC",
			Message: "Les archives ORC n'ont pas pu être ajoutées au projet. Vérifiez que vous avez donné des fichiers au bon format (issus d'une collecte avec l'exécutable DFIR-ORC de l'ANSSI)",
		})
		return ""
	}
	return filepath.Dir(orcs[0])
}

/*
	Fonction permettant de valider la création d'une nouvelle analyse

@return : vrai si et seulement la validation a fonctionné
*/
func (a *App) ValidationCreationProjet(nomAnalyste string, description string) bool {
	err := gestionprojet.EcritureFichierAqua(nomAnalyste, description, time.Now(), time.Time{}, chemin_projet)
	if err != nil {
		log.Println("ERR | Problème dans l'écriture du fichier .aqua : ", err.Error())
		a.signalerErreur(err)
		return false
	}
	return true
}

/*
	Fonction permettant de valider la création d'un nouveau modèle

@return : vrai si et seulement la validation a fonctionné
*/
func (a *App) ValidationCreationModele(nomModele string, description string, supprimerOrc bool) bool {
	err := gestionprojet.EcritureFichierModeleAqua(nomModele, description, time.Now(), chemin_projet)
	if err != nil {
		a.signalerErreur(err)
		return false
	}
	a.ExtraireArborescence(false)
	return true
}

/***************************************************************************************/
/************************* Extraction FUNCTIONS PAGE **********************************/
/***************************************************************************************/

/* Fonction renvoyant la liste des éléments pouvant être extraits de l'ORC
 */
func (a *App) ListeExtractionsPossibles() map[string]extraction.InfosExtracteur {
	resultat, err := extraction.ListeExtracteursHtml(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

/*
	Fonction de permettant de lancer une extraction

@param module : le nom du module à utiliser pour l'extraction
@param description : la description du module à extraire
*/
func (a *App) ExtraireElements(module string, description string) {
	err := extraction.Extraction(module, chemin_projet)
	if err != nil {
		log.Println("Erreur dans l’extraction du module", module, ":", err.Error())
		a.signalerErreur(err)
	} else {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Extraction réussie",
			Message: "L’extraction du module " + description + " s'est terminée avec succès !",
		})
	}
}

/*
	Fontion permettant d'annuler une extraction en cours

@return : vrai si et seulement si l'annulation a bien fonctionné
*/
func (a *App) AnnulerExtraction(module string) bool {
	return extraction.AnnulerExtraction(module)
}

/* Fonction permettant de connaitre le pourcentage de progression d'une extraction*/
func (a *App) ProgressionExtraction(idExtracteur string) float32 {
	return extraction.ProgressionExtraction(chemin_projet, idExtracteur)
}

/* Fonction permettant de lancer l'extraction de la table chronologie */
func (a *App) ExtractionChronologie() bool {
	err := extraction.ExtraireTableChronologie(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
		return false
	}
	return true
}

/***************************************************************************************/
/************************* Arborescence FUNCTIONS PAGE ********************************/
/***************************************************************************************/

/*
	Fonction qui renvoie les enfants d'un élément dans l'arborescence de la machine analysée

@param cheminDossier : le chemin du dossier duquel on veut connaître les enfants
@return : la liste des enfants
*/
func (a *App) ArborescenceMachineAnalysee(cheminDossier []int) []arborescence.MetaDonnees {
	res, err := arborescence.RecupEnfantsArbo(chemin_projet, cheminDossier)
	if err != nil {
		a.signalerErreur(err)
	}
	return res
}

func (a *App) ExtraireArborescence(avecModele bool) arborescence.Arborescence {
	var cheminModele = ""
	var err error
	if avecModele {
		cheminModele, err = runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
			Title:   "Choisissez le modele",
			Filters: []runtime.FileFilter{{DisplayName: "Modèles aqua", Pattern: "modele.aqua"}},
		})
		if err != nil || cheminModele == "" {
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.ErrorDialog,
				Message: "Vous devez choisir un fichier modele.aqua. \nSi vous n'avez pas de modèle, il est possible d'en créer un en vous rendant sur la page d'accueil.\nSi vous ne souhaitez pas utiliser de modèle, décochez l'option « Comparer l'arborescence avec celle d'un modèle. »",
			})
			return arborescence.Arborescence{}
		}
	}
	res, err := arborescence.ExtraireArborescence(chemin_projet, filepath.Dir(cheminModele))

	if err != nil {
		a.signalerErreur(err)
	}

	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "Extraction terminée",
		Message: "L'extraction de l'arborescence s'est terminée avec succès !",
	})
	return res
}

/***************************************************************************************/
/************************* DB_INFO PAGE ************************************************/
/***************************************************************************************/
func (a *App) Get_db_info() map[string]string {
	adb := aquabase.InitDB_Extraction(chemin_projet)
	return adb.GetAllTableNames()
}

func (a *App) Get_header_table(tableName string, limitJS string) []map[string]interface{} {
	limit, _ := strconv.Atoi(limitJS)
	adb := aquabase.InitDB_Extraction(chemin_projet)
	return adb.SelectAllFrom(tableName, limit)
}

/***************************************************************************************/
/************************* Detection FUNCTIONS PAGE ********************************/
/***************************************************************************************/

/* Fonction permettant d'obtenir une liste des règles de détection
 */
func (a *App) ListeReglesDetection(lancerRegles bool) map[string]map[string]int {
	regles, reglesEnErreur, err := detection.ListeReglesDetection(chemin_projet, lancerRegles)
	if err != nil {
		a.signalerErreur(err)
	}
	if len(reglesEnErreur) != 0 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Certaines règles contiennent des erreurs",
			Message: "Les règles suivantes ont renvoyé une erreur.\n - " + strings.Join(reglesEnErreur, "\n - ") + "\nNous vous conseillons de vérifier leur syntaxe.",
		})
	}
	return regles
}

func (a *App) InfosRegleDetection(nomRegle string) detection.Regle {
	regles, err := detection.DetailsRegleDetection(chemin_projet, nomRegle)
	if err != nil {
		a.signalerErreur(err)
	}
	return regles
}

func (a *App) CreationReglesDetection(json_rule string) bool {
	err := detection.NewDetectionRule(chemin_projet, json_rule)
	if err != nil {
		a.signalerErreur(err)
	}
	return true
}

func (a *App) Delete_rule(nomRegle string) {
	err := detection.SuppressionRegleDetection(chemin_projet, nomRegle)
	if err != nil {
		a.signalerErreur(err)
	}
}

func (a *App) ResultatRegleDetection(nomRegle string) int {
	resultat, err := detection.ResultatRegleDetection(chemin_projet, nomRegle)
	if err != nil {
		a.signalerErreur(err)
	}
	if resultat == 0 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Certaines règles contiennent des erreurs",
			Message: "L'exécution de cette règle a renvoyé une erreur.\nNous vous conseillons de vérifier sa syntaxe.",
		})
	}

	return resultat
}

func (a *App) ResultatsSQL(nomRegle string) []map[string]interface{} {
	resultat, err := detection.ResultatSQL(chemin_projet, nomRegle)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (a *App) StatutReglesDetection() []map[string]interface{} {
	resultat, err := detection.StatutReglesDetection(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

/***************************************************************************************/
/*************************** Chronologie FUNCTIONS PAGE ********************************/
/***************************************************************************************/

func (app *App) ValeursTableChronologie(debut int, taille int) []map[string]interface{} {
	return extraction.ValeursTableChronologie(chemin_projet, debut, taille)
}

func (app *App) ResultatRequeteSQLExtraction(requete string, debut int, taille int) []map[string]interface{} {
	requete = fmt.Sprintf("%s LIMIT %d OFFSET %d", requete, taille, debut)
	log.Println("[INFO] - Execution depuis JS de la requete ", requete)
	var base aquabase.Aquabase = *aquabase.InitDB_Extraction(chemin_projet)
	return base.ResultatRequeteSQL(requete)
}

func (app *App) TailleRequeteSQLExtraction(requete string) int {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(chemin_projet)
	return base.TailleRequeteSQL(requete)
}

func (app *App) GetListeTablesExtraction() []string {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(chemin_projet)
	return base.GetListeTablesDansBDD()
}

/***************************************************************************************/
/******************************* Rapport FUNCTIONS PAGE ********************************/
/***************************************************************************************/

func (app *App) AjouterPisteDansRapport(titre string, description string) {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	err := rprt.AjouterPiste(titre, description)
	if err != nil {
		app.signalerErreur(err)
	} else {
		runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Piste ajoutée avec succès",
			Message: "La piste « " + titre + " » a correctement été ajoutée !",
		})
	}
}

func (app *App) AjouterEtapeDansRapport(requeteSQL string, lignesAEnregistrer []map[string]interface{}, idPiste string, commentaire string) {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	err := rprt.AjouterEtape(idPiste, commentaire, requeteSQL, lignesAEnregistrer)
	if err != nil {
		app.signalerErreur(err)
	} else {
		runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Piste ajoutée avec succès",
			Message: "La table a été correctement enregistrée dans votre rapport !",
		})
	}
}

func (app *App) ListePistesRapport() []map[string]interface{} {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	return rprt.GetPistes()
}

func (app *App) ListeEtapesRapport(idPiste int) []rapport.EtapeAnalyse {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	return rprt.GetEtapesPiste(idPiste)
}

func (app *App) DonneesTableRapport(nomTable string) []map[string]interface{} {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	log.Println(nomTable)
	return rprt.GetDonnesTableSauvegardee(nomTable)
}
