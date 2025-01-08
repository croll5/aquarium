/**
 * Execution au chargement de la page
 */
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

    afficher_regles(false);
};


/**
 * Recuperation de l'ensemble des regles de detection au format json au chargement de la page
 * @param {boolean} lancer - Si True alors execute la requete SQL de chaque regle
 */
function afficher_regles(lancer) {
    const conteneur = document.getElementById("regles");
    if (lancer) {
        conteneur.innerHTML = "";
    }
    let listRegles;
    parent.window.go.main.App.ListeReglesDetection(lancer).then(resultat => {
        listRegles = new Map(Object.entries(resultat));
        listRegles.forEach((values, regle) => {
            // value for html balises
            const contenu_id = regle//.replace(" ", "");
            const contenu_state = values['state'];
            const detail_smiley = (() => {
                switch (contenu_state) {
                    case 0: return " 😴";
                    case 1: return " 🤓";
                    case 2: return " 🥸";
                    default: return "";
                }
            })();
            // html block
            const detail_regle = document.createElement('details');
            detail_regle.innerHTML = `
                    <div id="${contenu_id}" class="contenuDetails" resultat="${contenu_state}"></div>
                    <summary>${regle}<strong class="etatRegle">${detail_smiley}</strong></summary>
                `;
            detail_regle.onclick = () => informations_regle(contenu_id, regle);
            // Add element in the page
            conteneur.appendChild(detail_regle);
        });
        update_summary()
    });
}

function update_summary() {
    parent.window.go.main.App.ListeReglesDetection(false).then(resul => {
        let nbRules = new Map(Object.entries(resul)).size;
        parent.window.go.main.App.StatutReglesDetection().then(resultat => {
            const errorCount = resultat.filter(item => item.isError === 1).length;
            const nbElement = resultat.length;
            document.getElementById("total").innerHTML = "Total:<br>"+nbRules;
            document.getElementById("notExecuted").innerHTML = "Inconnu:<br>"+(nbRules-nbElement);
            document.getElementById("valided").innerHTML = "Validation:<br>"+(nbElement-errorCount)
            document.getElementById("detected").innerHTML = "Detection:<br>"+errorCount
        });
    });
}

/**
 * Affiche les informations d'une regle au clic
 */
function informations_regle(id, nom_regle) {
    if (document.getElementById(id).childElementCount <= 1) {
        parent.window.go.main.App.InfosRegleDetection(nom_regle).then(resultat => {
            // rule html balise and parameters
            const regle = document.getElementById(id);
            const criticiteColor = ["#18C700", "#72C702", "#C2D16C", "#E3A500", "#F06136", "#D42222"][resultat["criticite"]] || "#000000";
            // html block to add
            const detail_regle_open = `
                <div>
                    <strong>Criticité : </strong>${resultat["criticite"]}
                    <input type="range" min="0" max="5" value="${resultat["criticite"]}" style="accent-color:${criticiteColor}" oninput="this.value=${resultat['criticite']}">
                </div>
                <p class="code sql"><strong>SQL : </strong>${resultat["sql"]}</p>
                <p><strong>Description : </strong>${resultat["description"]}</p>
                <p><strong>Auteur : </strong>${resultat["auteur"]}</p>
            `;
            // Add element in the page
            regle.innerHTML += detail_regle_open;
            // Button to print the dataframe result
            if (regle.getAttribute("resultat") == 0) {
                const bouton = document.createElement("button");
                bouton.className = "bouton_sombre";
                bouton.innerText = "Lancer cette règle";
                bouton.onclick = () => lancer_regle(id, nom_regle);
                regle.appendChild(bouton);
            } else if (regle.getAttribute("resultat") == 2) {
                const bouton = document.createElement("button");
                bouton.className = "bouton_sombre";
                bouton.innerText = "Afficher le résultat";
                bouton.onclick = () => afficher_resultat_regle(id, nom_regle);
                regle.appendChild(bouton);
            }
            // Button for local rules
            if (!resultat["IsGlobal"]) {
                // delete button
                const boutonModif = document.createElement("button");
                boutonModif.onclick = () => modifier_regle_panel(id, nom_regle);
                boutonModif.className = "bouton_sombre";
                boutonModif.innerText = "Modifier";
                boutonModif.style = "background-color:orange; margin-left:1rem;"
                regle.appendChild(boutonModif);
                const boutonDel = document.createElement("button");
                boutonDel.onclick = () => supprimer_regle(id, nom_regle);
                boutonDel.className = "bouton_sombre";
                boutonDel.innerText = "Supprimer";
                boutonDel.style = "background-color:red; margin-left:1rem;"
                regle.appendChild(boutonDel);

            }
            update_summary()
        });
    }
}


/**
 * Lance une sequence d'execution de la requete SQL d'une regle au clic
 */
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
        update_summary()
    });
}


/**
 * Ferme tous les volets pliant
 */
function closeAllRules() {
    const detailsElements = document.querySelectorAll("#regles details");
    detailsElements.forEach(detail => {
        detail.removeAttribute("open");
    });
}

/**
 * Verification de la requete SQL du formulaire de creation d'une regle
 */
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

/**
 * Creation d'une nouvelle regle
 */
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

    document.getElementById("popup-newRule").style.display = "none";
    document.getElementById("regles").innerHTML = "";
    afficher_regles(false);
    update_summary()
}

function supprimer_regle(id, nom_regle) {
    parent.window.go.main.App.Delete_rule(nom_regle).then(_ => {
        document.getElementById("regles").innerHTML = "";
        afficher_regles(false);
        update_summary()
    });
}

function modifier_regle_panel(id, nom_regle) {
    parent.window.go.main.App.InfosRegleDetection(nom_regle).then(resultat => {
        document.getElementById("nom").value = nom_regle;
        document.getElementById("auteur").value = resultat["auteur"];
        document.getElementById("description").value = resultat["description"];
        document.getElementById("criticite").value = resultat["criticite"];
        document.getElementById("sql").value = resultat["sql"];

        document.getElementById("popup-newRule").style.display = "block";
        update_summary()
    });
}



function afficher_resultat_regle(id, nom_regle) {
     parent.window.go.main.App.ResultatsSQL(nom_regle).then(resultat => {
        if (!resultat) {return;}
        // add results in the popup
        let div_db_infos = document.querySelector("#popup-resultRule #table_values");
        div_db_infos.textContent = ''

        // Créer un conteneur avec une barre de défilement horizontal
        let scrollContainer = document.createElement('div');
        scrollContainer.style.maxHeight = '450px'
        scrollContainer.style.overflow = 'auto';

        creer_tableau_depuis_dico(resultat, scrollContainer)
        div_db_infos.appendChild(scrollContainer);

        document.querySelector("#popup-resultRule .modal-content").style.width = "90%";
        document.getElementById("popup-resultRule").style.display = "block";
        update_summary()
    });
}


