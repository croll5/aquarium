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
function creer_tableau_depuis_dico(dico, divOuMettreTableau, afficherFiltres, filtres, consignes_filtres, order_by){
    // Créer un tableau Bootstrap
    let table = document.createElement('table');
    table.className = 'table table-striped table-bordered';
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
            if(filtres.has(key)){
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
            if (consignes_filtres.has(key)){
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
            td.textContent = value;
            tr.appendChild(td);
        }
        tbody.appendChild(tr);
    }
    table.appendChild(tbody);

    divOuMettreTableau.appendChild(table);
}