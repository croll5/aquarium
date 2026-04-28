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

import fs from 'node:fs';
import path from 'node:path';
import { describe, it, expect, beforeEach, vi } from 'vitest';

const scriptPath = path.resolve(process.cwd(), 'frontend/src/js/parametres.js');
const scriptContent = fs.readFileSync(scriptPath, 'utf8');
// Laisse le temps aux .then/.catch de se resoudre avant les assertions.
const flushPromises = () => new Promise((resolve) => setTimeout(resolve, 0));

function chargerScriptParametres() {
  // Mock des fonctions globales utilisees par parametres.js.
  vi.stubGlobal('contrastes_eleves', vi.fn());
  vi.stubGlobal('police_dyslexie', vi.fn());
  // Charge le script tel qu'il tourne dans le navigateur.
  window.eval(scriptContent);
}

describe('parametres.js - quitter_parametres', () => {
  beforeEach(() => {
    vi.restoreAllMocks();

    document.body.innerHTML = `
      <div id="toast_parametres"></div>
      <div id="message_toast_parametres"></div>
      <input id="contrastes" type="checkbox" />
      <input id="dyslexie" type="checkbox" />
      <input id="non_aux_bubulles" type="checkbox" />
    `;

    global.parent = {
      contrastes: false,
      dyslexie: false,
      non_aux_bubulles: false,
      window: {
        go: {
          main: {
            App: {
              SauvegarderParametres: vi.fn(),
            },
          },
        },
      },
    };

    chargerScriptParametres();
  });

  it('envoie les 3 bons booleens a SauvegarderParametres', () => {
    // Contrat principal frontend -> backend: ordre et valeurs des 3 flags.
    document.getElementById('contrastes').checked = true;
    document.getElementById('dyslexie').checked = false;
    document.getElementById('non_aux_bubulles').checked = true;

    parent.window.go.main.App.SauvegarderParametres.mockResolvedValue(true);

    quitter_parametres();

    expect(parent.window.go.main.App.SauvegarderParametres).toHaveBeenCalledTimes(1);
    expect(parent.window.go.main.App.SauvegarderParametres).toHaveBeenCalledWith(true, false, true);
  });

  it('gere le retour backend true', async () => {
    // Si sauvegarde OK, on attend un feedback succes avec redirection.
    const toastSpy = vi.fn();
    afficher_toast_parametres = toastSpy;

    parent.window.go.main.App.SauvegarderParametres.mockResolvedValue(true);

    quitter_parametres();
    await flushPromises();

    expect(toastSpy).toHaveBeenCalledWith('Les paramètres ont bien été enregistrés.', true, 'succes');
  });

  it('gere le retour backend false', async () => {
    // Si backend renvoie false, le flux doit rester en echec non bloquant.
    const toastSpy = vi.fn();
    afficher_toast_parametres = toastSpy;

    parent.window.go.main.App.SauvegarderParametres.mockResolvedValue(false);

    quitter_parametres();
    await flushPromises();

    expect(toastSpy).toHaveBeenCalledWith("L'enregistrement des paramètres a échoué.", false, 'echec');
  });

  it('gere une exception backend', async () => {
    // En cas d'exception, le comportement attendu est identique a un echec.
    const toastSpy = vi.fn();
    afficher_toast_parametres = toastSpy;

    parent.window.go.main.App.SauvegarderParametres.mockRejectedValue(new Error('backend indisponible'));

    quitter_parametres();
    await flushPromises();

    expect(toastSpy).toHaveBeenCalledWith("L'enregistrement des paramètres a échoué.", false, 'echec');
  });
});
