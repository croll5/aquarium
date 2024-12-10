function change_onglet(destination, id_onglet){
    document.getElementsByTagName("iframe")[0].src = destination;
    // On met tous les autres onglets à la couleur standard
    let onglets = document.getElementsByClassName("onglet");
    for (const onglet of onglets) {
        onglet.style.backgroundColor = "#856B0D";
        onglet.style.color = "#fff";
    }
    // On met l'onglet sur lequel on va à la couleur de la page
    let onglet_courant = parent.document.getElementById(id_onglet);
    onglet_courant.style.backgroundColor = "#FCF5DC";
    onglet_courant.style.color = "#000";
}

function accueil(){
    document.getElementsByTagName("iframe")[0].src = "html/accueil.html";
    document.getElementsByTagName("header")[0].style.display = "none";
}

var contrastes = false;
var dyslexie = false;
var non_aux_bubulles = false;