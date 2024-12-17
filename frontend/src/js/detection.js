afficher_regles(false);

function afficher_regles(lancer){
    if(lancer){
        document.getElementById("regles").innerHTML = "";
    }
    parent.window.go.main.App.ListeReglesDetection(lancer).then(resultat =>{
        let conteneur = document.getElementById("regles");
        let regles = Object.keys(resultat);
        regles.forEach(regle => {
            let detail_regle =  document.createElement("details");
            let contenu_regle = document.createElement("div")
            contenu_regle.id = regle.replace(" ", "");
            contenu_regle.classList.add("contenuRegle");
            contenu_regle.setAttribute("resultat", resultat[regle]);
            detail_regle.appendChild(contenu_regle)
            detail_regle.onclick = function(ev){return informations_regle(contenu_regle.id, regle)};
            let titre_regle = document.createElement("summary");
            titre_regle.textContent = regle;
            let smiley = document.createElement("strong");
            smiley.classList.add("etatRegle");
            switch (resultat[regle]) {
                case 0:
                    smiley.textContent = " 😴"
                    break;
                case 1:
                    smiley.textContent = " 🤓"
                    break;
                case 2:
                    smiley.textContent += " 🥸"
                    break;
            }
            titre_regle.appendChild(smiley);
            detail_regle.appendChild(titre_regle);
            conteneur.appendChild(detail_regle);
        });
    })
}

function informations_regle(id, nom_regle){
    if(document.getElementById(id).childElementCount <= 1){
        parent.window.go.main.App.InfosRegleDetection(nom_regle).then(resultat =>{
            let regle = document.getElementById(id);
            // Criticité 
            let rangeCriticite = document.createElement("input");
            rangeCriticite.type = "range";
            rangeCriticite.min = 0;
            rangeCriticite.max = 5;
            rangeCriticite.value = resultat["criticite"];
            rangeCriticite.readOnly = true;
            switch (resultat["criticite"]) {
                case 0:
                    rangeCriticite.style.accentColor = "#18C700";
                    break;
                case 1:
                    rangeCriticite.style.accentColor = "#72C702";
                    break;
                case 2:
                    rangeCriticite.style.accentColor = "#C2D16C";
                    break;
                case 3:
                    rangeCriticite.style.accentColor = "#E3A500";
                    break;
                case 4:
                    rangeCriticite.style.accentColor = "#F06136";
                    break;
                case 5:
                    rangeCriticite.style.accentColor = "#D42222";
                    break;
                default:
                    break;
            }
            let criticite = document.createElement("div");
            criticite.innerHTML = "<strong>Criticité : </strong>" + resultat["criticite"] + " ";
            criticite.appendChild(rangeCriticite);
            regle.appendChild(criticite);
            // Requete SQL
            let sql = document.createElement("p");
            sql.innerText = resultat["sql"];
            sql.className = "code sql";
            regle.appendChild(sql);
            // Description de la regle
            let description = document.createElement("p");
            description.innerHTML = "<strong>Description : </strong>" + resultat["description"];
            regle.appendChild(description);
            // Auteur de la regle
            let auteur = document.createElement("p");
            auteur.innerHTML = "<strong>Auteur : </strong>" + resultat["auteur"];
            regle.appendChild(auteur);
            
            // Boutons à afficher
            // Récupération de la valeur de résultat de la règle
            let resultatRegle = regle.getAttribute("resultat");
            if(resultatRegle == 0){
                let lancerRegle = document.createElement("button");
                lancerRegle.innerText = "Lancer cette règle";
                lancerRegle.onclick = function(ev){return lancer_regle(id, nom_regle)};
                lancerRegle.classList.add("bouton_sombre");
                regle.appendChild(lancerRegle);
            }else if(resultatRegle == 2){
                let afficherResulataRegle = document.createElement("button");
                afficherResulataRegle.innerText = "Afficher le résultat";
                afficherResulataRegle.classList.add("bouton_sombre");
                regle.appendChild(afficherResulataRegle);
            }
        })
    }
}

function lancer_regle(id, nom_regle){
    parent.window.go.main.App.ResultatRegleDetection(nom_regle).then(resultat =>{
        let regle = document.getElementById(id);
        let bouton = regle.querySelector("button");
        if(resultat == 1){
            regle.parentNode.querySelector(".etatRegle").textContent = " 🤓";
            regle.removeChild(bouton);
        }else if(resultat == 2){
            regle.parentNode.querySelector(".etatRegle").textContent = " 🥸";
            bouton.innerText = "Afficher le résultat";
            bouton.onclick = "";
        }
    })
}