/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025) 

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l'analyse des collectes
traçologiques effectuées avec le logiciel DFIR-ORC. 

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

function ajouter_bouton_suppression(bloc){
    if(bloc.parentElement != null && bloc.parentElement.firstElementChild === bloc){
        return;
    }
    if(bloc.getElementsByClassName("bouton_suppression_ligne").length > 0){
        return;
    }
    let suppr = document.createElement("button");
    suppr.innerText = "🗑️";
    suppr.className = "bouton_invisible bouton_suppression_ligne";
    suppr.onclick = () => {
        let bouton_ajout_dans_bloc = bloc.querySelector(".bouton_ajout_liste");
        if(bouton_ajout_dans_bloc != null){
            let bloc_precedent = bloc.previousElementSibling;
            if(bloc_precedent != null && bloc_precedent.classList.contains("bloc_de_liste")){
                bloc_precedent.appendChild(bouton_ajout_dans_bloc);
            }
        }
        bloc.remove();
    };
    let bouton_ajout = bloc.getElementsByClassName("bouton_ajout_liste");
    if(bouton_ajout.length > 0){
        bloc.insertBefore(suppr, bouton_ajout.item(0));
    } else {
        bloc.appendChild(suppr);
    }
}

function ajouter_ligne_liste(id_section){
    let noeud_principal = document.getElementById(id_section);
    let liste_elements = noeud_principal.children;
    let numero_champ = liste_elements.length - 1;
    let element_de_base = liste_elements.item(numero_champ);
    let clone = element_de_base.cloneNode(true);
    let bouton_ajout = noeud_principal.querySelector(".bouton_ajout_liste");
    if(bouton_ajout != null){
        bouton_ajout.remove();
    }
    let bouton_ajout_clone = clone.getElementsByClassName("bouton_ajout_liste");
    if(bouton_ajout_clone.length > 0){
        bouton_ajout_clone.item(0).remove();
    }
    let boutons_suppression_clone = clone.getElementsByClassName("bouton_suppression_ligne");
    while(boutons_suppression_clone.length > 0){
        boutons_suppression_clone.item(0).remove();
    }
    vider_champs(clone);
    noeud_principal.appendChild(clone);
    if(bouton_ajout != null){
        clone.appendChild(bouton_ajout);
    }
    ajouter_bouton_suppression(clone);
}

function initialiser_suppression_lignes(id_section){
    let section = document.getElementById(id_section);
    if(section == null){
        return;
    }
    for(let i = 1; i < section.children.length; i++){
        let bloc = section.children.item(i);
        if(bloc.classList.contains("bloc_de_liste")){
            ajouter_bouton_suppression(bloc);
        }
    }
}

function ajuster_champs(id_section, ajout_auto = true){
    let noeud_principal = document.getElementById(id_section)
    let liste_elements = noeud_principal.children;
    let numero_champ = liste_elements.length - 1
    if(ajout_auto && section_remplie(liste_elements.item(numero_champ), false)){
        ajouter_ligne_liste(id_section);
    } else if(ajout_auto && numero_champ > 0 && !section_remplie(liste_elements.item(numero_champ - 1), false)){
        let liste_boutons = liste_elements.item(numero_champ-1).getElementsByClassName("bouton_suppression_ligne");
        for(bouton of liste_boutons){
            bouton.remove();
        }
        liste_elements.item(numero_champ).remove()
    }
}

window.addEventListener("load", () => {
    initialiser_suppression_lignes("main_courante");
    initialiser_suppression_lignes("liste_contacts");
});

function vider_champs(objet){
    let liste_elements = objet.children;
    for(let i = 0; i < liste_elements.length; i++){
        let element = liste_elements.item(i);
        if(liste_elements.item(i).hasChildNodes()){
            vider_champs(element);
        } else {
            element.value = "";
            if(element.hasAttribute("required")){
                element.removeAttribute("required");
            }
        }
    }
}

function section_remplie(objet, completement = true) {
    let contenu_element = objet.children;
    for(let i = 0; i < contenu_element.length; i++){
        if(contenu_element.item(i).hasChildNodes()){
            let remplissage_enfant = section_remplie(contenu_element.item(i), completement)
            if(completement && !remplissage_enfant){
                return false
            } else if(!completement && remplissage_enfant){
                return true
            }
        }else if(contenu_element.item(i).value != undefined && contenu_element.item(i).value != ""){
            if(!completement){
                return true;
            }
        } else if(completement && contenu_element.item(i).hasAttribute("required") && (contenu_element.item(i).value == undefined || contenu_element.item(i).value == "")){
                return false;
        }
    }
    return completement;
}

function afficher_bloc(a_afficher, id){
    let bloc = document.getElementById(id);
    if(a_afficher){
        bloc.style.display = "inline";
    } else{
        bloc.style.display = "none";
    }
}

function quitter_nouvelle_analyse(){
    window.location.replace("../html/accueil.html");
}

function verifier_remplissage(section, prochaine_etape){
    let element = document.getElementById(section);
    if(section_remplie(element, true)){
        if(document.getElementById(prochaine_etape).hasAttribute("disabled")){
            document.getElementById(prochaine_etape).removeAttribute("disabled");
        }
    } else{
        document.getElementById(prochaine_etape).setAttribute("disabled", true);
    }
}

function selection_dossier(id_paragraphe, id_input, id_section, id_suivant){
    parent.window.go.main.App.ChoisirDossier("Sélectionnez un dossier d’enregistrement").then(resultat => {
        if(resultat == ""){
            document.getElementById(id_paragraphe).textContent = "Aucun dossier sélectionné...";
            document.getElementById(id_input).value = ""; 
        } else{
            document.getElementById(id_paragraphe).textContent = resultat;
            document.getElementById(id_input).value = resultat; 
        }
        verifier_remplissage(id_section, id_suivant);
    })
}

function valider_creation_analyse(){
    let donnees_analyse = donnees_conf_analyse();
    donnees_analyse["Machines"] = get_donnees_reseau();
    donnees_analyse["liaisons_reseau"] = get_liens_reseau();
    document.getElementById("formulaire").style.display = "none";
    document.getElementById("patientez").style.display = "inline";
    parent.window.go.main.App.CreationNouveauProjet(donnees_analyse).then(resultat =>{
        if (resultat != "") {
            window.location.replace("../html/config_analyse.html");
            parent.document.getElementsByTagName("header")[0].style.display = "flex";
            let onglet_courant = parent.document.getElementById("onglet_extraction");
            onglet_courant.classList.add("onglet_selectionne");
        } else{
            document.getElementById("formulaire").style.display = "inline";
            document.getElementById("patientez").style.display = "none";
        }
    });
}

function donnees_conf_analyse(base = document, profondeur = 0){
    // Création de la variable résultat
    let resultat = {}
    // Récupération des données 
    let tags = ["input", "textarea", "select"];
    for(let tag of tags){
        let elements = base.getElementsByTagName(tag);
        for(let element of elements){
            if (element.hasAttribute("aqua_champ") && element.value != "" && (!element.hasAttribute("aqua_prof") || element.getAttribute("aqua_prof") == String(profondeur))){
                resultat[element.getAttribute("aqua_champ")] = element.value;
                if(element.hasAttribute("type") && element.getAttribute("type") == "datetime-local"){
                    resultat[element.getAttribute("aqua_champ")] = element.value + ":00Z";
                }
            }
        }
    }
    // Gestion des divs
    let divs = base.getElementsByTagName("div");
    for(let div of divs){
        if(div.hasAttribute("aqua_champ")){
            let nom_liste = div.getAttribute("aqua_champ");
            console.log(nom_liste);
            if(resultat[nom_liste] == null){
                resultat[nom_liste] = []
            }
            let donnees = donnees_conf_analyse(div, profondeur+1);
            if(Object.keys(donnees).length != 0){
                resultat[nom_liste].push(donnees);
            }
        }
    }
    return resultat
}
