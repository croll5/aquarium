/**parent.window.go.main.App.ArborescenceMachineAnalysee().then(resultat =>{
    if(resultat["nom"] == "" && resultat["enfants"] == undefined){
        document.getElementById("extraction_arborescence").style.display = "inline";
        document.getElementById("patientez").style.display = "none";
    }else{
        construireArborescence(resultat, "arborescence")
        document.getElementById("patientez").style.display = "none";

    }
})*/

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
                let legitimite = "😇"
                if(resultat[i]["EnfantsSuspects"] > 0){
                    legitimite = "🥴"
                }else if(resultat[i]["EnfantsInconnus"] > 0){
                    legitimite = "😵"
                }
                enfant.id = String.prototype.concat(id_racine, "_", i);
                enfant.className = "fichier_arborescence";
                enfant.textContent = String.prototype.concat(resultat[i]["Nom"], " ", legitimite);
                racine.appendChild(enfant);
            }
        }
        document.body.style.cursor = "default"; 
    })
    
}
/**
function construireArborescence(dossier, id_racine, numero_dossier){
    if(dossier["enfants"] != undefined){
        let contenant = document.createElement("details");
        let id_contenant = String.prototype.concat(id_racine, "_", numero_dossier);
        contenant.id = id_contenant;
        contenant.className = "dossier_arborescence";
        let nom_dossier = document.createElement("summary");
        nom_dossier.textContent = dossier["nom"];
        contenant.appendChild(nom_dossier);
        document.getElementById(id_racine).appendChild(contenant);
        for(let i = 0; i < dossier["enfants"].length; i++){
            construireArborescence(dossier["enfants"][i], id_contenant, i);
        }
    }
    else{
        let nom_fichier = document.createElement("p");
        let legitimite = "😵"
        if(dossier["legitimite"] == 1){
            legitimite = "🥴"
        }else if(dossier["legitimite"] == 2){
            legitimite = "😇"
        }
        nom_fichier.textContent = dossier["nom"] + legitimite;
        nom_fichier.className = "fichier_arborescence";
        document.getElementById(id_racine).appendChild(nom_fichier);
    }
    return true;
}
     */

function extraire_arborescence(){
    document.getElementById("extraction_arborescence").style.display = "none";
    document.getElementById("patientez").style.display = "inline";
    parent.window.go.main.App.ExtraireArborescence(true).then(resultat =>{
        document.getElementById("patientez").style.display = "none";
        construireArborescence(resultat, "arborescence");
    })
}