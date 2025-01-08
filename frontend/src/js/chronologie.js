let position_dans_table = 0
let requete = "SELECT extracteur, horodatage, message, source FROM chronologie";
let taille_requete = 0;
let filtres = new Map();

affichage_table(true);

/* À l'affichage de la zone de recherche click-bouton, on affiche la liste des tables */
let divChangementRequete = document.getElementById("changement_requete");
divChangementRequete.addEventListener("toggle", (event) => {
    let click_bouton = document.getElementById("zone_recherche_click_bouton").style.display != "none";
    if (divChangementRequete.open && click_bouton) {
        let selecteurTable = document.getElementById("choix_table");
        selecteurTable.innerHTML = "";
        parent.window.go.main.App.GetListeTablesExtraction().then(resultat => {
            for (const i in resultat) {
                let nom_table = document.createElement("option");
                nom_table.value = resultat[i];
                nom_table.textContent = resultat[i];
                selecteurTable.appendChild(nom_table)
            }
        });
    }
  });
  

function affichage_table(majTaille){
    let emplacement_resultat = document.getElementById("emplacement_table");
    parent.window.go.main.App.ResultatRequeteSQLExtraction(requete, position_dans_table, 5).then(resultat =>{
    document.getElementById("indicateur_page").textContent = position_dans_table + "-" + (position_dans_table+5);
    emplacement_resultat.innerHTML = "";
        creer_tableau_depuis_dico(resultat, emplacement_resultat, true, filtres);
        if (majTaille){
            parent.window.go.main.App.TailleRequeteSQLExtraction(requete).then(nbLignes =>{
                console.log(nbLignes);
                taille_requete = nbLignes;
            });
            document.getElementById("requete_sql").value = requete;
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
    if (taille_requete != 0){
        position_dans_table = Math.min(taille_requete-1, position_dans_table)
    }
    affichage_table(false);
}

function appliquer_filtre(colonne){
    let valeur_filtre = document.getElementById("filtre_" + colonne).textContent;
    filtres.set(colonne, valeur_filtre)
    if(requete.includes("WHERE")){
        let demi_requetes = requete.split("WHERE");
        let conditions = demi_requetes[1].split("AND");
        for(let i = 0; i < conditions.length; i++){
            if (conditions[i].includes(colonne)){
                conditions.splice(i, 1);
            }
        }
        conditions.push(colonne + " LIKE \"%" + valeur_filtre + "%\"")
        requete = demi_requetes[0] + " WHERE " + conditions.join(" AND ");
    }else{
        requete += " WHERE " + colonne + " LIKE \"%" + valeur_filtre + "%\""; 
    }
    position_dans_table = 0;
    affichage_table(true);
}

document.getElementById("emplacement_table").focus()
document.onkeydown = function (e) {
    switch (e.code){
        case "ArrowDown":
            tourner_page(0, 1);
            break;
        case "ArrowUp":
            tourner_page(0, -1);
            break;
    }
};

function nouvelle_recherche_sql(){
    requete = document.getElementById("requete_sql").value;
    position_dans_table = 0;
    affichage_table(true);
    document.getElementById("changement_requete").removeAttribute("open");
}

function nouvelle_recherche_click_bouton(){
    let selecteurTable = document.getElementById("choix_table");
    requete = "SELECT * FROM " + selecteurTable.value;
    position_dans_table = 0;
    affichage_table(true);
    document.getElementById("changement_requete").removeAttribute("open");
}

function changer_type_recherche(){
    let sql = document.getElementById("zone_recherche_sql");
    let click_bouton = document.getElementById("zone_recherche_click_bouton");
    if (click_bouton.style.display == "none"){
        sql.style.display = "none";
        click_bouton.style.display = "inline";
    }else{
        sql.style.display = "inline";
        click_bouton.style.display = "none";
    }
}