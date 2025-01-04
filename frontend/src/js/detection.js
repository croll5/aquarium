afficher_regles(false);

function afficher_regles(lancer) {
    if (lancer) {
        document.getElementById("regles").innerHTML = "";
    }
    parent.window.go.main.App.ListeReglesDetection(lancer).then(resultat => {
        const conteneur = document.getElementById("regles");
        const regles = Object.keys(resultat);
        regles.forEach(regle => {
            const detail_regle = document.createElement("details");
            const contenu_regle = document.createElement("div");
            contenu_regle.id = regle.replace(" ", "");
            contenu_regle.classList.add("contenuRegle");
            contenu_regle.setAttribute("resultat", resultat[regle]);

            detail_regle.appendChild(contenu_regle);
            detail_regle.onclick = () => informations_regle(contenu_regle.id, regle);

            const titre_regle = document.createElement("summary");
            titre_regle.textContent = regle;

            const smiley = document.createElement("strong");
            smiley.classList.add("etatRegle");
            smiley.textContent = getSmiley(resultat[regle]);

            titre_regle.appendChild(smiley);
            detail_regle.appendChild(titre_regle);
            conteneur.appendChild(detail_regle);
        });
    });
}

function getSmiley(etat) {
    switch (etat) {
        case 0: return " 😴";
        case 1: return " 🤓";
        case 2: return " 🥸";
        default:return "";
    }
}

function informations_regle(id, nom_regle) {
    if (document.getElementById(id).childElementCount <= 1) {
        parent.window.go.main.App.InfosRegleDetection(nom_regle).then(resultat => {
            const regle = document.getElementById(id);
            regle.appendChild(createCriticiteElement(resultat["criticite"]));
            regle.appendChild(createParagraph("SQL", "code sql", resultat["sql"]));
            regle.appendChild(createParagraph("Description", "", resultat["description"]));
            regle.appendChild(createParagraph("Auteur", "", resultat["auteur"]));
            appendActionButtons(regle, regle.getAttribute("resultat"), id, nom_regle);
        });
    }
}

function createParagraph(label, className, content) {
    const p = document.createElement("p");
    p.innerHTML = `<strong>${label} : </strong>${content}`;
    p.className = className;
    return p;
}

function createCriticiteElement(criticite) {
    const rangeCriticite = document.createElement("input");
    rangeCriticite.type = "range";
    rangeCriticite.min = 0;
    rangeCriticite.max = 5;
    rangeCriticite.value = criticite;
    rangeCriticite.readOnly = true;
    rangeCriticite.style.accentColor = getCriticiteColor(criticite);

    const criticiteDiv = document.createElement("div");
    criticiteDiv.innerHTML = `<strong>Criticité : </strong>${criticite} `;
    criticiteDiv.appendChild(rangeCriticite);
    return criticiteDiv;
}

function getCriticiteColor(criticite) {
    const colors = ["#18C700", "#72C702", "#C2D16C", "#E3A500", "#F06136", "#D42222"];
    return colors[criticite] || "#000000";
}

function appendActionButtons(regle, resultatRegle, id, nom_regle) {
    if (resultatRegle == 0) {
        const lancerRegle = document.createElement("button");
        lancerRegle.innerText = "Lancer cette règle";
        lancerRegle.classList.add("bouton_sombre");
        lancerRegle.onclick = () => lancer_regle(id, nom_regle);
        regle.appendChild(lancerRegle);
    } else if (resultatRegle == 2) {
        const afficherResulataRegle = document.createElement("button");
        afficherResulataRegle.innerText = "Afficher le résultat";
        afficherResulataRegle.classList.add("bouton_sombre");
        regle.appendChild(afficherResulataRegle);
    }
}


function lancer_regle(id, nom_regle) {
    parent.window.go.main.App.ResultatRegleDetection(nom_regle).then(resultat => {
        const regle = document.getElementById(id);
        const bouton = regle.querySelector("button");
        const etatRegle = regle.parentNode.querySelector(".etatRegle");

        if (resultat == 1) {
            etatRegle.textContent = " 🤓";
            bouton.remove();
        } else if (resultat == 2) {
            etatRegle.textContent = " 🥸";
            bouton.innerText = "Afficher le résultat";
            bouton.onclick = () => afficher_resultat_regle(id, nom_regle);
        }
    });
}



function validateSQL() {
    const input = document.getElementById("sql");
    const sqlPattern = /SELECT\s.*\sFROM\s.*\sWHERE\s.*/ig;
    input.value = input.value.trim();

    if (sqlPattern.test(input.value)) {
        input.classList.remove('invalid');
        return true;
    } else {
        input.classList.add('invalid');
        alert("La requête SQL doit être au format 'SELECT % FROM % WHERE %': " + input.value);
        return false;
    }
}


window.onload = function() {
    document.querySelectorAll(".button").forEach(btn => {
        btn.onclick = () => {
            const modal = btn.getAttribute("data-modal");
            document.getElementById(modal).style.display = "block";
        };
    });

    window.onclick = event => {
        if (event.target.className === "modal") {
            event.target.style.display = "none";
        }
    };

    const value = document.querySelector("#criticite_value");
    const input = document.querySelector("#criticite");
    value.textContent = input.value;
    input.addEventListener("input", event => {
        value.textContent = event.target.value;
    });

    // Ajoutez l'événement au bouton pour fermer toutes les règles
    const closeButton = document.querySelector('button[onclick="closeAllRules()"]');
    closeButton.addEventListener("click", closeAllRules);
};

function creation_regle() {
    const regle = {
        "nom": document.getElementById("nom").value,
        "auteur": document.getElementById("auteur").value,
        "description": document.getElementById("description").value,
        "criticite": parseInt(document.getElementById("criticite").value),
        "sql": document.getElementById("sql").value
    };

    const jsonString = JSON.stringify(regle, null, 2);
    parent.window.go.main.App.CreationReglesDetection(jsonString);
}


function afficher_resultat_regle(id, nom_regle) {
    res = parent.window.go.main.App.ResultatsSQL(nom_regle).then(resultat => {
        if (resultat == null) {return;}
        // add results in the popup
        let div_db_infos = document.querySelector("#popup-resultRule #table_values");


        console.table(resultat);
        div_db_infos.textContent = ''
        // Créer un conteneur avec une barre de défilement horizontal
        let scrollContainer = document.createElement('div');
        scrollContainer.style.maxHeight = '450px'
        scrollContainer.style.overflowX = 'auto';
        scrollContainer.style.overflowY = 'auto';


        // Créer un tableau Bootstrap
        let table = document.createElement('table');
        table.className = 'table table-striped table-bordered';
        table.style.fontSize = 'smaller'; // Réduire la taille du texte

        // Créer l'en-tête du tableau
        let thead = document.createElement('thead');
        let headerRow = document.createElement('tr');

        // Assumer que le premier objet dans résultat contient les colonnes
        let firstRow = Object.values(resultat)[0];
        for (let key in firstRow) {
            let th = document.createElement('th');
            th.textContent = key;
            headerRow.appendChild(th);
        }
        thead.appendChild(headerRow);
        table.appendChild(thead);

        // Créer le corps du tableau
        let tbody = document.createElement('tbody');
        for (let [idRow, valueRows] of Object.entries(resultat)) {
            let tr = document.createElement('tr');
            for (let [key, value] of Object.entries(valueRows)) {
                let td = document.createElement('td');
                td.textContent = value;
                tr.appendChild(td);
            }
            tbody.appendChild(tr);
        }
        table.appendChild(tbody);

        scrollContainer.appendChild(table);
        div_db_infos.appendChild(scrollContainer);

        document.querySelector("#popup-resultRule .modal-content").style.width = "90%";
        document.getElementById("popup-resultRule").style.display = "block";
    });
}

function closeAllRules() {
    const detailsElements = document.querySelectorAll("#regles details");
    detailsElements.forEach(detail => {
        detail.removeAttribute("open");
    });
}
