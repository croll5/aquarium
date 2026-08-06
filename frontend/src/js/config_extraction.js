let dossier_selectionne;
let selection_dossiers = {}
let config_dossier_a_enlever = "";

let params = new URLSearchParams(document.location.search);

// Afficher le nom de la machine dans le titre de la page
document.getElementById("nom_machine").textContent = params.get("nom_machine");

// Enlever le menu contextuel quand on clique ailleurs
document.addEventListener("click", function (event){
    let menus = document.getElementsByClassName("menu_contextuel");
    for(let [, menu] of Object.entries(menus)){
        menu.style.display = "none";
    }
})

// Afficher les dossiers analysables
parent.window.go.main.App.ListeFichiersAnalysables(params.get("machine")).then(resultat =>{
    let emplacement_arbo = document.getElementById("arborescence");
    ajouter_contenu_dossier(emplacement_arbo, resultat);
})

// Ajouter la liste des configurations disponibles au menu déroulant
parent.window.go.main.App.ListeConfigExtractionsDisponibles().then(resultat =>{
    let select_extractions = document.getElementById("select_config_extraction");
    resultat.forEach(nomConfig => {
        let nouvelle_option = document.createElement("option");
        nouvelle_option.textContent = nomConfig.replaceAll("_", " ").replaceAll(".xml", "");
        nouvelle_option.value = nomConfig.replaceAll(".xml", "");
        select_extractions.appendChild(nouvelle_option);
    });
})
 
/** Fonction récursive qui ajoute le contenu des dossiers analysables
 * @param contenant : dossier à remplir
 * @param ajouts : éléments à ajouter dans le dossier
 * @returns : rien, ajoute des éléments sur la page HTML
**/
function ajouter_contenu_dossier(contenant, ajouts){
    // Ajouter les sous-dossiers
    if (ajouts.DossiersEnfants != undefined){
        for(let [str_nom_dossier, dossier] of Object.entries(ajouts.DossiersEnfants)){
            let nouveau_dossier = document.createElement("details");
            nouveau_dossier.classList.add("dossier_arborescence");
            if(dossier.Est7z){
                nouveau_dossier.setAttribute("aqua_archive", true);
            }
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
    // Ajouter les fichiers
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
    // Ajouter le nombre d'éléments s’ils ne peuvent pas être tous ajoutés
    if (ajouts.NbFichiers > 10){
        let indication_nb_fichiers = document.createElement("i");
        indication_nb_fichiers.textContent = "... et " + (ajouts.NbFichiers-10) + " autres fichiers";
        contenant.appendChild(indication_nb_fichiers);
    }
}

/*** Fonction qui affiche un menu contextuel différent en fonction de 
 * si le dossier est sélectionné ou non
 * @param event : l’évènement de click, avec notamment la position du menu
 * @param nouveau_dossier : le dossier sélectionné 
***/
function afficher_menu_contextuel(event, nouveau_dossier){
    event.preventDefault();
    let menu_a_afficher;
    let autre_menu;
    let deja_selectionne = false;
    let dossier_a_regarder = nouveau_dossier;
    if (dossier_a_regarder.classList.contains("fichier_arborescence")){
        dossier_a_regarder = dossier_a_regarder.parentElement;
    }
    for(let [nom_config, contenu] of Object.entries(selection_dossiers)){
        for(let i = 0; i < contenu.length; i++){
            if(contenu[i].dossier == dossier_a_regarder){
                deja_selectionne=true;
                config_dossier_a_enlever = {config:nom_config, index:i}
                break
            }
        }
        if(deja_selectionne){
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

/*** Fonction permettant de valider un filtre sur les fichiers 
 * d'un dossier
 */
function valider_selection_filtre(){
    fermer_popup("popup_select_fichiers");
    afficher_popup_config_extraction();
}

/** Fonction permettant d’ajouter un dossier à la sélection
 */
function ajout_dossier(){
    dossier_selectionne.classList.add("dossier_selectionne");
    let nom_config = document.getElementById("select_config_extraction").value;
    if(selection_dossiers[nom_config] == undefined){
        selection_dossiers[nom_config] = [];
    }
    let input_filtre = document.getElementById("filtre_choix_fichiers");
    selection_dossiers[nom_config].push({dossier:dossier_selectionne, filtre:input_filtre.value});
    input_filtre.value = "*";
    fermer_popup('popup_select_config_extraction');
    document.getElementById("select_config_extraction").value = "";
    
}

/** Fonction permettant de supprimer le dossier sélectionné */
function retrait_dossier(){
    dossier_selectionne.classList.remove("dossier_selectionne");
    selection_dossiers[config_dossier_a_enlever.config].splice(config_dossier_a_enlever.index, 1)
    if(selection_dossiers[config_dossier_a_enlever.config].length == 0){
        delete selection_dossiers[config_dossier_a_enlever.config];
    }
}

/** Fonctionn affichant les fichiers suivants au scroll */
function scroll_liste_fichiers(){
    let contenant_liste_fichiers = document.getElementById("contenant_fichiers_filtres");
    let scrollCourant = contenant_liste_fichiers.scrollTop;
    if(scrollCourant + contenant_liste_fichiers.getBoundingClientRect().height > 4*contenant_liste_fichiers.scrollHeight/5){
        let ul_liste_fichiers = document.getElementById("fichiers_filtres");
        let taille_liste_affichee = ul_liste_fichiers.childNodes.length;
        let liste_fichiers = dossier_selectionne.getAttribute("fichiers").split(",");
        let div_liste_fichiers = document.getElementById("fichiers_filtres");
        for(let i = taille_liste_affichee; i < taille_liste_affichee+5 && i < liste_fichiers.length; i++){
            let nouveau_fichier = document.createElement("li");
            nouveau_fichier.textContent = liste_fichiers[i];
            nouveau_fichier.classList.add("fichier_liste");
            div_liste_fichiers.appendChild(nouveau_fichier);
        }
    }
}

/** Fonction affichant la fenêtre contextuelle permettant 
 * de choisir un filtre à appliquer aux fichiers
 */
function afficher_popup_select_fichiers(filtre="*"){
    document.getElementById("popup_select_fichiers").style.display = "flex";
    document.getElementById("fond_popup").style.display = "block";
    document.getElementById("filtre_choix_fichiers").value = filtre;
    let liste_fichiers = dossier_selectionne.getAttribute("fichiers").split(",");
    let div_liste_fichiers = document.getElementById("fichiers_filtres");
    div_liste_fichiers.textContent = "";
    for(let i = 0; i < 50; i++){
        if(liste_fichiers.length <= i){
            break;
        }
        let p_fichier = document.createElement("li");
        p_fichier.textContent = liste_fichiers[i];
        p_fichier.classList.add("fichier_liste");
        div_liste_fichiers.appendChild(p_fichier);
    }
    document.getElementById("nom_dossier_cible").textContent = dossier_selectionne.getElementsByTagName("summary")[0].textContent;
    document.getElementById("contenant_fichiers_filtres").scrollTop = 0;
}

/** Affichage de la fenêtre contextuelle pour choisir la 
 * configuration d’extraction des éléments d’un dossier */
function afficher_popup_config_extraction(){
    document.getElementById("popup_select_config_extraction").style.display = "block";
    document.getElementById("fond_popup").style.display = "block";
}

/** Fonction permettant de valider les informations de la configruation */
function valider_infos_config(){
    document.getElementById("details_config").style.display = "inline";
    document.getElementById("infos_peripheriques").style.display = "none";
}

/** Application du filtre renseigné par l’utilisateur aux fichiers affichés */
function actualisation_fichiers_filtres(){
    let fichiers_filtres = document.getElementById("fichiers_filtres").childNodes;
    // Obtention de l’expression régulière entrée dans la zone de 
    let valeurRecherchee = document.getElementById("filtre_choix_fichiers").value;
    fichiers_filtres.forEach(fichier =>{
        parent.window.go.main.App.CorrespondanceCheminModele(fichier.textContent, valeurRecherchee).then(resultat =>{
            if(resultat){
                fichier.classList.remove("fichier_non_compris");
            }else{
                fichier.classList.add("fichier_non_compris");
            }
        })
    })
}

/** Fonction permettant de modifier le filtre appliqué aux fichiers 
 * du dossier sélectionné
 */
function modifier_filtres_fichiers(){
    let filtre_actuel = selection_dossiers[config_dossier_a_enlever.config][config_dossier_a_enlever.index].filtre;
    document.getElementById("select_config_extraction").value = config_dossier_a_enlever.config;
    afficher_popup_select_fichiers(filtre_actuel);
    actualisation_fichiers_filtres();
    retrait_dossier();
}

function modifier_config_dossier(){
    let filtre_actuel = selection_dossiers[config_dossier_a_enlever.config][config_dossier_a_enlever.index].filtre;
    let input_filtre = document.getElementById("filtre_choix_fichiers");
    input_filtre.value = filtre_actuel;
    document.getElementById("select_config_extraction").value = config_dossier_a_enlever.config;
    afficher_popup_config_extraction();
    retrait_dossier();
}

/** Fonction permettant d’enregistrer la nouvelle configuration */
function valider_configuration(){
    let nouvelle_config = conversion_config_html_objet()
    if(Object.keys(selection_dossiers).length == 0){
        alert("Vous n’avez ajouté aucun dossier à la configuration !😯\nVous pouvez le faire à l’aide d’un clic droit sur le dossier que vous souhaitez ajouter.");
        return
    }
    // Récupération du nom de la configuration
    let nom_config = document.getElementById("nom_config").value;
    nom_config = nom_config.replaceAll(" ", "_");
    let reutilisable = document.getElementById("permanente").checked;
    if(confirm("Voulez-vous enregistrer cette configuration ?🙃")){
        parent.window.go.main.App.EnregistrerConfigMachine(nouvelle_config, nom_config, reutilisable, params.get("machine")).then(resultat =>{
            let titre_cr = document.createElement("h1");
            titre_cr.textContent = "Enregistrement réussi 🐬";
            let texte_cr = document.createElement("p");
            texte_cr.textContent = "La configuration a bien été enregistrée. Vous pouvez aller à la vue d’ensemble pour lancer l’extraction des données. 🫧"
            parent.fermer_onglet_courant(titre_cr, texte_cr)
        })
    }
}

function conversion_config_html_objet(){
    let nouvelle_config = []
    for(let [config, dossiers] of Object.entries(selection_dossiers)){
        nouvelle_config.push({
            Id:config,
            Chemins:[]
        })
        for(const dossier of dossiers){
            let liste_dossiers = [];
            let num_archive = -1;
            let dossier_courant = dossier.dossier;
            let i = 0
            while(dossier_courant.classList.contains("dossier_arborescence")){
                liste_dossiers.push(dossier_courant.querySelector("summary").textContent);
                if(dossier_courant.hasAttribute("aqua_archive")){
                    num_archive = i;
                }
                dossier_courant = dossier_courant.parentElement;
                i++;
            }
            if(num_archive == -1){
                nouvelle_config.at(-1).Chemins.push({
                    Dossiers: liste_dossiers,
                    Fichier: dossier.filtre
                })
            }else{
                if(num_archive > 0){
                    dossier.filtre = liste_dossiers.slice(0, num_archive).join("/") + "/" + dossier.filtre;
                }
                nouvelle_config.at(-1).Chemins.push({
                    Dossiers: liste_dossiers.slice(num_archive+1, liste_dossiers.length),
                    Archive: liste_dossiers[num_archive],
                    Fichier: dossier.filtre
                })
            }
            
        }
    }
    return nouvelle_config
}