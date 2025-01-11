let position_dans_table = 0;
let position_debut_recuperation = 0;
let requete = "SELECT extracteur, horodatage, message, source FROM chronologie";
let taille_requete = 0;
let valeurs_filtres = new Map();
let consignes_filtres = new Map();
let order_by = "riendutout";
let tableRecuperee = new Object();

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
    if(majTaille || position_dans_table > position_debut_recuperation + 995 || position_dans_table < position_debut_recuperation){
        if(position_dans_table > position_debut_recuperation + 995 || position_dans_table < position_debut_recuperation){
            position_debut_recuperation = Math.max(0,position_dans_table - 500);
        }
        document.body.style.cursor = "wait";
        parent.window.go.main.App.ResultatRequeteSQLExtraction(requete, position_debut_recuperation, 1000).then(resultat =>{
            document.body.style.cursor = "default";
            tableRecuperee = resultat;
            document.getElementById("indicateur_page").textContent = position_dans_table + "-" + (position_dans_table+5);
            emplacement_resultat.innerHTML = "";
            console.log(resultat);
            creer_tableau_depuis_dico(resultat.slice(position_dans_table - position_debut_recuperation, position_dans_table - position_debut_recuperation + 5), emplacement_resultat, true, valeurs_filtres, consignes_filtres, order_by);
            if (majTaille){
                parent.window.go.main.App.TailleRequeteSQLExtraction(requete).then(nbLignes =>{
                    console.log(nbLignes);
                    taille_requete = nbLignes;
                });
                document.getElementById("requete_sql").value = requete;
            }
        })
    }else{
        emplacement_resultat.innerHTML = "";
        console.log(position_dans_table - position_debut_recuperation);
        document.getElementById("indicateur_page").textContent = position_dans_table + "-" + (position_dans_table+5);
        creer_tableau_depuis_dico(tableRecuperee.slice(position_dans_table - position_debut_recuperation, position_dans_table - position_debut_recuperation + 5), emplacement_resultat, true, valeurs_filtres, consignes_filtres, order_by);
    }
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
    let valeur_filtre = document.getElementById("valeur_filtre_" + colonne).textContent;
    let consigne_filtre = document.getElementById("consigne_filtre_" + colonne).value;
    let avant_valeur = " LIKE \"%";
    let apres_valeur = "%\"";
    switch(consigne_filtre){
        case "🔎🔤":
            avant_valeur = " LIKE \"";
            break;
        case "🔤🔎":
            apres_valeur = "\"";
            break;
        case "🔤 = 🔎":
            avant_valeur = " = \"";
            apres_valeur = "\"";
            break;
        case "🔤 > 🔎":
            avant_valeur = " > \"";
            apres_valeur = "\"";
            break;
        case "🔤 < 🔎":
            avant_valeur = " < \"";
            apres_valeur = "\"";
    }
    let changement_valeur = (valeur_filtre == "" && !valeurs_filtres.has(colonne)) || (valeurs_filtres.has(colonne) && valeurs_filtres.get(colonne) == valeur_filtre);
    let changement_consigne = (consigne_filtre == "🔤🔎🔤" && !consignes_filtres.has(colonne)) || (consignes_filtres.has(colonne) && consignes_filtres.get(colonne) == consigne_filtre);
    if(changement_valeur && changement_consigne){
        return;
    }
    valeurs_filtres.set(colonne, valeur_filtre);
    consignes_filtres.set(colonne, consigne_filtre);
    if(requete.includes(" WHERE ")){
        let demi_requetes = requete.split(" WHERE ");
        let filtrage_order_by = demi_requetes[1].split(" ORDER ")
        let conditions = filtrage_order_by[0].split(" AND ");
        for(let i = 0; i < conditions.length; i++){
            if (conditions[i].includes(colonne)){
                conditions.splice(i, 1);
            }
        }
        if (valeur_filtre != ""){
            conditions.push(colonne + avant_valeur + valeur_filtre + apres_valeur)
        }
        if (conditions.length > 0 ){
            requete = demi_requetes[0] + " WHERE " + conditions.join(" AND ");
        }else{
            requete = demi_requetes[0];
        }
        if (filtrage_order_by.length > 1){
            requete += " ORDER" + filtrage_order_by[1];
        }
    }else if (requete.includes(" ORDER ")){
        let demi_requetes = requete.split(" ORDER ");
        requete = demi_requetes[0] + " WHERE " + colonne + avant_valeur + valeur_filtre + apres_valeur + " ORDER " + demi_requetes[1]
    }else{
        requete += " WHERE " + colonne + avant_valeur + valeur_filtre + apres_valeur; 
    }
    position_dans_table = 0;
    position_debut_recuperation = 0;
    affichage_table(true);
}

function trier_par(colonne){
    if (order_by != colonne){
        order_by = colonne;
        requete = requete.split(" ORDER")[0];
        requete += " ORDER BY " + colonne;
        affichage_table(true);
    }
    else{
        order_by = "riendutout";
        requete = requete.split(" ORDER")[0];
        affichage_table(true);
    }
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
        case "ArrowLeft":
            tourner_page(0, -5);
            break;
        case "ArrowRight":
            tourner_page(0, 5);
            break;
    }
};

function nouvelle_recherche_sql(){
    valeurs_filtres.clear();
    requete = document.getElementById("requete_sql").value;
    position_dans_table = 0;
    position_debut_recuperation = 0;
    affichage_table(true);
    document.getElementById("changement_requete").removeAttribute("open");
}

function nouvelle_recherche_click_bouton(){
    valeurs_filtres.clear();
    let selecteurTable = document.getElementById("choix_table");
    requete = "SELECT * FROM " + selecteurTable.value;
    position_dans_table = 0;
    position_debut_recuperation = 0;
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