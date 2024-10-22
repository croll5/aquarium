function choix_enregistrement(){
    parent.window.go.main.App.CreationDossierNouveauModele().then(resultat =>{
        document.getElementById("enregistrement").value = resultat;
        document.getElementById("archives").value = "";
        document.getElementById("valider").style.display = "none";
    })
}

function choix_orc(){
    let chemin_enreg = document.getElementById("enregistrement").value;
    if(chemin_enreg == ""){
        alert("Vous deviez d'abord choisir où vous voulez enregistrer votre modèle")
    }
    else{
        document.getElementById("patientez").style.display = "inline";
        document.getElementById("formulaire").style.display = "none";
        parent.window.go.main.App.AjoutORCNouveauProjet().then(resultat =>{
            document.getElementById("archives").value = resultat;
            if(resultat == ""){
                document.getElementById("valider").style.display = "none";
            }
            else{
                let nom_modele = document.getElementById("nom_modele").value;
                if(nom_modele != ""){
                    document.getElementById("valider").style.display = "inline";
                }
            }
            document.getElementById("patientez").style.display = "none";
            document.getElementById("formulaire").style.display = "inline";
        })
    }
}

function change_nom_modele(){
    let nom_modele = document.getElementById("nom_modele").value;
    let enregistrement = document.getElementById("enregistrement").value;
    let archives = document.getElementById("archives").value;
    if(nom_modele == "" || enregistrement == "" || archives == ""){
        document.getElementById("valider").style.display = "none";
    }
    else{
        document.getElementById("valider").style.display = "inline";
    }
}

function validation(){
    let nom_modele = document.getElementById("nom_modele").value;
    let description = document.getElementById("description").value;
    let supprimerORC = document.getElementById("avec_nettoyage").checked;
    document.getElementById("patientez_analyse").style.display = "inline";
    parent.window.go.main.App.ValidationCreationModele(nom_modele, description, supprimerORC).then(resultat =>{
        window.location.replace("../index.html"); //"../accueil/index.html"
    })
}