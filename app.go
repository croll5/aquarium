package main

import (
	"aquarium/modules/arborescence"
	"aquarium/modules/extraction"
	"aquarium/modules/gestionprojet"
	"context"
	"log"
	"path/filepath"
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

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
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
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) signalerErreur(erreur error) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Erreur dans l'écriture du projet",
		Message: "Félicitation ! Vous venez de trouver un bogue dans le logiciel Aquarium !\n C'est cadeau : \n " + erreur.Error(),
	})
}

func (a *App) CreationNouveauProjet() string {
	// Partie création du squelette de l'analyse
	projet, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choix de l'emplacement de l'analyse"})
	if err != nil {
		return ""
	}
	chemin_projet = projet
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

func (a *App) OuvrirAnalyseExistante() bool {
	fichier, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Ouvrir une analyse existante",
		Filters: []runtime.FileFilter{{DisplayName: "Aquarium", Pattern: "*.aqua"}},
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

func (a *App) ListeExtractionsPossibles() map[string]string {
	return extraction.ListeExtracteursHtml()
}

func (a *App) ExtraireElements(module string, description string) {
	err := extraction.Extraction(module)
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

func (a *App) ArborescenceMachineAnalysee() arborescence.Arborescence {
	res, err := arborescence.GetArborescence(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
	}
	return res
}

func (a *App) ExtraireArborescence() arborescence.Arborescence {
	res, err := arborescence.ExtraireArborescence(chemin_projet)
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
