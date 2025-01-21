/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025) 

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à [rappeler les
caractéristiques techniques de votre logiciel]. 

Ce logiciel est régi par la licence CeCILL soumise au droit français et
respectant les principes de diffusion des logiciels libres. Vous pouvez
utiliser, modifier et/ou redistribuer ce programme sous les conditions
de la licence CeCILL telle que diffusée par le CEA, le CNRS et l'INRIA 
sur le site "http://www.cecill.info".

En contrepartie de l'accessibilité au code source et des droits de copie,
de modification et de redistribution accordés par cette licence, il n'est
offert aux utilisateurs qu'une garantie limitée.  Pour les mêmes raisons,
seule une responsabilité restreinte pèse sur l'auteur du programme,  le
titulaire des droits patrimoniaux et les concédants successifs.

A cet égard  l'attention de l'utilisateur est attirée sur les risques
associés au chargement,  à l'utilisation,  à la modification et/ou au
développement et à la reproduction du logiciel par l'utilisateur étant 
donné sa spécificité de logiciel libre, qui peut le rendre complexe à 
manipuler et qui le réserve donc à des développeurs et des professionnels
avertis possédant  des  connaissances  informatiques approfondies.  Les
utilisateurs sont donc invités à charger  et  tester  l'adéquation  du
logiciel à leurs besoins dans des conditions permettant d'assurer la
sécurité de leurs systèmes et ou de leurs données et, plus généralement, 
à l'utiliser et l'exploiter dans les mêmes conditions de sécurité. 

Le fait que vous puissiez accéder à cet en-tête signifie que vous avez 
pris connaissance de la licence CeCILL, et que vous en avez accepté les
termes.
*/

let position_dans_table = 0;
let position_debut_recuperation = 0;
let requete = "SELECT id, extracteur, horodatage, message, source FROM chronologie";
let taille_requete = 0;
let valeurs_filtres = new Map();
let consignes_filtres = new Map();
let order_by = "riendutout";
let tableRecuperee = new Object();
let liste_id_a_enregistrer = Array(0);
let evenements_a_enregistrer = new Map();

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
    if(majTaille){
        liste_id_a_enregistrer = Array(0);
        evenements_a_enregistrer.clear();
    }
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
            creer_tableau_depuis_dico(resultat.slice(position_dans_table - position_debut_recuperation, position_dans_table - position_debut_recuperation + 5), emplacement_resultat, true, valeurs_filtres, consignes_filtres, order_by, position_dans_table, liste_id_a_enregistrer);
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
        document.getElementById("indicateur_page").textContent = position_dans_table + "-" + (position_dans_table+5);
        creer_tableau_depuis_dico(tableRecuperee.slice(position_dans_table - position_debut_recuperation, position_dans_table - position_debut_recuperation + 5), emplacement_resultat, true, valeurs_filtres, consignes_filtres, order_by, position_dans_table, liste_id_a_enregistrer);
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
            requete += " ORDER " + filtrage_order_by[1];
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
        requete = requete.split(" ORDER ")[0];
        requete += " ORDER BY " + colonne;
        affichage_table(true);
    }
    else{
        order_by = "riendutout";
        requete = requete.split(" ORDER ")[0];
        affichage_table(true);
    }
}

function enregistrement_id(id){
    for (let i = 0; i < liste_id_a_enregistrer.length; i++){
        if(liste_id_a_enregistrer[i] == id){
            liste_id_a_enregistrer.splice(i, 1);
            evenements_a_enregistrer.delete(id)
        }
    }
    if(document.getElementById("casacocher_" + id).checked){
        liste_id_a_enregistrer.push(id);
        evenements_a_enregistrer.set(id, tableRecuperee[id-position_debut_recuperation])
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

/* FONCTIONS D'ENREGISTREMENT DANS LE RAPPORT */

let divEnregistrementTable = document.getElementById("enregistrement_table");
divEnregistrementTable.addEventListener("toggle", (event) => {
    if (divEnregistrementTable.open) {
        let selecteurPiste = document.getElementById("choix_piste");
        selecteurPiste.innerHTML = "";
        parent.window.go.main.App.ListePistesRapport().then(resultat => {
            for (let ligne of resultat) {
                let nom_piste = document.createElement("option");
                nom_piste.textContent = ligne["titre"];
                nom_piste.value = ligne["id"];
                selecteurPiste.appendChild(nom_piste)
            }
        });
    }
  });

function enregistrer_table_dans_rapport(){
    if(liste_id_a_enregistrer.length == 0){
        alert("Vous devez sélectionner des lignes à enregistrer 🧐");
        return
    }
    let idPiste = document.getElementById("choix_piste").value;
    let commentaire = document.getElementById("commentaire_analyste").value;
    if(commentaire == ""){
        alert("Vous devez ajouter un commentaire sur ces évènements 🤓");
        return
    }
    let tableau_a_enregistrer = Array.from(evenements_a_enregistrer, ([_, valeur]) => valeur)
    parent.window.go.main.App.AjouterEtapeDansRapport(requete, tableau_a_enregistrer, idPiste, commentaire);
    document.getElementById("enregistrement_table").removeAttribute("open");
    document.getElementById("commentaire_analyste").value = "";
}