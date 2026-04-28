# Module `params` (fullstack)

Ce module gère la persistance des préférences utilisateur de l’application :

- `contrastes`
- `dyslexie`
- `non_aux_bubulles`

Les préférences sont stockées dans `./config/params.xml`.

## Format du fichier

Le fichier attendu est :

```xml
<?xml version="1.0" encoding="UTF-8"?>
<parametres>
  <contrastes>false</contrastes>
  <dyslexie>false</dyslexie>
  <non_aux_bubulles>true</non_aux_bubulles>
</parametres>
```

## Backend (`Go`)

Fichier : `modules/params/params.go`

- `type ParametresXML`
  - structure de mapping XML (`<parametres>...`).
- `SauvegarderParametres(cheminBase, contrastes, dyslexie, nonAuxBubulles)`
  - sérialise en XML,
  - crée `./config` si nécessaire,
  - écrit `./config/params.xml`.
- `ChargerParametres(cheminBase)`
  - lit et désérialise `./config/params.xml`,
  - si le fichier n’existe pas, renvoie des valeurs par défaut (`false/false/false`) sans erreur bloquante.

Intégration dans `app.go` :

- `App.SauvegarderParametres(...)` : API Wails appelée par le frontend pour enregistrer.
- `App.GetParametres()` : API Wails appelée au démarrage pour charger l’état courant.

## Frontend (`JS/HTML`)

### Chargement au démarrage

Fichier : `frontend/src/templates/index.js`

- `initialiser_parametres()` appelle `window.go.main.App.GetParametres()`.
- Les variables globales de la fenêtre principale sont mises à jour :
  - `contrastes`
  - `dyslexie`
  - `non_aux_bubulles`
- L’iframe est rechargée pour appliquer immédiatement le style/comportement sur la page affichée.

### Enregistrement depuis la page Paramètres

Fichier : `frontend/src/js/parametres.js`

- Le bouton **Enregistrer** appelle `quitter_parametres()`.
- Cette fonction lit l’état des 3 cases puis appelle `window.go.main.App.SauvegarderParametres(...)`.
- Une toast notification indique le résultat :
  - succès : message + redirection vers l’accueil,
  - échec : message d’erreur, sans redirection forcée.

### Fermeture sans enregistrement

La croix de fermeture (`parametres.html`) appelle `fermer_parametres_sans_enregistrer()` :

- retour à l’accueil,
- aucune écriture dans `params.xml`.

## Résumé du flux

1. Démarrage app : frontend appelle `GetParametres()`.
2. Le backend lit `config/params.xml` (si présent) et renvoie l’état.
3. Le frontend applique l’état global (contraste, police, bulles).
4. L’utilisateur modifie les options dans l’écran Paramètres.
5. Si **Enregistrer** : appel backend `SauvegarderParametres(...)` puis écriture XML.
