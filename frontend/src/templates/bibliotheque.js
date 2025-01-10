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
function creer_tableau_depuis_dico(dico, divOuMettreTableau, afficherFiltres, filtres, order_by){
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
        for (let key in firstRow) {
            let td = document.createElement('td');
            td.contentEditable = true;
            td.className = "input_massif";
            td.placeholder = "filtrer...";
            td.id = "filtre_" + key;
            if(filtres.has(key)){
                td.textContent = filtres.get(key);
            }
            td.onblur = () => appliquer_filtre(key);
            ligneFiltres.appendChild(td);
        }
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