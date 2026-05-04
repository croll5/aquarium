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



EXEMPLES EVENT :

{"level":"info","ts":"2026-05-04T20:46:46.337+0200","caller":"utilitaires/utilitaires.go:155","msg":"application_startup","event":"application.startup","attrs":{"oui_au_debug":true,"source":"backend"}}

# explication :
level: "info" : événement informatif.
ts: "2026-05-04T20:46:46.337+0200" : horodatage précis (4 mai 2026, 20:46:46.337, UTC+2).
caller: "utilitaires/utilitaires.go:155" : point d’écriture dans le logger central.
msg: "application_startup" : message libre fourni par le backend au démarrage.
event: "application.startup" : nom d’événement structuré.
attrs :
oui_au_debug: true : la session debug est active.
source: "backend" : événement émis côté Go.

---------

{"level":"info","ts":"2026-05-04T20:46:51.607+0200","caller":"utilitaires/utilitaires.go:155","msg":"[PARAMETRES] case_changee | option=contrastes | valeur=true","event":"parametres.case_changee","attrs":{"option":"contrastes","source":"frontend","utilisateur_session_debug":true,"valeur":true}}

# explication :
level: "info" : événement informatif.
ts: "2026-05-04T20:46:51.607+0200" : horodatage (environ 5.27 s après le démarrage).
caller: "utilitaires/utilitaires.go:155" : même point d’entrée central.
msg: "[PARAMETRES] case_changee | option=contrastes | valeur=true" : message manuel lisible.
event: "parametres.case_changee" : événement structuré “une case paramètres a changé”.
attrs :
option: "contrastes" : la case concernée.
valeur: true : nouvel état coché.
source: "frontend" : événement émis depuis JS via wrapper Wails.
utilisateur_session_debug: true : garde-fou de session debug actif.