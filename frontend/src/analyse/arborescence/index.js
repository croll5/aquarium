parent.window.go.main.App.ArborescenceMachineAnalysee().then(resultat =>{
    if(resultat["nom"] == "" && resultat["enfants"] == undefined){
        document.getElementById("extraction_arborescence").style.display = "inline";
        document.getElementById("patientez").style.display = "none";
    }else{
        construireArborescence(resultat, "arborescence")
        document.getElementById("patientez").style.display = "none";

    }
})

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

function extraire_arborescence(){
    document.getElementById("extraction_arborescence").style.display = "none";
    document.getElementById("patientez_extraction").style.display = "inline";
    parent.window.go.main.App.ExtraireArborescence(true).then(resultat =>{
        document.getElementById("patientez").style.display = "inline";
        document.getElementById("patientez_extraction").style.display = "none";
        construireArborescence(resultat, "arborescence");
        document.getElementById("patientez").style.display = "none";
    })
}