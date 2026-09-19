/**
Copyright Cécile Rolland, (6 août 2026) 

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
-->
 */

let demande_patienter = true;
let chargement_en_cours = false;
let logo_aquarium = document.getElementById("logo_aquarium_patientez")

function lancer_fonction_patientez(){
    demande_patienter = true;
    if(chargement_en_cours){
        return;
    }
    let scroll = document.getElementById("body_arbo").scrollTop;
    console.log(scroll)
    chargement_en_cours = true;
    let pos_a = "100vw", pos_b = "-30vw";
    let img_a = "../assets/images/dessin_patientez.png"
    let img_b = "../assets/images/dessin_patientez_inverse.png"
    console.log(logo_aquarium.style.left)
    if(logo_aquarium.style.left == pos_a){
        chemin_patientez(pos_b, pos_a, img_b, img_a)
    }else{
        chemin_patientez(pos_a, pos_b, img_a, img_b);
    }
}

function chemin_patientez(pos_a, pos_b, img_a, img_b){
    if(!demande_patienter){
        chargement_en_cours = false;
        return
    }
    logo_aquarium.style.left = pos_a;
    logo_aquarium.src = img_a
    setTimeout(() => {
        chemin_patientez(pos_b, pos_a, img_b, img_a)
    }, 5000);
}

let params = new URLSearchParams(document.location.search);
let machine_a_afficher = params.get("machine");
let nb_evenements_affiches = 0;

const TAILLE_MAX_EVT_SUCCINT = 300;

document.getElementById("nom_machine").textContent = params.get("nom_machine");

afficher_nouveau_filtre("machine", params.get("nom_machine"));

function afficher_suite(){
    afficher_evenements_suivants(50, nb_evenements_affiches)
}

function afficher_precedent(){
    let liste_evenements = document.getElementById("liste_evenements");
    afficher_evenements_suivants(50, nb_evenements_affiches-liste_evenements.childElementCount-50, true);
}

function repositionner_date(){
    // Récupération de la date
    lancer_fonction_patientez();
    let date_choisie = new Date(document.getElementById("selecteur_date").value)
    parent.window.go.main.App.PositionDateDansChronologie(params.get("nachine"), date_choisie).then(resultat =>{
        document.getElementById("liste_evenements").textContent = "";
        nb_evenements_affiches = resultat;
        document.getElementById("bouton_evt_precedents").style.display = "inline";
        afficher_evenements_suivants(50, resultat);
    })
}

afficher_evenements_suivants(50, 0);

function afficher_evenements_suivants(nombre, decalage, debut=false){
    demande_patienter = true;
    setTimeout(()=>{
        if(demande_patienter){
            lancer_fonction_patientez();
        }
    },1000)
    parent.window.go.main.App.ContenuEvenementsChronologie(params.get("machine"), decalage, nombre).then(resultat =>{
        demande_patienter = false;
        let premier_element;
        let liste_evenements = document.getElementById("liste_evenements");
        if(debut){
            premier_element = liste_evenements.firstChild;
        }
        for(let evenement of resultat){
            let ligne_evenement = creer_html_ligne_evenement(evenement)
            if(debut){
                premier_element.before(ligne_evenement)
            }else{
                liste_evenements.appendChild(ligne_evenement);
            }
        }
        if(debut){
            nb_evenements_affiches -= nombre
        }else{
            nb_evenements_affiches += nombre;
        }
        if(!debut && liste_evenements.childElementCount > nombre*5){
            document.getElementById("bouton_evt_precedents").style.display = "inline";
            for(let i = 0; i < nombre; i++){
                liste_evenements.firstChild.remove()
            }
        } else if(debut && liste_evenements.childElementCount > nombre*5){
            for(let i = 0; i < nombre; i++){
                liste_evenements.lastChild.remove()
            }
            if(decalage == 0){
                document.getElementById("bouton_evt_precedents").style.display = "none";
            }
        }
    })
}

function creer_html_ligne_evenement(evenement){
    let ligne_evenement = document.createElement("li");
    let horodatage = document.createElement("h3");
    horodatage.textContent = evenement.aqua_horodatage;
    ligne_evenement.appendChild(horodatage);
    let contenu_detaille = document.createElement("div");
    contenu_detaille.classList.add("details_evenement");
    let bouton_fermer = document.createElement("button");
    bouton_fermer.textContent = "×";
    bouton_fermer.classList.add("bouton_fermer");
    bouton_fermer.onclick = () => {
        contenu_detaille.style.display = "none";
        contenu_succint.style.display = "block";
    }
    contenu_detaille.appendChild(bouton_fermer);
    let contenu_succint = document.createElement("p");
    contenu_succint.classList.add("evenement_tronque");
    // Permettre d’afficher le détail d’un évènement
    contenu_succint.onclick = () => {
        contenu_succint.style.display = "none";
        contenu_detaille.style.display = "block";
    }
    let complet = false;
    for(let [attribut, valeur] of Object.entries(evenement)){
        if(attribut == "aqua_horodatage"){
            continue;
        }
        if(valeur == ""){
            continue;
        }
        // Ajout au contenu succint de l’évènement
        let champ_detaille = document.createElement("h4");
        champ_detaille.textContent = attribut;
        contenu_detaille.appendChild(champ_detaille);
        let detail_valeur = document.createElement("p");
        detail_valeur.textContent = valeur
        contenu_detaille.appendChild(detail_valeur)
        // Ajout au contenu détaillé de l’évènement
        if(complet){
            continue
        }
        let attribut_html = document.createElement("strong");
        attribut_html.textContent = (attribut + " : ").slice(0,Math.max(0,TAILLE_MAX_EVT_SUCCINT-contenu_succint.textContent.length));
        contenu_succint.appendChild(attribut_html);
        let valeur_html = document.createElement("span");
        valeur_html.textContent = (valeur+" ").slice(0,Math.max(0,TAILLE_MAX_EVT_SUCCINT-contenu_succint.textContent.length));
        contenu_succint.appendChild(valeur_html);
        if(contenu_succint.textContent.length >= TAILLE_MAX_EVT_SUCCINT){
            let suspension = document.createElement("span");
            suspension.textContent = "…";
            contenu_succint.appendChild(suspension);
            complet = true;
        }
    }
    ligne_evenement.appendChild(contenu_detaille);
    ligne_evenement.appendChild(contenu_succint);
    return ligne_evenement
}


function afficher_nouveau_filtre(nom_champ, valeur_filtre){
    let zone_recherche = document.getElementById("zone_recherche")
    let liste_filtres = document.getElementById("barre_recherche");
    let filtre = document.createElement("div");
    filtre.textContent = nom_champ + ": " + valeur_filtre;
    filtre.classList.add("filtre_recherche");
    let bouton_fermer = document.createElement("button");
    bouton_fermer.textContent = "✖";
    bouton_fermer.onclick = function (event) {
        filtre.remove();
    }
    filtre.appendChild(bouton_fermer);
    zone_recherche.before(filtre);
}