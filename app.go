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
	"aquarium/modules/config"
	"aquarium/modules/detection"
	"aquarium/modules/extraction"
	"aquarium/modules/gestionprojet"
	"aquarium/modules/params"
	"aquarium/modules/rapport"
	"aquarium/modules/utilitaires"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var chemin_projet string

const DOSSIER_ERREURS = "erreurs"

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
	runtime.WindowMaximise(ctx)
	a.ctx = ctx

	cheminBase, err := utilitaires.GetCheminBaseApplication()
	if err != nil {
		log.Printf("initialisation logger impossible (chemin base): %v", err)
		return
	}
	parametres, err := params.ChargerParametres(cheminBase)
	if err != nil {
		log.Printf("initialisation logger impossible (lecture params): %v", err)
		return
	}
	if err = utilitaires.InitLogger(cheminBase, parametres.OuiAuDebug); err != nil {
		log.Printf("initialisation logger impossible: %v", err)
		return
	}
	utilitaires.LogEvent("info", "application.startup", map[string]interface{}{
		"source":       "backend",
		"oui_au_debug": parametres.OuiAuDebug,
	}, "application_startup")
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	a.ctx = ctx
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
	utilitaires.LogEvent("info", "application.shutdown", map[string]interface{}{
		"source": "backend",
	}, "application_shutdown")
	utilitaires.SyncLogger()
}

// Call this function when a bug appear
func (a *App) signalerErreur(erreur error) {
	utilitaires.LogEvent("error", "application.erreur", map[string]interface{}{
		"source": "backend",
		"erreur": fmt.Sprintf("%+v", erreur),
	}, "erreur_application")
	// Récupérer l’horodatage
	nomFichierErreur := time.Now().Format("2006010215040599999_erreur.txt")
	// Informer l’utilisateur de l’erreur
	runtime.WindowExecJS(a.ctx, fmt.Sprintf("signaler_erreur('%s','%s')", nomFichierErreur, url.QueryEscape(fmt.Sprintf("%+v", erreur))))
	// Récupérer le chemin d'enregistrement de l’erreur
	cheminErreurs := filepath.Join(chemin_projet, DOSSIER_ERREURS)
	if chemin_projet == "" {
		cheminExecutable, err := os.Executable()
		if err != nil {
			a.alerterEnregistrementErreurImpossible()
			return
		}
		cheminErreurs = filepath.Join(filepath.Dir(cheminExecutable), DOSSIER_ERREURS)
	}
	// Créer le dossier des erreur s’il n’existe pas
	err := os.MkdirAll(cheminErreurs, 0o755)
	if err != nil {
		a.alerterEnregistrementErreurImpossible()
		return
	}
	// Enregistrer le contenu de l’erreur
	fichierErr, err := os.Create(filepath.Join(cheminErreurs, nomFichierErreur))
	if err != nil {
		a.alerterEnregistrementErreurImpossible()
		return
	}
	_, err = fichierErr.Write([]byte(fmt.Sprintf("%+v", erreur)))
	if err != nil {
		a.alerterEnregistrementErreurImpossible()
	}
}

func (a *App) alerterEnregistrementErreurImpossible() {
	runtime.WindowExecJS(a.ctx, "details_erreur()")
}

func (a *App) LoggerEvent(niveau string, evenement string, attributs map[string]interface{}, message string) bool {
	utilitaires.LogEvent(niveau, evenement, attributs, message)
	return true
}

func (a *App) SauvegarderParametres(contrastes bool, dyslexie bool, nonAuxBubulles bool, ouiAuDebug bool) bool {

	utilitaires.LogParametresEnregistrerDebut(contrastes, dyslexie, nonAuxBubulles, ouiAuDebug)

	cheminBase, err := utilitaires.GetCheminBaseApplication()
	if err != nil {
		utilitaires.LogParametresEnregistrerEchec("GetCheminBaseApplication", err.Error())
		a.signalerErreur(err)
		return false
	}
	err = params.SauvegarderParametres(cheminBase, contrastes, dyslexie, nonAuxBubulles, ouiAuDebug)
	if err != nil {
		utilitaires.LogParametresEnregistrerEchec("params.SauvegarderParametres", err.Error())
		a.signalerErreur(err)
		return false
	}
	utilitaires.LogParametresEnregistrerSucces()
	return true
}

func (a *App) GetParametres() params.ParametresXML {
	cheminBase, err := utilitaires.GetCheminBaseApplication()
	if err != nil {
		a.signalerErreur(err)
		return params.ParametresXML{}
	}
	parametres, err := params.ChargerParametres(cheminBase)
	if err != nil {
		a.signalerErreur(err)
		return params.ParametresXML{}
	}
	return parametres
}

/***************************************************************************************/
/************************* INDEX FUNCTIONS **********************************/
/***************************************************************************************/

/*
	Cette fonction permet l'ouverture d'une analyse aquarium, à partir d'un fichier .aqua

@return : vrai si et seulement si l'analyse a été correctement ouverte
*/
func (a *App) OuvrirAnalyseExistante() bool {
	options := runtime.OpenDialogOptions{
		Title: "Ouvrir une analyse existante",
	}
	if goruntime.GOOS != "darwin" {
		options.Filters = []runtime.FileFilter{{DisplayName: "Aquarium", Pattern: "*.aqua"}}
	}
	fichier, err := runtime.OpenFileDialog(a.ctx, options)
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
	if !strings.EqualFold(filepath.Base(fichier), "analyse.aqua") {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Fichier invalide",
			Message: "Veuillez sélectionner un fichier nommé analyse.aqua",
		})
		return false
	}
	chemin_projet = filepath.Dir(fichier)
	extraction.CreationBaseAnalyse(chemin_projet)
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

func (a *App) ListeConfigurationsDisponibles() []string {
	resultat, err := config.GetListeConfigurationsDisponibles(false)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (a *App) ListeConfigExtractionsDisponibles() []string {
	resultat, err := config.GetListeConfigurationsDisponibles(true)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (a *App) VerifierPrerequisDemarrage() []string {
	manquants, err := config.VerifierPrerequisDemarrage()
	if err != nil {
		a.signalerErreur(err)
		return []string{}
	}
	return manquants
}

/*
	Fonction permettant la création d'un nouveau projet

@return : le chemin vers le nouveau projet
*/
func (a *App) CreationNouveauProjet(configuration config.AquaConfig) string {
	// Création de l’arborescence de l’analyse
	chemin_projet = configuration.DossierAnalyse
	configuration.DebutAnalyse = time.Now()
	err := gestionprojet.CreationArborescence(&chemin_projet, configuration)
	if err != nil {
		a.signalerErreur(err)
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

/***************************************************************************************/
/************************* Extraction FUNCTIONS PAGE **********************************/
/***************************************************************************************/

func (a *App) ListeMachinesAnalysees() map[string]config.AquaConfigMachine {
	liste, err := config.ListeMachinesAnalysees(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
	}
	return liste
}

/* Fonction renvoyant la liste des éléments pouvant être extraits de l'ORC
 */
func (a *App) ListeExtractionsPossibles() map[string]extraction.ExtractionMachine {
	extractions, err := extraction.ListeExtractionsHtml(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
	}
	return extractions
}

func (a *App) LancerExtraction(ordreExtractions []map[string]string) {
	erreursConfig, err := extraction.LancerExtractions(chemin_projet, ordreExtractions)
	if err != nil {
		a.signalerErreur(err)
	}
	for _, erreurConf := range erreursConfig {
		if erreurConf.ConfigInexistante {
			a.signalerErreur(fmt.Errorf("La configuration de l’extraction %s est introuvable", erreurConf.IdExtraction))
		} else if len(erreurConf.DoublonTable) > 0 {
			for _, nomTable := range erreurConf.DoublonTable {
				a.signalerErreur(fmt.Errorf("La table %s est utilisée dans deux configurations différentes (%s)", nomTable, erreurConf.IdExtraction))
			}
		} else if len(erreurConf.ParametresManquants) > 0 {
			a.signalerErreur(fmt.Errorf("Des paramètres de la configuration « %s » sont manquants : %v", erreurConf.IdExtraction, erreurConf.ParametresManquants))
		}
	}
}

/* Fonction permettant de connaitre le pourcentage de progression d'une extraction*/
func (a *App) ProgressionExtraction() map[string]string {
	return extraction.ProgressionExtraction(chemin_projet)
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
func (a *App) ArborescenceMachineAnalysee(cheminDossier []string, idMachine string) []arborescence.MetaDonnees {
	res, err := arborescence.RecupEnfantsArbo(chemin_projet, cheminDossier, idMachine)
	if err != nil {
		a.signalerErreur(err)
	}
	return res
}

func (a *App) ExtractionEnCours() bool {
	return arborescence.ExtractionEnCours()
}

func (a *App) ArborescenceEnCache() string {
	return arborescence.ArborescenceEnCache()
}

func (a *App) DetailsFichierArborescence(idFichier int64, idMachine string) []map[string]interface{} {
	resultat, err := arborescence.DetailsFichier(chemin_projet, idFichier, idMachine)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

/***************************************************************************************/
/************************* DB_INFO PAGE ************************************************/
/***************************************************************************************/
func (a *App) Get_db_info() map[string]string {
	adb := aquabase.InitDB_Extraction(chemin_projet)
	resultat, err := adb.GetAllTableNames()
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (a *App) Get_header_table(tableName string, limitJS string) []map[string]interface{} {
	limit, _ := strconv.Atoi(limitJS)
	adb := aquabase.InitDB_Extraction(chemin_projet)
	resultat, err := adb.SelectAllFrom(tableName, limit)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
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

func (a *App) ResultatRequeteSQLExtraction(requete string, debut int, taille int) []map[string]interface{} {
	requete = strings.TrimSpace(requete)
	for strings.HasSuffix(requete, ";") {
		requete = strings.TrimSpace(strings.TrimSuffix(requete, ";"))
	}
	requete = fmt.Sprintf("SELECT * FROM (%s) AS requete_utilisateur LIMIT %d OFFSET %d", requete, taille, debut)
	log.Println("[INFO] - Execution depuis JS de la requete ", requete)
	var base aquabase.Aquabase = *aquabase.InitDB_Extraction(chemin_projet)
	resultat, err := base.ResultatRequeteSQL(requete)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (app *App) TailleRequeteSQLExtraction(requete string) int {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(chemin_projet)
	return base.TailleRequeteSQL(requete)
}

func (a *App) GetListeTablesExtraction() []string {
	var base *aquabase.Aquabase = aquabase.InitDB_Extraction(chemin_projet)
	resultat, err := base.GetListeTablesDansBDD()
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
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

func (a *App) ListePistesRapport() []map[string]interface{} {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	resultat, err := rprt.GetPistes()
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (a *App) ListeEtapesRapport(idPiste int) []rapport.EtapeAnalyse {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	resultat, err := rprt.GetEtapesPiste(idPiste)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (a *App) DonneesTableRapport(nomTable string) []map[string]interface{} {
	var rprt *rapport.Rapport = rapport.InitRapport(chemin_projet)
	log.Println(nomTable)
	resultat, err := rprt.GetDonnesTableSauvegardee(nomTable)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

/***************************************************************************************/
/******************************* FONCTIONS UTILITAIRES  ********************************/
/***************************************************************************************/

func (app *App) ChoisirFichier(ordre string) []string {
	chemin, err := runtime.OpenMultipleFilesDialog(app.ctx, runtime.OpenDialogOptions{
		Title: ordre})
	if err != nil {
		runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
			Title: "Erreur lors de l'ouverture du fichier",
			Type:  runtime.ErrorDialog,
		})
	}
	return chemin
}

func (app *App) ChoisirDossier(ordre string) string {
	chemin, err := runtime.OpenDirectoryDialog(app.ctx, runtime.OpenDialogOptions{
		Title: ordre})
	if err != nil {
		runtime.MessageDialog(app.ctx, runtime.MessageDialogOptions{
			Title: "Erreur lors de l'ouverture du dossier",
			Type:  runtime.ErrorDialog,
		})
	}
	return chemin
}

/*   FONCTIONS VUE D’ENSEMBLE */

func (a *App) GetAquaConfig() config.AquaConfig {
	aquaConfig, err := config.GetAquaConfig(chemin_projet)
	if err != nil {
		a.signalerErreur(err)
	}
	return aquaConfig
}

/* ---------- FONCTIONS CONFIGURATION DES EXTRACTIONS ---------- */

func (a *App) ListeFichiersAnalysables(idMachine string) extraction.DossierAnalysable {
	resultat, err := extraction.ListeFichiersAnalysables(filepath.Join(chemin_projet, config.DOSSIER_FICHIERS_A_ANALYSER), idMachine)
	if err != nil {
		a.signalerErreur(err)
	}
	return resultat
}

func (a *App) CorrespondanceCheminModele(filename string, modele string) bool {
	result, err := filepath.Match(modele, filename)
	if err != nil {
		a.signalerErreur(err)
	}
	return result
}
