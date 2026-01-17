const PROCHAINE_ETAPE = "etape_3";

// Ajout de configurations disponibles

let select_configuration = document.getElementById("fichier_config_analyse");
parent.window.go.main.App.ListeConfigurationsDisponibles().then(liste_config => {
    for (nom_config of liste_config){
        let option = document.createElement("option");
        option.value = nom_config;
        option.textContent = nom_config.replace(".xml", "");
        select_configuration.appendChild(option);
    }
})

// Création du réseau

let noeudEnCoursDeModif;
let fonctionCallback;

// create an array with nodes
let nodes = new vis.DataSet([
    { id: 1, shape:'image', image:'../assets/images/internet.png' },
]);

let contenu_noeuds = {};
let liste_liens = {};

// create an array with edges
let edges = new vis.DataSet([
]);

// create a network
let container = document.getElementById("vue_generale");
let data = {
nodes: nodes,
edges: edges,
};
let options = {
locale: "fr",
manipulation: {
    addNode: (donneesNoeud, callback) => ajout_modif_noeud(donneesNoeud, callback, true),
    editNode: (donneesNoeud, callback) => ajout_modif_noeud(donneesNoeud, callback, false),
    deleteNode: (donneesNoeud, callback) => suppr_noeud(donneesNoeud, callback),
    addEdge: (donneesNoeud, callback) => ajout_liaison(donneesNoeud, callback),
    deleteEdge:(donneesNoeud, callback) => suppr_noeud(donneesNoeud, callback)
}
};
let network = new vis.Network(container, data, options);

function ajout_liaison(donneesNoeud, callback){
    if(donneesNoeud["from"] == donneesNoeud["to"]){
        alert("Vous n’avez pas compris la consigne ! 🙃\nVous devez cliquer sur le premier objet, maintenir la souris enfoncée, aller jusqu’au deuxième objet et relâcher la souris.");
        callback(null);
    }else{
        liste_liens[donneesNoeud["id"]] = {"from":donneesNoeud["from"], "to":donneesNoeud["to"]}
        callback(donneesNoeud);
    }
}

function suppr_noeud(donneesNoeud, callback){
    for(let idNoeud of donneesNoeud["nodes"]){
        if(idNoeud == 1){
            alert("Vous ne pouvez la supprimer Internet 😅");
            callback(null);
            return
        }
    }
    if(donneesNoeud["nodes"].length > 0){
        if(confirm("Voulez-vous vraiment supprimer cet appareil ? 😧")){
            for(let id of donneesNoeud["nodes"]){
                delete contenu_noeuds[id];
            }
            for(let id of donneesNoeud["edges"]){
                delete liste_liens[id];
            }
            callback(donneesNoeud)
        }else{
            callback(null)
        }
    }
    for(let id of donneesNoeud["nodes"]){
        delete contenu_noeuds[id];
    }
    for(let id of donneesNoeud["edges"]){
        delete liste_liens[id];
    }
    verifier_archi_utilisable();
    callback(donneesNoeud)
}

function ajout_modif_noeud(donneesNoeud, callback, ajout){
    if(donneesNoeud.id == 1){
        alert("Vous ne pouvez la modifier Internet 😅");
        callback(null);
    }else{
        if(!ajout){
            preremplissage_noeud(donneesNoeud["id"]);
        }
        document.getElementById("popup-config-objet").style.display = "block";
        document.getElementById("fond-popup").style.display = "block";
        noeudEnCoursDeModif = donneesNoeud;
        fonctionCallback = callback;
    }
}

function preremplissage_noeud(id){
    document.getElementById("nature_equipement").value = contenu_noeuds[id]["nature_equipement"];
    document.getElementById("fichier_config_analyse").value = contenu_noeuds[id]["config"];
    document.getElementById("nom_machine").value = contenu_noeuds[id]["nom"];
    let div_adresses = document.getElementById("adresses_equipement");
    for(const adresse of contenu_noeuds[id]["adresses"]){
        let bloc = div_adresses.children.item(div_adresses.children.length-1);
        div_adresses.appendChild(bloc.cloneNode(true));
        bloc.children.item(0).value = adresse["adresse"];
        bloc.children.item(1).value = adresse["commentaire"];
        let bouton = document.createElement('button');
        bouton.textContent = "🗑️";
        bouton.classList.add("bouton_invisible");
        bouton.onclick = () =>{
            bloc.remove();
        }
        bloc.appendChild(bouton)
    }
    if (contenu_noeuds[id]["a_analyser"]){
        console.log(contenu_noeuds[id]);
        document.getElementById("oui_ana").click();
    }
    afficher_fichiers_selectionnes(id);
}

function validerDonneesNoeud(){
    if(contenu_noeuds[noeudEnCoursDeModif["id"]] == null){
        contenu_noeuds[noeudEnCoursDeModif["id"]] = {};
    }
    // Récupération du nom de l’appareil
    let nom_appareil = document.getElementById("nom_machine");
    if (nom_appareil.value == ""){
        alert("Vous devez donner un nom à cet appareil");
        return
    } 
    contenu_noeuds[noeudEnCoursDeModif["id"]]["nom"] = nom_appareil.value;
    nom_appareil.value = "";
    // Ajout d’une légende avec la liste des adresses
    let texte_adresses = recupere_liste_adresses();
    noeudEnCoursDeModif.label = texte_adresses;
    // Récupération des informations sur l’analyse
    recuperer_informations_analyse();
    // Modification de l’apparence en fonction du type d’objet
    let type_equipement = document.getElementById("nature_equipement").value;
    if(type_equipement == ""){
        alert("Vous devez renseigner un type d’équipement.");
        return
    }
    contenu_noeuds[noeudEnCoursDeModif["id"]]["nature_equipement"] = type_equipement
    document.getElementById("nature_equipement").value = "";
    noeudEnCoursDeModif.shape = "image";
    noeudEnCoursDeModif.image = "../assets/images/" + type_equipement + ".png";
    document.getElementById("popup-config-objet").style.display = "none";
    document.getElementById("fond-popup").style.display = "none";
    verifier_archi_utilisable();
    fonctionCallback(noeudEnCoursDeModif);
}

function recuperer_informations_analyse() {
    // On regarde s’il y a une analyse ou non
    contenu_noeuds[noeudEnCoursDeModif["id"]]["a_analyser"] = document.getElementById("oui_ana").checked;
    document.getElementById("non_ana").click();
    // On remet le sélecteur de fichiers à « aucun fichier sélectionné »
    let affichage = document.getElementById("liste_fichiers_analyses");
    affichage.innerHTML = "Aucun fichier sélectionné...";
    // On récupère le fichier de configuration choisi
    contenu_noeuds[noeudEnCoursDeModif["id"]]["config"] = document.getElementById("fichier_config_analyse").value;
    document.getElementById("fichier_config_analyse").value = "AUCUNE_SELECTION";
}

function recupere_liste_adresses(){
    let div_adresses = document.getElementById("adresses_equipement");
    let liste_adresses = div_adresses.children;
    let texte_adresses = "";
    if(contenu_noeuds[noeudEnCoursDeModif["id"]] != undefined){
        texte_adresses = contenu_noeuds[noeudEnCoursDeModif["id"]]["nom"]
    }
    let table_adresses = [];
    for(let i = 0; i < liste_adresses.length-1; i++){
        let liste_valeurs = liste_adresses.item(i).children;
        if(liste_valeurs.length > 1){
            if(liste_valeurs.item(0).value != ""){
                table_adresses.push({"adresse":liste_valeurs.item(0).value, "commentaire":liste_valeurs.item(1).value});
                texte_adresses += "\n" + liste_valeurs.item(0).value;
                if(liste_valeurs.item(1).value != ""){
                    texte_adresses += " (" + liste_valeurs.item(1).value + ")";
                }
            }
            liste_valeurs.item(0).value = "";
            liste_valeurs.item(1).value = "";
        }
    }
    contenu_noeuds[noeudEnCoursDeModif["id"]]["adresses"] = table_adresses;
    // Suppression des adresses en trop
    while(liste_adresses.length > 1){
        liste_adresses.item(0).remove();
    }
    return texte_adresses
}

function selection_fichier(){
    parent.window.go.main.App.ChoisirFichier("Sélectionnez les fichiers à analyser").then(resultat =>{
        if(contenu_noeuds[noeudEnCoursDeModif["id"]] == null){
            contenu_noeuds[noeudEnCoursDeModif["id"]] = {};
            if(typeof(resultat) == String){
                contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"] = [resultat];
            }else{
                contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"] = resultat;
            }
        }else{
            for (const chemin of resultat) {
                if(!(chemin in contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"])){
                    contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"].push(chemin);
                }
            }
        }
        afficher_fichiers_selectionnes(noeudEnCoursDeModif["id"])
    })
}

function afficher_fichiers_selectionnes(idNoeud) {
    let affichage = document.getElementById("liste_fichiers_analyses");
    affichage.innerHTML = "";
    if (contenu_noeuds[idNoeud]["fichiers"] == null){
        contenu_noeuds[idNoeud]["fichiers"] = []
    }
    for (let i = 0; i < contenu_noeuds[idNoeud]["fichiers"].length; i++){
        let item = document.createElement("i");
        let texte = contenu_noeuds[idNoeud]["fichiers"][i].split("\\");
        item.textContent = "✅ " + texte[texte.length-1];
        item.value = contenu_noeuds[idNoeud]["fichiers"][i];
        affichage.appendChild(item);
        let bouton = document.createElement('button');
        bouton.textContent = "🗑️";
        bouton.classList.add("bouton_invisible");
        bouton.onclick = () =>{
            let index_valeur = contenu_noeuds[idNoeud]["fichiers"].indexOf(item.value);
            contenu_noeuds[idNoeud]["fichiers"].splice(index_valeur,1);
            item.remove();
        }
        item.appendChild(bouton);
        if (i < contenu_noeuds[idNoeud]["fichiers"].length){
            let br = document.createElement("br");
            item.appendChild(br);
        }
    }
    if (affichage.innerHTML == ""){
        affichage.innerHTML = "Aucun fichier sélectionné...";
    }
}

function verifier_archi_utilisable(){
    let etape_suivante = document.getElementById(PROCHAINE_ETAPE);
    let utilisable = false;
    if(contenu_noeuds.length < 1){
        utilisable = false;
    }else{
        for(let id_noeud of Object.keys(contenu_noeuds)){
            if (contenu_noeuds[id_noeud]["fichiers"] != null && contenu_noeuds[id_noeud]["fichiers"].length > 0){
                utilisable = true;
                break;
            }
        }
    }
    if(utilisable && etape_suivante.hasAttribute("disabled")){
        etape_suivante.removeAttribute("disabled");
    } else if (!utilisable){
        etape_suivante.setAttribute("disabled", true);
    }
}

function get_donnees_reseau(){
    return contenu_noeuds;
}