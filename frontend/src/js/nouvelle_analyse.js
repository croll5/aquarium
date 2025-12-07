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

function ajuster_champs(id_section){
    let noeud_principal = document.getElementById(id_section)
    let liste_elements = noeud_principal.children;
    let numero_champ = liste_elements.length - 1
    if(section_remplie(liste_elements.item(numero_champ), false)){
        let element_de_base = liste_elements.item(numero_champ);
        let clone = element_de_base.cloneNode(true);
        vider_champs(clone);
        noeud_principal.appendChild(clone);
        // Ajout d'un bouton « supprimer »
        let suppr = document.createElement("button");
        suppr.innerText = "❌";
        suppr.className = "bouton_invisible";
        suppr.onclick = (ev) => {
            element_de_base.remove();
        }
        element_de_base.appendChild(suppr);
    } else if(numero_champ > 0 && !section_remplie(liste_elements.item(numero_champ - 1), false)){
        let liste_boutons = liste_elements.item(numero_champ-1).getElementsByClassName("bouton_invisible");
        for(bouton of liste_boutons){
            bouton.remove();
        }
        liste_elements.item(numero_champ).remove()
    }
}

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
    console.log(donnees_analyse);
}

function donnees_conf_analyse(base = document, profondeur = 0){
    // Création de la variable résultat
    let resultat = {}
    // Gestion des inputs
    let inputs = base.getElementsByTagName("input");
    for(let input of inputs){
        if (input.hasAttribute("aqua_champ") && (!input.hasAttribute("aqua_prof") || input.getAttribute("aqua_prof") == String(profondeur))){
            resultat[input.getAttribute("aqua_champ")] = input.value;
        }
    }
    // Gestion des text area
    let textareas = base.getElementsByTagName("textarea");
    for(let textarea of textareas){
        if (textarea.hasAttribute("aqua_champ") && (!textarea.hasAttribute("aqua_prof") || textarea.getAttribute("aqua_prof") == String(profondeur))){
            resultat[textarea.getAttribute("aqua_champ")] = textarea.value;
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
            resultat[nom_liste].push(donnees_conf_analyse(div, profondeur+1));
        }
    }
    return resultat
}

/* ANCIENNES FONCTIONS */

function avt_choix_enregistrement(){
    parent.window.go.main.App.CreationNouveauProjet().then(resultat =>{
        document.getElementById("enregistrement").value = resultat;
        document.getElementById("archives").value = "";
        document.getElementById("valider").style.display = "none";
    })
}

function avt_choix_orc(){
    let chemin_enreg = document.getElementById("enregistrement").value;
    if(chemin_enreg == ""){
        alert("Vous deviez d'abord choisir où vous voulez enregistrer votre ORC")
    } 
    else{
        document.getElementById("patientez").style.display = "inline";
        document.getElementById("formulaire").style.display = "none";
        parent.window.go.main.App.AjoutORCNouveauProjet().then(resultat =>{
            document.getElementById("archives").value = resultat;
            if(resultat == ""){
                document.getElementById("valider").style.display = "none";
            }
            else{
                let nom_auteur = document.getElementById("auteur").value;
                if(nom_auteur != ""){
                    document.getElementById("valider").style.display = "inline";
                }
            }
            document.getElementById("patientez").style.display = "none";
            document.getElementById("formulaire").style.display = "inline";
        })
    }
}

function avt_change_auteur(){
    let nom_auteur = document.getElementById("auteur").value;
    let enregistrement = document.getElementById("enregistrement").value;
    let archives = document.getElementById("archives").value;
    if(nom_auteur == "" || enregistrement == "" || archives == ""){
        document.getElementById("valider").style.display = "none";
    }
    else{
        document.getElementById("valider").style.display = "inline";
    }
}

function avt_validation(){
    let auteur = document.getElementById("auteur").value;
    let description = document.getElementById("description").value;
    parent.window.go.main.App.ValidationCreationProjet(auteur, description).then(resultat =>{
        if(resultat){ 
            window.location.replace("../html/extraction.html");
            parent.document.getElementsByTagName("header")[0].style.display = "inline";
            let onglet_courant = parent.document.getElementById("onglet_extraction");
            onglet_courant.style.backgroundColor = "#FCF5DC";
            onglet_courant.style.color = "#000";
        }else{
            window.location.replace("../html/accueil.html")
        }
    })
}