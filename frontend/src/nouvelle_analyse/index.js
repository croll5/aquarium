function choix_enregistrement(){
    parent.window.go.main.App.CreationNouveauProjet().then(resultat =>{
        document.getElementById("enregistrement").value = resultat;
        document.getElementById("archives").value = "";
        document.getElementById("valider").style.display = "none";
    })
}

function choix_orc(){
    let chemin_enreg = document.getElementById("enregistrement").value;
    if(chemin_enreg == ""){
        alert("Vous deviez d'abord choisir où vous voulez enregistrer votre ORC")
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
                let nom_auteur = document.getElementById("auteur").value;
                if(nom_auteur != ""){
                    document.getElementById("valider").style.display = "inline";
                }
            }
            document.getElementById("patientez").style.display = "none";
            document.getElementById("formulaire").style.display = "inline";
        })
    }
}

function change_auteur(){
    let nom_auteur = document.getElementById("auteur").value;
    let enregistrement = document.getElementById("enregistrement").value;
    let archives = document.getElementById("archives").value;
    if(nom_auteur == "" || enregistrement == "" || archives == ""){
        document.getElementById("valider").style.display = "none";
    }
    else{
        document.getElementById("valider").style.display = "inline";
    }
}

function validation(){
    let auteur = document.getElementById("auteur").value;
    let description = document.getElementById("description").value;
    parent.window.go.main.App.ValidationCreationProjet(auteur, description).then(resultat =>{
        if(resultat){ 
            window.location.replace("../analyse/extraction/index.html");
            parent.document.getElementsByTagName("header")[0].style.display = "inline";
            let onglet_courant = parent.document.getElementById("onglet_extraction");
            onglet_courant.style.backgroundColor = "#FCF5DC";
            onglet_courant.style.color = "#000";
        }else{
            window.location.replace("../index.html")
        }
    })
}