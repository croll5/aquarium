# Tests frontend

Aller voir : 
https://v3alpha.wails.io/guides/e2e-testing/

## Pourquoi cette stack

- `Vitest` : tests unitaires JS rapides pour valider la logique (fonctions, transformations, états UI simples).
- `Playwright` : tests E2E pour valider les parcours complets utilisateur dans un navigateur réel.
- `Mock Wails (window.go)` : permet de tester le frontend sans dépendre du backend Go en simulant les appels `window.go.main.App.*`.

## Limites
On ne teste PAS :
- la vraie app desktop native
- la WebView directement
On teste via  :
- la version web exposée par Wails
C'est poiur cela qu'il faut 2 session de Terminal actives.

## Setup minimal

1. Installer les dépendances JS de test dans le workspace frontend (Vitest + Playwright).
> npm install -D vitest jsdom

> npm install -D @playwright/test  
 
2. Initialize les browsers Firefox + WebKit

> npx playwright install  

exemple pour macOs:
    
    Downloading Firefox 148.0.2 (playwright firefox v1511) from https://cdn.playwright.dev/dbazure/download/playwright/builds/firefox/1511/firefox-mac-arm64.zip
    97.1 MiB [====================] 100% 0.0s
    
    Downloading WebKit 26.4 (playwright webkit v2272) from https://cdn.playwright.dev/dbazure/download/playwright/builds/webkit/2272/webkit-mac-15-arm64.zip
    75.4 MiB [====================] 100% 0.0s
    


3. Configurer Vitest en environnement `jsdom`.
4. Créer un mock partagé de `window.go` (ex: `test/mocks/wails.mock.js`).
5. Charger ce mock dans les tests unitaires (setup Vitest) et E2E (injection Playwright avant chargement de page).

## Règles d'usage

- Unitaire (`Vitest`) : tester la logique isolée, sans backend réel.
- E2E (`Playwright`) : tester les scénarios critiques (navigation, sauvegarde paramètres, erreurs visibles).
- Toujours contrôler les réponses Wails via mock (`succès`, `échec`, `exception`) pour couvrir les cas robustes.

## lancer les tests JS unitaires
npx vitest run --config frontend/tests/vitest.config.mjs

## Lancer les tests Playwright (E2E)
1. Démarrer l'application dans un terminal :
`wails dev`
2. Dans un second terminal, lancer les tests E2E :
`npx playwright test --config frontend/tests/setup/playwright.config.js`

Commandes utiles :
- Un seul fichier :
`npx playwright test --config frontend/tests/setup/playwright.config.js frontend/tests/e2e/parametres.non_aux_bubulles.test.spec.js`
- Mode UI :
`npx playwright test --ui --config frontend/tests/setup/playwright.config.js`

Erreur fréquente :
- `ERR_CONNECTION_REFUSED` : vérifier que `wails dev` est bien lancé et que le port de `baseURL` dans `frontend/tests/setup/playwright.config.js` correspond au port réel.
