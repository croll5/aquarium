let nodes = new vis.DataSet([
]);

// create an array with edges
let edges = new vis.DataSet([
]);

let infos_machines = {}

// create a network
let container = document.getElementById("vue_generale");
let network;
let noeud_selectionne;

document.addEventListener("click", function(e){
    document.getElementById("menu_contextuel").style.display = "none";
})

parent.window.go.main.App.GetAquaConfig().then(aquaConfig => {
    let noeuds = []
    // On ajoute les machines
    for(let [idMachine, infosMachine] of Object.entries(aquaConfig["machines"])){
        noeuds.push({
            id:idMachine, 
            shape:"image", 
            image:'../assets/images/' + infosMachine.nature_equipement + '.png', x:infosMachine.x, y:infosMachine.y,
            label:getNomMachine(infosMachine)
        });
        console.log(infosMachine);
        infos_machines[idMachine] = {nom:infosMachine["nom"], a_analyser:infosMachine["a_analyser"]}; 
    }
    nodes = new vis.DataSet(noeuds);
    // On ajoute les liens entre les machines
    let cables = [];
    for(let [idCable, infosLien] of Object.entries(aquaConfig["liaisons_reseau"])){
        cables.push({
            id:idCable, 
            from:infosLien.source,
            to:infosLien.destination
        })
    }
    edges = new vis.DataSet(cables);
    let data = {
        nodes: nodes,
        edges: edges,
    };
    let options = {
        locale: "fr"
    };
    network = new vis.Network(container, data, options);
    network.on("oncontext", function (params) {
        // On enregistre le noeud sélectionné
        noeud_selectionne = network.getNodeAt(params.pointer.DOM);
        if(noeud_selectionne != undefined && infos_machines[noeud_selectionne].a_analyser){
            // On affiche le menu contextuel
            let menu = document.getElementById("menu_contextuel")
            menu.style.display = "flex";
            menu.style.left = params.event.clientX + "px";
            menu.style.top = params.event.clientY + "px";
        }
        params.event.preventDefault();
    })
})

function getNomMachine(param){
    console.log(param)
    let resultat = param.nom;
    for(let i in param.adresses){
        resultat += "\n" + param.adresses[i].adresse
        if(param.adresses[i].commentaire != ""){
            resultat += " (" + param.adresses[i].commentaire + ")";
        }
    }
    return resultat
}

function affichage_page_personnalisee(nomPage) {
    parent.creer_onglet('html/' + nomPage + '.html?machine=' + encodeURI(noeud_selectionne) + "&nom_machine=" + encodeURI(infos_machines[noeud_selectionne].nom), infos_machines[noeud_selectionne].nom + ' - ' + nomPage)
}