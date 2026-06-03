/*
Copyright ou © ou Copr. Didier Hoizé, (10 mai 2026)

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

const nouvelleAnalyseScriptPath = path.resolve(process.cwd(), 'frontend/src/js/nouvelle_analyse.js');
const nouvelleAnalyseScript = fs.readFileSync(nouvelleAnalyseScriptPath, 'utf8');

const nouvelleAnalyseVisScriptPath = path.resolve(process.cwd(), 'frontend/src/js/nouvelle_analyse_vis.js');
const nouvelleAnalyseVisScript = fs.readFileSync(nouvelleAnalyseVisScriptPath, 'utf8');

const nouvelleAnalyseHtmlPath = path.resolve(process.cwd(), 'frontend/src/html/nouvelle_analyse.html');
const nouvelleAnalyseHtml = fs.readFileSync(nouvelleAnalyseHtmlPath, 'utf8');

// Attend la resolution des .then() utilises par les appels Wails mockes.
const flushPromises = () => new Promise((resolve) => setTimeout(resolve, 0));

// Mock commun de parent.window.go.main.App pour les scripts frontend testes.
function setParentMock() {
  global.parent = {
    window: {
      go: {
        main: {
          App: {
            ListeConfigurationsDisponibles: vi.fn().mockResolvedValue([]),
            ChoisirFichier: vi.fn().mockResolvedValue([]),
          },
        },
      },
    },
  };
}

// Remet le DOM a un etat connu avant chaque scenario.
function resetDom(html = '') {
  document.body.innerHTML = html;
}

// Charge le script frontend tel qu'il est execute dans la page.
function loadNouvelleAnalyseScript() {
  window.eval(nouvelleAnalyseScript);
}

// Charge le vrai HTML de la page pour valider le couplage DOM/JS.
function loadNouvelleAnalyseHtml() {
  const parser = new DOMParser();
  const doc = parser.parseFromString(nouvelleAnalyseHtml, 'text/html');
  document.documentElement.innerHTML = doc.documentElement.innerHTML;
}

// Mock minimal de vis.js requis par nouvelle_analyse_vis.js.
function installVisMock() {
  class DataSet {
    constructor(initial = []) {
      this.items = initial;
    }
  }

  class Network {
    constructor() {
      this.body = { nodes: {}, edges: {} };
      global.__networkInstance = this;
    }
  }

  vi.stubGlobal('vis', { DataSet, Network });
}

// Charge la couche JS visuelle (reseau, popup, interactions).
function loadNouvelleAnalyseVisScript() {
  window.eval(nouvelleAnalyseVisScript);
}

describe('nouvelle_analyse.js - priorite haute', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    vi.spyOn(console, 'log').mockImplementation(() => {});
    resetDom();
    setParentMock();
    loadNouvelleAnalyseScript();
  });

  describe('donnees_conf_analyse', () => {
    it('extrait les champs aqua_champ non vides et ignore les champs vides', () => {
      resetDom(`
        <input aqua_champ="Nom" value="Analyse Demo" />
        <textarea aqua_champ="Description">Description test</textarea>
        <input aqua_champ="Vide" value="" />
        <select aqua_champ="Priorite">
          <option value="haute" selected>Haute</option>
        </select>
      `);

      const resultat = donnees_conf_analyse();

      expect(resultat).toEqual({
        Nom: 'Analyse Demo',
        Description: 'Description test',
        Priorite: 'haute',
      });
      expect(resultat.Vide).toBeUndefined();
    });

    it('respecte aqua_prof pour filtrer les champs selon la profondeur', () => {
      resetDom(`
        <div aqua_champ="Machines">
          <input aqua_champ="NomMachine" aqua_prof="1" value="PC-01" />
          <input aqua_champ="Ignore" aqua_prof="2" value="ne-doit-pas-apparaitre" />
        </div>
      `);

      const resultat = donnees_conf_analyse();

      expect(resultat).toEqual({
        Machines: [
          {
            NomMachine: 'PC-01',
          },
        ],
      });
    });

    it('convertit les champs datetime-local au format attendu', () => {
      resetDom(`
        <input type="datetime-local" aqua_champ="DateDebut" value="2026-05-09T14:30" />
      `);

      const resultat = donnees_conf_analyse();

      expect(resultat).toEqual({
        DateDebut: '2026-05-09T14:30:00Z',
      });
    });

    it('construit une liste par aqua_champ avec uniquement les objets non vides', () => {
      resetDom(`
        <div aqua_champ="Contacts">
          <input aqua_champ="Nom" aqua_prof="1" value="Alice" />
        </div>
        <div aqua_champ="Contacts">
          <input aqua_champ="Nom" aqua_prof="1" value="" />
        </div>
        <div aqua_champ="Contacts">
          <input aqua_champ="Nom" aqua_prof="1" value="Bob" />
        </div>
      `);

      const resultat = donnees_conf_analyse();

      expect(resultat).toEqual({
        Contacts: [
          { Nom: 'Alice' },
          { Nom: 'Bob' },
        ],
      });
    });
  });

  describe('section_remplie et verifier_remplissage', () => {
    it('section_remplie(true): retourne false si un champ required est vide', () => {
      resetDom(`
        <div id="section">
          <input required value="ok" />
          <input required value="" />
        </div>
      `);

      expect(section_remplie(document.getElementById('section'), true)).toBe(false);
    });

    it('section_remplie(true): retourne true quand tous les required sont remplis', () => {
      resetDom(`
        <div id="section">
          <input required value="ok" />
          <input required value="rempli" />
        </div>
      `);

      expect(section_remplie(document.getElementById('section'), true)).toBe(true);
    });

    it('section_remplie(false): retourne true si au moins un champ est rempli', () => {
      resetDom(`
        <div id="section">
          <input value="" />
          <input value="quelque chose" />
        </div>
      `);

      expect(section_remplie(document.getElementById('section'), false)).toBe(true);
    });

    it('section_remplie(false): retourne false si tous les champs sont vides', () => {
      resetDom(`
        <div id="section">
          <input value="" />
          <input value="" />
        </div>
      `);

      expect(section_remplie(document.getElementById('section'), false)).toBe(false);
    });

    it('verifier_remplissage active le bouton suivant si la section est complete', () => {
      resetDom(`
        <div id="section_identite">
          <input required value="Analyse 01" />
        </div>
        <button id="etape_suivante" disabled>Suivant</button>
      `);

      verifier_remplissage('section_identite', 'etape_suivante');

      expect(document.getElementById('etape_suivante').hasAttribute('disabled')).toBe(false);
    });

    it('verifier_remplissage desactive le bouton suivant si la section est incomplete', () => {
      resetDom(`
        <div id="section_identite">
          <input required value="" />
        </div>
        <button id="etape_suivante">Suivant</button>
      `);

      verifier_remplissage('section_identite', 'etape_suivante');

      expect(document.getElementById('etape_suivante').hasAttribute('disabled')).toBe(true);
    });
  });
});

describe('nouvelle_analyse_vis.js - priorite haute', () => {
  beforeEach(async () => {
    vi.restoreAllMocks();
    loadNouvelleAnalyseHtml();
    installVisMock();
    vi.stubGlobal('alert', vi.fn());
    vi.stubGlobal('confirm', vi.fn(() => true));
    setParentMock();
    loadNouvelleAnalyseVisScript();
    await flushPromises();
  });

  it('structure HTML: les ids critiques existent', () => {
    expect(document.getElementById('etape_3')).not.toBeNull();
    expect(document.getElementById('vue_generale')).not.toBeNull();
    expect(document.getElementById('popup-config-objet')).not.toBeNull();
    expect(document.getElementById('fond-popup')).not.toBeNull();
    expect(document.getElementById('nom_machine')).not.toBeNull();
    expect(document.getElementById('nature_equipement')).not.toBeNull();
    expect(document.getElementById('oui_ana')).not.toBeNull();
    expect(document.getElementById('non_ana')).not.toBeNull();
    expect(document.getElementById('liste_fichiers_analyses')).not.toBeNull();
    expect(document.getElementById('fichier_config_analyse')).not.toBeNull();
    expect(document.getElementById('adresses_equipement')).not.toBeNull();
  });

  it('ajout_liaison: refuse une auto-liaison', () => {
    const cb = vi.fn();

    ajout_liaison({ id: 'e1', from: '2', to: '2' }, cb);

    expect(alert).toHaveBeenCalledTimes(1);
    expect(cb).toHaveBeenCalledWith(null);
  });

  it('ajout_liaison: enregistre une liaison valide', () => {
    const cb = vi.fn();
    const liaison = { id: 'e2', from: '2', to: '3' };

    ajout_liaison(liaison, cb);

    expect(cb).toHaveBeenCalledWith(liaison);
  });

  it('verifier_archi_utilisable: active etape_3 si au moins un noeud a des fichiers', async () => {
    ajout_modif_noeud({ id: '2' }, vi.fn(), true);
    parent.window.go.main.App.ChoisirFichier.mockResolvedValueOnce(['C:\\tmp\\collecte.zip']);

    selection_fichier();
    await flushPromises();
    verifier_archi_utilisable();

    expect(document.getElementById('etape_3').hasAttribute('disabled')).toBe(false);
  });

  it('verifier_archi_utilisable: desactive etape_3 si aucun fichier', async () => {
    document.getElementById('etape_3').removeAttribute('disabled');
    ajout_modif_noeud({ id: '2' }, vi.fn(), true);
    parent.window.go.main.App.ChoisirFichier.mockResolvedValueOnce([]);

    selection_fichier();
    await flushPromises();
    verifier_archi_utilisable();

    expect(document.getElementById('etape_3').hasAttribute('disabled')).toBe(true);
  });

  it('validerDonneesNoeud: bloque si nom_machine est vide', () => {
    const cb = vi.fn();
    ajout_modif_noeud({ id: '2' }, cb, true);
    document.getElementById('nom_machine').value = '';
    document.getElementById('nature_equipement').value = 'ordinateur';

    validerDonneesNoeud();

    expect(alert).toHaveBeenCalled();
    expect(cb).not.toHaveBeenCalled();
  });

  it('validerDonneesNoeud: bloque si nature_equipement est vide', () => {
    const cb = vi.fn();
    ajout_modif_noeud({ id: '2' }, cb, true);
    document.getElementById('nom_machine').value = 'PC-01';
    document.getElementById('nature_equipement').value = '';

    recupere_liste_adresses = vi.fn(() => 'PC-01\n10.0.0.1');
    recuperer_informations_analyse = vi.fn();

    validerDonneesNoeud();

    expect(alert).toHaveBeenCalled();
    expect(cb).not.toHaveBeenCalled();
  });

  it('validerDonneesNoeud: succes met a jour le noeud et ferme le popup', () => {
    const cb = vi.fn();
    const noeud = { id: '2' };
    ajout_modif_noeud(noeud, cb, true);

    document.getElementById('nom_machine').value = 'PC-01';
    document.getElementById('nature_equipement').value = 'ordinateur';

    recupere_liste_adresses = vi.fn(() => 'PC-01\n10.0.0.1');
    recuperer_informations_analyse = vi.fn();
    verifier_archi_utilisable = vi.fn();

    validerDonneesNoeud();

    expect(noeud.label).toBe('PC-01\n10.0.0.1');
    expect(noeud.shape).toBe('image');
    expect(noeud.image).toBe('../assets/images/ordinateur.png');
    expect(document.getElementById('popup-config-objet').style.display).toBe('none');
    expect(document.getElementById('fond-popup').style.display).toBe('none');
    expect(verifier_archi_utilisable).toHaveBeenCalledTimes(1);
    expect(cb).toHaveBeenCalledWith(noeud);
  });

  it('get_donnees_reseau: ajoute les positions x/y sauf les edgeId*', () => {
    __networkInstance.body.nodes = {
      '1': { x: 1, y: 2 },
      '2': { x: 10, y: 20 },
      edgeId42: { x: 99, y: 99 },
    };

    const resultat = get_donnees_reseau();

    expect(resultat['1'].x).toBe(1);
    expect(resultat['1'].y).toBe(2);
    expect(resultat['2'].x).toBe(10);
    expect(resultat['2'].y).toBe(20);
    expect(resultat.edgeId42).toBeUndefined();
  });

  it('get_liens_reseau: retourne source/destination pour chaque lien', () => {
    __networkInstance.body.edges = {
      e1: { from: { id: '1' }, to: { id: '2' } },
      e2: { from: { id: '2' }, to: { id: '3' } },
    };

    const resultat = get_liens_reseau();

    expect(resultat).toEqual({
      e1: { source: '1', destination: '2' },
      e2: { source: '2', destination: '3' },
    });
  });
});
