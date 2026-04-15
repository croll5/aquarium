let dossier_selectionne;
let selection_dossiers = {}
let config_dossier_a_enlever = "";

let params = new URLSearchParams(document.location.search);


document.getElementById("nom_machine").textContent = params.get("nom_machine");

document.addEventListener("click", function (event){
    let menus = document.getElementsByClassName("menu_contextuel");
    for(let [_, menu] of Object.entries(menus)){
        menu.style.display = "none";
    }
})

parent.window.go.main.App.ListeFichiersAnalysables(params.get("machine")).then(resultat =>{
    let emplacement_arbo = document.getElementById("arborescence");
    ajouter_contenu_dossier(emplacement_arbo, resultat);
})

parent.window.go.main.App.ListeConfigExtractionsDisponibles().then(resultat =>{
    console.log(resultat);
    let select_extractions = document.getElementById("select_config_extraction");
    resultat.forEach(nomConfig => {
        let nouvelle_option = document.createElement("option");
        nouvelle_option.textContent = nomConfig.replaceAll("_", " ").replaceAll(".xml", "");
        nouvelle_option.value = nomConfig.replaceAll(".xml", "");
        select_extractions.appendChild(nouvelle_option);
    });
})
 
function ajouter_contenu_dossier(contenant, ajouts){
    if (ajouts.DossiersEnfants != undefined){
        for(let [str_nom_dossier, dossier] of Object.entries(ajouts.DossiersEnfants)){
            let nouveau_dossier = document.createElement("details");
            nouveau_dossier.classList.add("dossier_arborescence");
            let nom_dossier = document.createElement("summary");
            nouveau_dossier.appendChild(nom_dossier);
            nom_dossier.textContent = str_nom_dossier;
            contenant.appendChild(nouveau_dossier);
            ajouter_contenu_dossier(nouveau_dossier, dossier);
            nom_dossier.oncontextmenu = function (event) {
                afficher_menu_contextuel(event, nouveau_dossier);
            }
        }
    }
    if (ajouts.Fichiers != undefined){
        contenant.setAttribute("fichiers", ajouts.Fichiers);
        let limite = 8;
        for(let nom_fichier of ajouts.Fichiers){
            let affichage_fichier = document.createElement("div");
            affichage_fichier.textContent = nom_fichier;
            affichage_fichier.classList.add("fichier_arborescence");
            affichage_fichier.oncontextmenu = function (event){
                afficher_menu_contextuel(event, contenant);
            }
            contenant.appendChild(affichage_fichier);
            if (limite < 0){
                break
            }
            limite--;
        }
    }
    if (ajouts.NbFichiers > 10){
        let indication_nb_fichiers = document.createElement("i");
        indication_nb_fichiers.textContent = "... et " + (ajouts.NbFichiers-10) + " autres fichiers";
        contenant.appendChild(indication_nb_fichiers);
    }
}

function afficher_menu_contextuel(event, nouveau_dossier){
    event.preventDefault();
    let menu_a_afficher;
    let autre_menu;
    let deja_selectionne = false;
    let dossier_a_regarder = nouveau_dossier;
    while (dossier_a_regarder.classList.contains("dossier_arborescence")){
        for(let [nom_config, contenu] of Object.entries(selection_dossiers)){
            if(contenu.includes(dossier_a_regarder)){
                deja_selectionne = true;
                config_dossier_a_enlever = nom_config;
                break;
            }
        }
        if (deja_selectionne){
            break
        }
        dossier_a_regarder = dossier_a_regarder.parentElement;
        if (dossier_a_regarder == undefined){
            break;
        }
    }
    if(deja_selectionne){
        menu_a_afficher = document.getElementById("menu_suppression");
        autre_menu = document.getElementById("menu_ajout");
    }else{
        menu_a_afficher = document.getElementById("menu_ajout");
        autre_menu = document.getElementById("menu_suppression");
    }
    autre_menu.style.display = "none";
    menu_a_afficher.style.top = event.pageY + "px";
    menu_a_afficher.style.left = event.pageX + "px";
    menu_a_afficher.style.display = "flex";
    dossier_selectionne = nouveau_dossier;
}

function ajout_dossier(){
    dossier_selectionne.classList.add("dossier_selectionne");
    let nom_config = document.getElementById("select_config_extraction").value;
    if(selection_dossiers[nom_config] == undefined){
        selection_dossiers[nom_config] = [];
    }
    selection_dossiers[nom_config].push(dossier_selectionne);
    fermer_popup('popup_select_config_extraction');
    console.log(selection_dossiers);
}

function retrait_dossier(){
    let dossier_a_supprimer = dossier_selectionne;
    while(dossier_a_supprimer.classList.contains("dossier_arborescence")){
        dossier_a_supprimer.classList.remove("dossier_selectionne");
        const index = selection_dossiers[config_dossier_a_enlever].indexOf(dossier_a_supprimer);
        if (index != -1){
            selection_dossiers[config_dossier_a_enlever].splice(index,1);
        }
        dossier_a_supprimer = dossier_a_supprimer.parentElement;
        if (dossier_a_supprimer == undefined){
            break;
        }
    }
    if(selection_dossiers[config_dossier_a_enlever].length == 0){
        delete selection_dossiers[config_dossier_a_enlever];
    }
}

function afficher_popup_select_fichiers(){
    document.getElementById("popup_select_fichiers").style.display = "block";
    document.getElementById("fond_popup").style.display = "block";
    let liste_fichiers = dossier_selectionne.getAttribute("fichiers").split(",");
    let div_liste_fichiers = document.getElementById("fichiers_filtres");
    div_liste_fichiers.textContent = "";
    for(let fichier of liste_fichiers){
        let p_fichier = document.createElement("li");
        p_fichier.textContent = fichier;
        p_fichier.classList.add("fichier_liste");
        p_fichier.onclick = function(event){
            p_fichier.classList.remove("fichier_liste");
            p_fichier.classList.add("fichier_non_compris");
        }
        div_liste_fichiers.appendChild(p_fichier);
    }
    document.getElementById("nom_dossier_cible").textContent = dossier_selectionne.getElementsByTagName("summary")[0].textContent;
}

function afficher_popup_config_extraction(){
    document.getElementById("popup_select_config_extraction").style.display = "block";
    document.getElementById("fond_popup").style.display = "block";
}

function valider_infos_config(){
    document.getElementById("arborescence").style.display = "inline";
    document.getElementById("infos_peripheriques").style.display = "none";
}

function actualisation_fichiers_filtres(){
    let fichiers_filtres = document.getElementById("fichiers_filtres").childNodes;
    // Obtention de l’expression régulière entrée dans la zone de 
    let valeurRecherchee = document.getElementById("filtre_choix_fichiers").value;
    let regex = new RegExp("^" + valeurRecherchee + "$"); 
    fichiers_filtres.forEach(fichier =>{
        if(regex.test(fichier.textContent)){
            console.log(fichier);
            fichier.classList.remove("fichier_non_compris");
        }else{
            fichier.classList.add("fichier_non_compris");
        }
    })
}