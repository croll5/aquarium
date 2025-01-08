let position_dans_table = 0
let requete = "SELECT extracteur, horodatage, message, source FROM chronologie";
let taille_requete = 0;

affichage_table(true);

function affichage_table(majTaille){
    let emplacement_resultat = document.getElementById("emplacement_table");
    parent.window.go.main.App.ResultatRequeteSQLExtraction(requete, position_dans_table, 5).then(resultat =>{
        
    emplacement_resultat.innerHTML = "";
        creer_tableau_depuis_dico(resultat, emplacement_resultat);
        if (majTaille){
            parent.window.go.main.App.TailleRequeteSQLExtraction(requete).then(nbLignes =>{
                console.log(nbLignes);
                taille_requete = nbLignes;
            });
        }
    })
    
}

function tourner_page(extremes, difference){
    if (extremes == -1){
        position_dans_table = 0;
    }
    if(extremes == 1){
        position_dans_table = taille_requete - 5;
    }
    position_dans_table = Math.max(0, position_dans_table + difference);
    affichage_table(false);
    document.getElementById("indicateur_page").textContent = position_dans_table + "-" + (position_dans_table+5);
}

let estEnTrainDeScroller = false
let divPOurScroller = document.getElementById("endroit_pour_scroller")
divPOurScroller.addEventListener("scroll", (event) => {
    if (!estEnTrainDeScroller){
        estEnTrainDeScroller = true
        if(divPOurScroller.scrollTop > 10000){
            tourner_page(0, 1);
        }else if(divPOurScroller.scrollTop < 10000){
            tourner_page(0,-1);
        }
        divPOurScroller.scrollTop = 10000;
        setTimeout(() => estEnTrainDeScroller = false, 500)
    }
});
