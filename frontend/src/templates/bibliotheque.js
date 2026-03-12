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

if(parent.contrastes){
    contrastes_eleves();
}
if(parent.dyslexie){
    police_dyslexie();
}

function contrastes_eleves(){
    let couleurs = document.documentElement;
    couleurs.style.setProperty('--bleu-1', '#000');
    couleurs.style.setProperty('--bleu-2', '#000');
    couleurs.style.setProperty('--bleu-3', '#000');
    couleurs.style.setProperty('--bleu-4', '#FFF');
    couleurs.style.setProperty('--bleu-5', '#FFF');
    couleurs.style.setProperty('--bleu-6', '#FFF');
    couleurs.style.setProperty('--or-1', '#000');
    couleurs.style.setProperty('--or-2', '#000');
    couleurs.style.setProperty('--or-3', '#FFF');
    couleurs.style.setProperty('--or-4', '#FFF');
    couleurs.style.setProperty('--or-5', '#FFF');
    couleurs.style.setProperty('--or-6', '#FFF');
}

function police_dyslexie(){
    document.body.style.fontFamily = "Open Dyslexic";
}


// TODO : ajouter un paramètre "colonnes_a_afficher"
function creer_tableau_depuis_dico(dico, divOuMettreTableau, afficherFiltres, filtres, consignes_filtres, order_by, offset, lignes_selectionnees){
    // Créer un tableau Bootstrap
    let table = document.createElement('table');
    table.className = 'table_chronologie';
    table.style.fontSize = 'smaller'; // Réduire la taille du texte

    // Créer l'en-tête du tableau
    let thead = document.createElement('thead');
    let headerRow = document.createElement('tr');

    // Assumer que le premier objet dans résultat contient les colonnes
    let firstRow = Object.values(dico)[0];
    for (let key in firstRow) {
        let th = document.createElement('th');
        th.textContent = key;
        if (key == order_by){
            th.textContent ="⬇️"+ key + "⬇️";
        }
        th.onclick = () => trier_par(key);
        headerRow.appendChild(th);
    }
    thead.appendChild(headerRow);
    table.appendChild(thead);

    // Créer le corps du tableau
    let tbody = document.createElement('tbody');

    // Créer une partie "filtres"
    if (afficherFiltres){
        let ligneFiltres = document.createElement('tr');
        let ligneChoixFiltre = document.createElement('tr');
        for (let key in firstRow) {
            // Case dans laquelle on met la valeur que l'on veut filtrer
            let td = document.createElement('td');
            td.contentEditable = true;
            td.className = "input_massif";
            td.id = "valeur_filtre_" + key;
            if(filtres && filtres.has(key)){
                td.textContent = filtres.get(key);
            }
            td.onblur = () => appliquer_filtre(key);
            ligneFiltres.appendChild(td);
            // Sélection de comment on veut filtrer
            let tdSelect = document.createElement('td');
            let selectFiltre = document.createElement('select');
            selectFiltre.id = "consigne_filtre_" + key;
            selectFiltre.className = "filtre";
            selectFiltre.onchange = () => appliquer_filtre(key)
            let contient = document.createElement('option');
            contient.textContent = "🔤🔎🔤";
            selectFiltre.appendChild(contient);
            let commence_par = document.createElement("option");
            commence_par.textContent = "🔎🔤";
            selectFiltre.appendChild(commence_par);
            let finit_par = document.createElement("option");
            finit_par.textContent = "🔤🔎";
            selectFiltre.appendChild(finit_par);
            let exactement = document.createElement("option");
            exactement.textContent = "🔤 = 🔎";
            selectFiltre.appendChild(exactement);
            let superieur_a = document.createElement("option");
            superieur_a.textContent = "🔤 > 🔎";
            selectFiltre.appendChild(superieur_a);
            let inferieur_a = document.createElement("option");
            inferieur_a.textContent = "🔤 < 🔎";
            selectFiltre.appendChild(inferieur_a);
            if (consignes_filtres && consignes_filtres.has(key)){
                selectFiltre.value = consignes_filtres.get(key);
            }
            tdSelect.appendChild(selectFiltre);
            ligneChoixFiltre.appendChild(tdSelect);
        }
        tbody.appendChild(ligneChoixFiltre);
        tbody.appendChild(ligneFiltres);
    }

    // Remplir le reste du tableau 
    for (let [idRow, valueRows] of Object.entries(dico)) {
        let tr = document.createElement('tr');
        for (let [key, value] of Object.entries(valueRows)) {
                let td = document.createElement('td');
                let contenutd = document.createElement("div");
                contenutd.textContent = value;
                contenutd.classList.add("contenutd");
                td.appendChild(contenutd);
                tr.appendChild(td);
        }
        // Ajouter une case à cocher pour enregistrer la table
        if (lignes_selectionnees != undefined){
            let tdCasacocher =  document.createElement("td");
            tdCasacocher.classList.add("tdcasacocher");
            let casacocher = document.createElement("input");
            casacocher.type = "checkbox";
            let idEvenement = Number(idRow) + offset;
            casacocher.id = "casacocher_" + (idEvenement)
            if(lignes_selectionnees?.includes(idEvenement)){
                casacocher.checked = true;
            }
            casacocher.onchange = () => enregistrement_id(idEvenement);
            tdCasacocher.appendChild(casacocher);
            tr.appendChild(tdCasacocher);
        }
        // Ajouter la ligne au tableau
        tbody.appendChild(tr);
    }
    table.appendChild(tbody);

    divOuMettreTableau.appendChild(table);
}

let nb_documents = 35;

function afficher_salle_d_attente(emplacement){
    // Récupérer un document aléatoire
    let num_doc = Math.floor(Math.random() * nb_documents);
    let iframe = document.createElement("iframe");
    iframe.src = "../assets/documents/" +  num_doc + ".pdf";
    emplacement.appendChild(iframe);
    emplacement.style.display = "inline";
}


function fermer_popup(idPopup){
    document.getElementById(idPopup).style.display = "none";
    document.getElementById("fond_popup").style.display = "none"
}