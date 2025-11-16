let noeudEnCoursDeModif;
let fonctionCallback;

// create an array with nodes
let nodes = new vis.DataSet([
    { id: 1, shape:'image', image:'../assets/images/votre_reseau.png' },
    { id: 2, shape:'image', image:'../assets/images/internet.png' },
]);

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
    addNode: true,
    editNode: (donneesNoeud, callback) => ajout_modif_noeud(donneesNoeud, callback)
}
};
let network = new vis.Network(container, data, options);

function ajout_modif_noeud(donneesNoeud, callback){
    document.getElementById("popup-config-objet").style.display = "block";
    document.getElementById("fond-popup").style.display = "block";
    noeudEnCoursDeModif = donneesNoeud;
    fonctionCallback = callback;
}

function validerDonneesNoeud(){
    document.getElementById("popup-config-objet").style.display = "none";
    document.getElementById("fond-popup").style.display = "none";
    // Ajout d’une légende avec la liste des adresses
    let div_adresses = document.getElementById("adresses_equipement");
    let liste_adresses = div_adresses.children
    let texte_adresses = ""
    for(let i = 0; i < liste_adresses.length-1; i++){
        let liste_valeurs = liste_adresses.item(i).children 
        if(liste_valeurs.length > 1){
            if(liste_valeurs.item(0).value != ""){
                texte_adresses += "\n" + liste_valeurs.item(0).value;
                if(liste_valeurs.item(1).value != ""){
                    texte_adresses += " (" + liste_valeurs.item(1).value + ")";
                }
            }
        }
    }
    noeudEnCoursDeModif.label = texte_adresses;
    // Modification de l’apparence en fonction du type d’objet
    let type_equipement = document.getElementById("nature_equipement").value;
    noeudEnCoursDeModif.shape = "image";
    noeudEnCoursDeModif.image = "../assets/images/" + type_equipement + ".png";
    fonctionCallback(noeudEnCoursDeModif);
}