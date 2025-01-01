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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var chemin_projet string
var chemin_bdd string

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

/****************************************************************************/
/************************* TEST FUNCTIONS **********************************/
/****************************************************************************/

// BlancPageFunction returns a greeting for the given name
func (a *App) BlancPageFunction(text string) string {
	return fmt.Sprintf("Hello %s, It's free1 time!", text)
}

/***************************************************************************************/
/************************* INDEX FUNCTIONS **********************************/
/***************************************************************************************/

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
	chemin_bdd = filepath.Join(chemin_projet, "analyse", "extractions.db")
	return true
}

func (a *App) CreationDossierNouveauModele() string {
	// Partie création du squelette de l'analyse
	projet, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choix de l'emplacement du modèle"})
	if err != nil {
		return ""
	}
	chemin_projet = projet
	chemin_bdd = filepath.Join(chemin_projet, "analyse", "extractions.db")
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

func (a *App) CreationNouveauProjet() string {
	// Partie création du squelette de l'analyse
	projet, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choix de l'emplacement de l'analyse"})
	if err != nil {
		return ""
	}
	chemin_projet = projet
	chemin_bdd = filepath.Join(chemin_projet, "analyse", "extractions.db")
	if !gestionprojet.CreationArborescence(chemin_projet) {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Problème dans la création de l'analyse",
			Message: "Les fichiers d'analyse n'ont pas pu être créés. Vérifiez que le dossier sélectionné est vide et que vous avez les droits en écriture :/"})
		return ""
	}
	return chemin_projet
}

// Greet returns a greeting for the given name
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

func (a *App) ValidationCreationProjet(nomAnalyste string, description string) bool {
	err := gestionprojet.EcritureFichierAqua(nomAnalyste, description, time.Now(), time.Time{}, chemin_projet)
	if err != nil {
		log.Println("ERR | Problème dans l'écriture du fichier .aqua : ", err.Error())
		a.signalerErreur(err)
		return false
	}
	return true
}

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

func (a *App) ListeExtractionsPossibles() map[string]extraction.InfosExtracteur {
	resultat, err := extraction.ListeExtracteursHtml(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

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

func (a *App) AnnulerExtraction(module string) bool {
	return extraction.AnnulerExtraction(module)
}

func (a *App) ProgressionExtraction(idExtracteur string) float32 {
	return extraction.ProgressionExtraction(chemin_projet, idExtracteur)
}

/***************************************************************************************/
/************************* Arborescence FUNCTIONS PAGE ********************************/
/***************************************************************************************/

func (a *App) ArborescenceMachineAnalysee(cheminDossier []int) []arborescence.MetaDonnees {
	res, err := arborescence.RecupEnfantsArbo(chemin_projet, cheminDossier)
	if err != nil {
		a.signalerErreur(err)
	}
	return res
}

/***************************************************************************************/
/************************* DB_INFO PAGE ************************************************/
/***************************************************************************************/
func (a *App) Get_db_info() map[string]string {
	adb := aquabase.Init(chemin_bdd)
	return adb.GetAllTableNames()
}

func (a *App) Get_header_table(tableName string, limitJS string) []map[string]interface{} {
	limit, _ := strconv.Atoi(limitJS)
	adb := aquabase.Init(chemin_bdd)
	return adb.SelectAllFrom(tableName, limit)
}

/***************************************************************************************/
/************************* ??? ********************************/
/***************************************************************************************/

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
/************************* Detection FUNCTIONS PAGE ********************************/
/***************************************************************************************/

func (a *App) ListeReglesDetection(lancerRegles bool) map[string]int {
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

func (a *App) CreationReglesDetection(json_rule string) bool {
	chemin_regles := filepath.Join(chemin_projet, "regles_detection")
	// Conversion de la chaîne JSON en une structure Go
	var regle map[string]interface{}
	if err := json.Unmarshal([]byte(json_rule), &regle); err != nil {
		a.signalerErreur(err)
	}
	// Récupération du nom à partir du JSON
	nom, ok := regle["nom"].(string)
	if !ok {
		a.signalerErreur(fmt.Errorf("Json without the variable: nom"))
	}
	// Conversion de la structure Go en JSON formaté
	data, err := json.MarshalIndent(regle, "", "  ")
	if err != nil {
		a.signalerErreur(err)
	}
	// Création du chemin complet du fichier avec le nom du JSON
	chemin_complet := filepath.Join(chemin_regles, nom+".json")
	// Écriture des données JSON dans un fichier
	if err := os.WriteFile(chemin_complet, data, 0644); err != nil {
		a.signalerErreur(err)
	}
	return true
}

func (a *App) ResultatsSQL(nomRegle string) []map[string]interface{} {
	chemin_regles := filepath.Join(chemin_projet, "regles_detection", nomRegle+".json")
	resultat, err := detection.ResultatSQL(chemin_projet, chemin_regles, nomRegle)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}
