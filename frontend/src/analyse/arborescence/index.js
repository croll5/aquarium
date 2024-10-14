construireArborescence("arborescence", []);

function construireArborescence(id_racine, chemin_num){
    let racine = document.getElementById(id_racine);
    if(racine == undefined || racine.children.length > 1){
        return
    }
    try {
        document.body.style.cursor = "wait"; 
    } catch (error) {
    }
    parent.window.go.main.App.ArborescenceMachineAnalysee(chemin_num).then(resultat =>{ 
        if(resultat.length == 0){
            document.getElementById("extraction_arborescence").style.display = "inline";
            document.getElementById("patientez").style.display = "none";
            return
        }
        // On cherche s'il faut affcher les indicateurs de légitimité
        let afficher_inconnu, afficher_ok, afficher_suspect;
        try{
        afficher_inconnu = document.getElementById("affiche_inconnu").checked;
        afficher_ok = document.getElementById("affiche_ok").checked;
        afficher_suspect = document.getElementById("affiche_suspect").checked;
        }catch(error){
            alert(document.getElementById("affiche_inconnu"));
            afficher_inconnu = true;
            afficher_ok = true;
            afficher_suspect = true;
        }
        for(let i=0; i < resultat.length; i++){
            if(resultat[i]["ADesEnfants"]){
                let enfant = document.createElement("details");
                enfant.id = String.prototype.concat(id_racine, "_", i);
                enfant.className = "dossier_arborescence";
                let chemin_enfant = chemin_num.concat([i]);
                let titre_enfant = document.createElement("summary");
                titre_enfant.textContent = resultat[i]["Nom"];
                titre_enfant.onclick = function(ev){return construireArborescence(enfant.id, chemin_enfant)};
                enfant.appendChild(titre_enfant);
                racine.appendChild(enfant);
            }else{
                let enfant = document.createElement("p");
                let legitimite = document.createElement("strong");
                legitimite.textContent = "😇";
                legitimite.className = "legitimite_ok";
                legitimite.style.display = afficher_ok ? "inline" : "none";
                if(resultat[i]["EnfantsSuspects"] > 0){
                    legitimite.textContent = "🥴";
                    legitimite.className = "legitimite_suspect";
                    legitimite.style.display = afficher_suspect ? "inline" : "none";
                }else if(resultat[i]["EnfantsInconnus"] > 0){
                    legitimite.textContent = "😵";
                    legitimite.className = "legitimite_aucune";
                    legitimite.style.display = afficher_inconnu ? "inline" : "none";
                }
                enfant.id = String.prototype.concat(id_racine, "_", i);
                enfant.className = "fichier_arborescence";
                enfant.textContent = resultat[i]["Nom"];
                enfant.appendChild(legitimite);
                racine.appendChild(enfant);
            }
        }
        document.body.style.cursor = "default"; 
        document.getElementById("affichage_arbo").style.display = "inline";
    })
    
}

function extraire_arborescence(){
    document.getElementById("extraction_arborescence").style.display = "none";
    document.getElementById("patientez").style.display = "inline";
    let avec_modele = document.getElementById("avec_modele").checked;
    parent.window.go.main.App.ExtraireArborescence(avec_modele).then(resultat =>{
        document.getElementById("patientez").style.display = "none";
        construireArborescence("arborescence", []);
    })
}

function affichage_legitimite(id_checkbox, nom_classe){
    let choix = document.getElementById(id_checkbox);
    let smileys = document.getElementsByClassName(nom_classe);
    if (choix.checked){
        for(const element of smileys){
            element.style.display = "inline";
        }
    } else{
        for(const element of smileys){
            element.style.display = "none";
        }
    }
}