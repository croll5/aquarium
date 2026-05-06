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

import { test, expect } from '@playwright/test';

// Ouvre l'onglet Parametres depuis l'iframe de la page d'accueil
// et retourne le frame locator prêt à être utilisé.
async function ouvrirPageParametres(page) {
  const accueilFrame = page.frameLocator('iframe[title="frame"]');
  // Depuis la page principale, ouvre la page Parametres via le bouton dedie.
  await accueilFrame.getByRole('button', { name: /Paramètres/i }).click();
  await expect(accueilFrame.locator('h1#titre_parametres')).toBeVisible();
  return accueilFrame;
}

// Simule la fin du parcours utilisateur apres sauvegarde:
// enregistrement, fermeture de la notification, puis retour accueil.
async function revenirAccueilDepuisParametres(frame) {
  await frame.getByRole('button', { name: 'Enregistrer' }).click();
  await expect(frame.locator('#toast_parametres')).toBeVisible();
  // Sur cette page, la fermeture de la notification confirme le flux et
  // declenche la redirection vers l'accueil quand l'enregistrement est un succes.
  await frame.getByRole('button', { name: /Fermer la notification/i }).click();
  await expect(frame.getByRole('heading', { name: /Bienvenue dans l’Aquarium/i })).toBeVisible();
}

// Normalise l'etat initial de la case "non_aux_bubulles" afin de rendre
// le test deterministe, meme si un etat precedent a ete persiste.
async function forcerEtatBubulles(page, checked) {
  const frame = await ouvrirPageParametres(page);
  const checkbox = frame.locator('#non_aux_bubulles');
  if ((await checkbox.isChecked()) !== checked) {
    if (checked) {
      await checkbox.check();
    } else {
      await checkbox.uncheck();
    }
  }
  await revenirAccueilDepuisParametres(frame);
}

// `test.describe` regroupe tout le parcours E2E lie aux parametres "bubulles".
// Le `beforeEach` definit un environnement Wails mocke commun a chaque execution
// du test pour garder des scenarios reproductibles et lisibles.
test.describe('Parametres - non_aux_bubulles', () => {
  test.beforeEach(async ({ page }) => {
    await page.addInitScript(() => {
      // Etat partage entre pages/iframes pour simuler une persistance Wails.
      const root = window.top;
      root.__wailsMockState = root.__wailsMockState ?? {
        Contrastes: false,
        Dyslexie: false,
        NonAuxBubulles: false,
      };

      window.go = {
        main: {
          App: {
            async GetParametres() {
              return { ...root.__wailsMockState };
            },
            async SauvegarderParametres(contrastes, dyslexie, nonAuxBubulles) {
              root.__wailsMockState.Contrastes = contrastes === true;
              root.__wailsMockState.Dyslexie = dyslexie === true;
              root.__wailsMockState.NonAuxBubulles = nonAuxBubulles === true;
              return true;
            },
          },
        },
      };
    });

    await page.goto('/');
  });

  test('Test 1.1 puis Test 1.2 - toggle et sauvegarde de "...bubulles"', async ({ page }) => {
    // Base line robuste: l'etat peut deja etre persiste par un precedent run.
    await forcerEtatBubulles(page, false);

    await test.step('Test 1.1: cocher "enlever les bubulles", enregistrer, fermer la page Parametres', async () => {
      const frame = await ouvrirPageParametres(page);
      const checkbox = frame.locator('#non_aux_bubulles');

      await expect(checkbox).not.toBeChecked();
      await checkbox.check();

      await revenirAccueilDepuisParametres(frame);
    });

    await test.step('Test 1.2: decocher "enlever les bubulles", enregistrer, fermer la page Parametres', async () => {
      const frame = await ouvrirPageParametres(page);
      const checkbox = frame.locator('#non_aux_bubulles');

      await expect(checkbox).toBeChecked();
      await checkbox.uncheck();

      await revenirAccueilDepuisParametres(frame);
    });
  });
});
