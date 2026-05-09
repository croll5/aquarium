package utilitaires

// LogParametresEnregistrerDebut journalise le debut de la sauvegarde des parametres.
func LogParametresEnregistrerDebut(contrastes bool, dyslexie bool, nonAuxBubulles bool, ouiAuDebug bool) {
	LogEvent("info", "parametres.enregistrer", map[string]interface{}{
		"status":           "debut",
		"contrastes":       contrastes,
		"dyslexie":         dyslexie,
		"non_aux_bubulles": nonAuxBubulles,
		"oui_au_debug":     ouiAuDebug,
	}, "Sauvegarde des parametres - debut")
}

// LogParametresEnregistrerEchec journalise un echec de la sauvegarde des parametres.
func LogParametresEnregistrerEchec(etape string, erreur string) {
	LogEvent("error", "parametres.enregistrer", map[string]interface{}{
		"status": "echec",
		"etape":  etape,
		"erreur": erreur,
	}, "Sauvegarde des parametres - echec")
}

// LogParametresEnregistrerSucces journalise le succes de la sauvegarde des parametres.
func LogParametresEnregistrerSucces() {
	LogEvent("info", "parametres.enregistrer", map[string]interface{}{
		"status": "succes",
	}, "Sauvegarde des parametres - succes")
}
