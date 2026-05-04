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

let version = "1.0"

function change_onglet(destination, onglet){
    document.getElementsByTagName("iframe")[0].src = destination;
    // On met tous les autres onglets à la couleur standard
    let onglets = document.getElementsByClassName("onglet");
    for (const onglet of onglets) {
        if(onglet.classList.contains("onglet_selectionne")){
            onglet.classList.remove("onglet_selectionne");
        }
    }
    // On met l'onglet sur lequel on va à la couleur de la page
    let onglet_courant = onglet.parentNode;
    onglet_courant.classList.add("onglet_selectionne");
}

function accueil(){
    document.getElementsByTagName("iframe")[0].src = "html/accueil.html";
    document.getElementsByTagName("header")[0].style.display = "none";
}

function creer_onglet(url, nomOnglet){
    // On crée le nouvel onglet
    let header = document.getElementsByTagName('header')[0];
    let nvel_onglet = document.createElement("div");
    nvel_onglet.classList.add("onglet");
    header.appendChild(nvel_onglet)
    // On ajoute le bouton (texte de l’onglet)
    let bouton_onglet = document.createElement("button");
    bouton_onglet.onclick = function (ev) {
        change_onglet(url, bouton_onglet);
    }
    bouton_onglet.textContent = nomOnglet;
    nvel_onglet.appendChild(bouton_onglet);
    // On ajoute le bouton pour fermer
    let bouton_fermer = document.createElement("button");
    bouton_fermer.classList.add("fermer_onglet");
    bouton_fermer.textContent = "✖";
    bouton_fermer.onclick = function (event) {
        if(nvel_onglet.classList.contains("onglet_selectionne")){
            nvel_onglet.previousElementSibling.getElementsByTagName("button")[0].click()
        }
        nvel_onglet.remove();
    }
    nvel_onglet.appendChild(bouton_fermer);
    change_onglet(url, bouton_onglet);
}

window.contrastes = false;
window.dyslexie = false;
window.non_aux_bubulles = false;
window.oui_au_debug = false;

function initialiser_parametres(tentatives_restantes = 20){
    const app = window?.go?.main?.App;
    if(!(app && typeof app.GetParametres === "function")){
        if(tentatives_restantes > 0){
            setTimeout(() => initialiser_parametres(tentatives_restantes - 1), 100);
        }
        return;
    }
    Promise.resolve(app.GetParametres()).then(parametres => {
        if(parametres == null){
            return;
        }
        window.contrastes = parametres.Contrastes === true;
        window.dyslexie = parametres.Dyslexie === true;
        window.non_aux_bubulles = parametres.NonAuxBubulles === true;
        window.oui_au_debug = parametres.OuiAuDebug === true;
        const iframe = document.getElementsByTagName("iframe")[0];
        if(iframe && iframe.getAttribute("src")){
            iframe.setAttribute("src", iframe.getAttribute("src"));
        }
    }).catch(() => {});
}

initialiser_parametres();

function signaler_erreur(fichierErreur, detailsErreur){
    document.getElementById("texte_details_erreur").textContent = decodeURIComponent(detailsErreur).replaceAll("+"," ")
    document.getElementById("nom_fichier_erreur").textContent = fichierErreur;
    document.getElementById("bandeau_erreur").style.display = "flex";
    lien_courriel = document.getElementById("courriel_erreur")
    lien_courriel.href = "mailto:aquarium@mailo.com";
    lien_courriel.href += "?Subject=" + encodeURI("[AQUA#" + Math.round(Math.random()*100000) + "] - Signalement de bogue");
    lien_courriel.href += "&body=" + encodeURI("Bonjour,\n\nJe me permets de vous signaler un bogue dans le logiciel Aquarium.\nVeuillez en trouver ci-dessous les détails :\n\n" + decodeURIComponent(detailsErreur).replaceAll("+"," ") + "\n\nPour votre bonne information, j'utilise la version " + version + " d’aquarium.\n\nRespectueusement,\n\n[Prénom] [Nom]")
}

function fermer_bandeau(){
    document.getElementById("bandeau_erreur").style.display = "none";
}

function details_erreur(){
    document.getElementById("fond_popup").style.display = "block";
    document.getElementById("popup_erreur").style.display = "block";
}
