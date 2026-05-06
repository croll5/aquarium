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

// Ce test charge parametres.js dans jsdom et verifie que l'etat de la case
// `#non_aux_bubulles` est bien synchronise vers `parent.non_aux_bubulles`.
const scriptPath = path.resolve(process.cwd(), 'frontend/src/js/parametres.js');
const scriptContent = fs.readFileSync(scriptPath, 'utf8');

function chargerScriptParametres() {
  vi.stubGlobal('contrastes_eleves', vi.fn());
  vi.stubGlobal('police_dyslexie', vi.fn());
  window.eval(scriptContent);
}

describe('parametres.js - enlever_bubulles', () => {
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
        go: { main: { App: {} } },
      },
    };

    chargerScriptParametres();
  });

  it('synchronise parent.non_aux_bubulles a true quand la checkbox est cochee', () => {
    const checkbox = document.getElementById('non_aux_bubulles');
    checkbox.checked = true;

    enlever_bubulles();

    expect(parent.non_aux_bubulles).toBe(true);
  });

  it('synchronise parent.non_aux_bubulles a false quand la checkbox est decochee', () => {
    const checkbox = document.getElementById('non_aux_bubulles');
    checkbox.checked = false;
    parent.non_aux_bubulles = true;

    enlever_bubulles();

    expect(parent.non_aux_bubulles).toBe(false);
  });
});
