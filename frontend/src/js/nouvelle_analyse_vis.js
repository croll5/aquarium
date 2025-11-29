let noeudEnCoursDeModif;
let fonctionCallback;

// create an array with nodes
let nodes = new vis.DataSet([
    { id: 1, shape:'image', image:'../assets/images/internet.png' },
]);

let contenu_noeuds = {};

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
    editNode: (donneesNoeud, callback) => ajout_modif_noeud(donneesNoeud, callback, false)
}
};
let network = new vis.Network(container, data, options);

function ajout_modif_noeud(donneesNoeud, callback, ajout){
    if(donneesNoeud.id == 1){
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
    let div_adresses = document.getElementById("adresses_equipement");
    for(const adresse of contenu_noeuds[id]["adresses"]){
        let bloc = div_adresses.children.item(div_adresses.children.length-1)
        div_adresses.appendChild(bloc.cloneNode(true))
        bloc.children.item(0).value = adresse["adresse"];
        bloc.children.item(1).value = adresse["commentaire"];
    }
}

function validerDonneesNoeud(){
    document.getElementById("popup-config-objet").style.display = "none";
    document.getElementById("fond-popup").style.display = "none";
    contenu_noeuds[noeudEnCoursDeModif["id"]] = {};
    // Ajout d’une légende avec la liste des adresses
    let texte_adresses = recupere_liste_adresses();
    noeudEnCoursDeModif.label = texte_adresses;
    // Modification de l’apparence en fonction du type d’objet
    let type_equipement = document.getElementById("nature_equipement").value;
    if(type_equipement == ""){
        fonctionCallback(null);
    }
    contenu_noeuds[noeudEnCoursDeModif["id"]]["nature_equipement"] = type_equipement
    document.getElementById("nature_equipement").value = "";
    noeudEnCoursDeModif.shape = "image";
    noeudEnCoursDeModif.image = "../assets/images/" + type_equipement + ".png";
    fonctionCallback(noeudEnCoursDeModif);
}

function recupere_liste_adresses(){
    let div_adresses = document.getElementById("adresses_equipement");
    let liste_adresses = div_adresses.children
    let texte_adresses = ""
    let table_adresses = []
    for(let i = 0; i < liste_adresses.length-1; i++){
        let liste_valeurs = liste_adresses.item(i).children 
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
            contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"] = resultat;
        }else{
            for (const chemin of resultat) {
                if(!(chemin in contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"])){
                    contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"].push(chemin);
                }
            }
        }
        let affichage = document.getElementById("liste_fichiers_analyses")
        affichage.innerHTML = ""
        for (let i = 0; i < contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"].length; i++){
            let item = document.createElement("i");
            let texte = contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"][i].split("\\");
            item.textContent = "✅ " + texte[texte.length-1];
            item.value = contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"][i];
            affichage.appendChild(item);
            let bouton = document.createElement('button');
            bouton.textContent = "❌";
            bouton.classList.add("bouton_invisible");
            bouton.onclick = () =>{
                let index_valeur = contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"].indexOf(item.value);
                contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"].splice(index_valeur,1);
                item.remove();
            }
            item.appendChild(bouton);
            if (i < contenu_noeuds[noeudEnCoursDeModif["id"]]["fichiers"].length){
                let br = document.createElement("br");
                item.appendChild(br);
            }
        }
    })
}