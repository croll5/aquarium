parent.window.go.main.App.ArborescenceMachineAnalysee().then(resultat =>{
    if(resultat["nom"] == "" && resultat["enfants"] == undefined){
        document.getElementById("extraction_arborescence").style.display = "inline";
    }else{
        construireArborescence(resultat, "arborescence");
    }
})

function construireArborescence(dossier, id_racine){
    let nom_dossier = document.createElement("p");
    nom_dossier.textContent = dossier["nom"];
    let sous_dossiers = document.createElement("div");
    sous_dossiers.id = "enfants_" + dossier["nom"];
    sous_dossiers.className = "dossier_arborescence";
    let contenant = document.getElementById(id_racine);
    contenant.appendChild(nom_dossier);
    contenant.appendChild(sous_dossiers);
    if(dossier["enfants"] == undefined){
        nom_dossier.className = "fichier_arborescence";
        return;
    }
    for (const enfant of dossier["enfants"]) {
        construireArborescence(enfant, "enfants_" + dossier["nom"]);
    }
}

function extraire_arborescence(){
    parent.window.go.main.App.ExtraireArborescence().then(resultat =>{
        alert("Chalut la compagnie !");
    })
}