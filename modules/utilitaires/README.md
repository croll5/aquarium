# Module `utilitaires`

Ce module contient des utilitaires techniques partages.

## Installation des dependances logger

Pour telecharger les librairies utilisees par le logger central :

```bash
go get go.uber.org/zap gopkg.in/natefinch/lumberjack.v2
```

## Logger central (zap + lumberjack)

Le logger est centralise dans `modules/utilitaires/utilitaires.go` :

- Initialisation unique : `InitLogger(cheminBase, debugActif)`
- Ecriture : `LogEvent(niveau, evenement, attributs, message)`
- Fin de session : `SyncLogger()`
- Helpers applicatifs optionnels : `modules/utilitaires/app_logging.go` (ex: flux `parametres.enregistrer`)

Un fichier de log est cree dans `./logs/` a chaque lancement quand le debug est actif, avec un nom horodate :
`aquarium_YYYYMMDD_HHMMSS.log`.
La rotation est geree par `lumberjack` (taille, retention, compression).

## Pourquoi zap + lumberjack

- `zap` : logger rapide, structure (JSON), niveaux natifs (`Debug`, `Info`, `Warn`, `Error`) et integration propre avec les metadonnees (champs).
- `lumberjack` : rotation simple des fichiers (`MaxSize`, `MaxBackups`, `MaxAge`, `Compress`) pour eviter la croissance infinie de `aquarium.log`.
- combinaison des deux : journalisation performante + gestion de cycle de vie des fichiers sans logique maison.

## Activation

L'activation du logger depend de `oui_au_debug` dans `config/params.xml`.

Au lancement de l'application (`app.startup`) :

1. lecture des parametres via `params.ChargerParametres(...)`,
2. appel `InitLogger(..., parametres.OuiAuDebug)`.

Si `oui_au_debug` vaut `false`, le logger reste inactif (no-op).

## Wrappers exposes au frontend

Le frontend n'appelle pas directement `utilitaires`.
Les points d'entree Wails sont dans `app.go` :

- `App.LoggerEvent(niveau, evenement, attributs, message) bool`

Ainsi, front et back passent par le meme service central.

### API recommandee : LoggerEvent (mode hybride)

`LoggerEvent` combine :

- un evenement structure (`evenement`) ;
- des attributs libres (`attributs`) ;
- un message manuel optionnel (`message`).

Exemple :

```js
app.LoggerEvent(
  "info",
  "parametres.case_changee",
  { option: "contrastes", valeur: true, source: "frontend" },
  "[PARAMETRES] case_changee | option=contrastes | valeur=true"
);
```

## Activer des evenements de log (JS et Go)

Prerequis : la session doit etre en debug (`oui_au_debug=true` dans `config/params.xml`, puis redemarrage de l'application).

### Exemple JS (frontend)

Depuis une page frontend, appeler `LoggerEvent` :

```js
const app = parent?.window?.go?.main?.App;
if (app && typeof app.LoggerEvent === "function" && parent.oui_au_debug === true) {
  const valeur = document.getElementById("contrastes").checked;
  app.LoggerEvent(
    "info",
    "parametres.case_changee",
    { option: "contrastes", valeur, source: "frontend", utilisateur_session_debug: true },
    "[PARAMETRES] case_changee | option=contrastes | valeur=" + (valeur ? "true" : "false")
  );
}
```

Points importants :

- verifier `parent.oui_au_debug === true` avant emission ;
- preferer `LoggerEvent` pour tous les evenements metier.
- pour les cases a cocher, preferer un evenement unique `case_changee` avec `valeur=true|false`.

### Exemple Go (backend)

Dans n'importe quel module Go, utiliser directement le service central :

```go
import "aquarium/modules/utilitaires"

func ExempleTraitement() {
	utilitaires.LogEvent("info", "traitement.debut", map[string]interface{}{
		"source": "backend",
		"id": 42,
	}, "Demarrage du traitement 42")
}
```

Le logger est no-op si la session n'est pas en debug, donc ces appels restent sans effet quand `oui_au_debug=false`.

## Politique de journalisation

Le logging n'est pas systematique sur toute l'application.

- logger uniquement les evenements a valeur operationnelle (debut/succes/echec d'actions metier, erreurs, transitions d'etat utiles) ;
- eviter le bruit (events ultra frequents sans valeur de diagnostic) ;
- conserver `signalerErreur(err)` pour les fichiers d'erreurs detailles, et `LogEvent(...)` pour la tracabilite generale.
