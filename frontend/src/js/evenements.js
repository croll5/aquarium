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
 */

let demande_patienter = true;
let chargement_en_cours = false;
let logo_aquarium = document.getElementById("logo_aquarium_patientez")

let filtres = [];

let params = new URLSearchParams(document.location.search);
let machine_a_afficher = params.get("machine");
let nb_evenements_affiches = 0;

const TAILLE_MAX_EVT_SUCCINT = 300;
const NB_EVENEMENTS_PAR_RECUP = 50;

document.getElementById("nom_machine").textContent = params.get("nom_machine");

afficher_nouveau_filtre(params.get("nom_machine"), "machine");

function recherche_indicateur(evt){
    if(evt.key == "Enter"){
        let zone_recherche = document.getElementById("zone_recherche");
        let indicateur = zone_recherche.value;
        afficher_nouveau_filtre(indicateur);
        zone_recherche.value = "";
        if(filtres.length > 0){
            filtrer_evenements_affiches(indicateur);
            filtres[filtres.length] = {"champ":"", "valeur":indicateur, "principal":false};
        }else{
            afficher_recherche_indicateur(indicateur);
            document.getElementById("bouton_selecteur_date").style.display = "none";
            filtres[filtres.length] = {"champ":"", "valeur":indicateur, "principal":true};
        }
    }
}

function filtrer_evenements_affiches(indicateur){
    let liste_evenements = document.getElementById("liste_evenements");
    let nb_suppressions = 0;
    for(let evenement of liste_evenements.childNodes){
        if(!evenement.textContent.includes(indicateur)){
            evenement.style.display = "none";
            nb_suppressions ++;
        }
    }
    if(nb_suppressions > 0){
        afficher_evenements_suivants(nb_suppressions, nb_evenements_affiches)
    }
}

function afficher_recherche_indicateur(indicateur){
    let liste_evenements = document.getElementById("liste_evenements");
    lancer_fonction_patientez()
    liste_evenements.textContent = "";
    parent.window.go.main.App.ChercherIndicateurDansBase(indicateur, NB_EVENEMENTS_PAR_RECUP).then((resultat) =>{
        demande_patienter = false;
        for(let evenement of resultat){
            let ligne_evenement = creer_html_ligne_evenement(evenement)
            liste_evenements.appendChild(ligne_evenement);
        }
        document.getElementById("bouton_evt_precedents").style.display = "none";
        if(resultat.length == NB_EVENEMENTS_PAR_RECUP){
            document.getElementById("bouton_evt_suivants").style.display = "inline";
        }else{
            document.getElementById("bouton_evt_suivants").style.display = "none";
        }
        if(filtres.length != 0){
            for(let filtre of filtres){
                filtrer_evenements_affiches(filtre["valeur"])
            }
        }
        nb_evenements_affiches = NB_EVENEMENTS_PAR_RECUP;
    })
}

function lancer_fonction_patientez(){
    demande_patienter = true;
    if(chargement_en_cours){
        return;
    }
    chargement_en_cours = true;
    let pos_a = "100vw", pos_b = "-30vw";
    let img_a = "../assets/images/dessin_patientez.png"
    let img_b = "../assets/images/dessin_patientez_inverse.png"
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

function afficher_suite(){
    console.log(nb_evenements_affiches)
    afficher_evenements_suivants(NB_EVENEMENTS_PAR_RECUP, nb_evenements_affiches)
}

function afficher_precedent(){
    let liste_evenements = document.getElementById("liste_evenements");
    console.log("avant : ", nb_evenements_affiches-liste_evenements.childElementCount-NB_EVENEMENTS_PAR_RECUP)
    afficher_evenements_suivants(NB_EVENEMENTS_PAR_RECUP, nb_evenements_affiches-liste_evenements.childElementCount-NB_EVENEMENTS_PAR_RECUP, true);
}

function repositionner_date(){
    // Récupération de la date
    lancer_fonction_patientez();
    let date_choisie = new Date(document.getElementById("selecteur_date").value)
    parent.window.go.main.App.PositionDateDansChronologie(params.get("machine"), date_choisie).then(resultat =>{
        document.getElementById("liste_evenements").textContent = "";
        nb_evenements_affiches = resultat;
        document.getElementById("bouton_evt_precedents").style.display = "inline";
        afficher_evenements_suivants(NB_EVENEMENTS_PAR_RECUP, resultat);
    })
}

function evenement_passe_filtres(evenement) {
    console.log(evenement);
    for(let filtre of filtres){
        let contientFiltre = false;
        for(let [_,valeur] of Object.entries(evenement)){
            if( valeur != null && valeur.toString().includes(filtre["valeur"])){
                contientFiltre = true;
                break
            }
        }
        if(!contientFiltre){
            return false
        }
    }
    return true
}

afficher_evenements_suivants(NB_EVENEMENTS_PAR_RECUP, 0);

function afficher_evenements_suivants(nombre, decalage, debut=false){
    console.log("décalage : ", decalage)
    demande_patienter = true;
    setTimeout(()=>{
        if(demande_patienter){
            lancer_fonction_patientez();
        }
    },1000)
    if(filtres.length > 0){
        parent.window.go.main.App.ChercherIndicateurEvenementSuivant(decalage, decalage+nombre).then(resultat =>{
            let liste_evenements = []
            for(let evenement of resultat){
                if(evenement_passe_filtres(evenement)){
                    liste_evenements[liste_evenements.length] = evenement
                }
            }
            afficher_evenements(liste_evenements, liste_evenements.length, debut)
            if(resultat.length > 0 && liste_evenements.length != nombre){
                nb_evenements_affiches = decalage+nombre;
                afficher_evenements_suivants(nombre-liste_evenements.length, decalage+nombre)
            } else if(resultat.length == 0){
                document.getElementById("bouton_evt_suivants").style.display = "none";
            }
        })
    }else{
        parent.window.go.main.App.ContenuEvenementsChronologie(params.get("machine"), decalage, nombre).then(resultat =>{
            afficher_evenements(resultat, nombre, debut)
        })
    }
}

function afficher_evenements(resultat, nombre, debut){
    let nombre_reel = resultat.length
    if(nombre_reel < nombre){
        document.getElementById("bouton_evt_suivants").style.display = "none";
    } else {
        document.getElementById("bouton_evt_suivants").style.display = "inline";
    }
    demande_patienter = false;
    let premier_element;
    let liste_evenements = document.getElementById("liste_evenements");
    if(debut){
        premier_element = liste_evenements.firstChild;
    }
    for(let evenement of resultat){
        let ligne_evenement = creer_html_ligne_evenement(evenement)
        if(debut && premier_element != undefined){
            premier_element.before(ligne_evenement)
        }else{
            liste_evenements.appendChild(ligne_evenement);
        }
    }
    if(!debut){
        nb_evenements_affiches += nombre_reel
    }else if(liste_evenements.childElementCount > NB_EVENEMENTS_PAR_RECUP*5){
        nb_evenements_affiches -= nombre_reel;
    }
    if(!debut && liste_evenements.childElementCount > NB_EVENEMENTS_PAR_RECUP*5){
        document.getElementById("bouton_evt_precedents").style.display = "inline";
        for(let i = 0; i < nombre_reel; i++){
            liste_evenements.firstChild.remove();
        }
    } else if(debut && liste_evenements.childElementCount > NB_EVENEMENTS_PAR_RECUP*5){
        document.getElementById("bouton_evt_suivants").style.display = "inline";
        for(let i = 0; i < nombre_reel; i++){
            liste_evenements.lastChild.remove()
        }
        if(decalage == 0){
            document.getElementById("bouton_evt_precedents").style.display = "none";
        }
    }
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


function afficher_nouveau_filtre(valeur_filtre, nom_champ=""){
    let zone_recherche = document.getElementById("zone_recherche")
    let filtre = document.createElement("div");
    if(nom_champ == ""){
        filtre.textContent = valeur_filtre;
    } else{
        filtre.textContent = nom_champ + ": " + valeur_filtre;
    }
    
    filtre.classList.add("filtre_recherche");
    let bouton_fermer = document.createElement("button");
    bouton_fermer.textContent = "✖";
    bouton_fermer.onclick = function (event) {
        filtre.remove();
        suppression_filtre(nom_champ, valeur_filtre)
    }
    filtre.appendChild(bouton_fermer);
    zone_recherche.before(filtre);
}

function suppression_filtre(nom_champ, valeur_filtre){
    console.log(filtres)
    let liste_evenements = document.getElementById("liste_evenements").children;
    let numFiltre = -1;
    for(let i = 0; i < filtres.length; i++){
        if(filtres[i]["valeur"] == valeur_filtre){
            numFiltre = i
        }
    }
    if(numFiltre == -1){
        return
    }
    if(filtres[numFiltre]["principal"] && filtres.length > 1){
        let nouveauNumFiltre = numFiltre+1 % filtres.length
        filtres[nouveauNumFiltre]["principal"] = nouveauNumFiltre
        afficher_recherche_indicateur(filtres[nouveauNumFiltre]["valeur"])
        filtres.splice(numFiltre, 1)
    }
    else if(filtres.length < 2){
        filtres = []
        document.getElementById("bouton_selecteur_date").style.display = "inline";
        let premier_evenement = liste_evenements[0]
        if(premier_evenement != undefined){
            let heure = premier_evenement.getElementsByTagName("h3")[0].textContent
            console.log(heure)
            if(heure != undefined){
                let date_debut = new Date(heure)
                document.getElementById("selecteur_date").value = date_debut.toISOString().slice(0,16);
            }
        }
        repositionner_date();
    }else{
        document.getElementById("liste_evenements").textContent = "";
        for(let i = 0; i < filtres.length; i++){
            if(filtres[i]["valeur"] == valeur_filtre){
                filtres.splice(i, 1)
                break;
            }
        }
        afficher_evenements_suivants(NB_EVENEMENTS_PAR_RECUP, 0)
    }
}

function suppression_filtre_secondaire(nom_champ, valeur_filtre){
    let liste_evenements = document.getElementById("liste_evenements").children;
    for(let i = 0; i < filtres.length; i++){
        if(filtres[i]["valeur"] == valeur_filtre){
            filtres.splice(i, 1)
            break;
        }
    }
    for(let evenement of liste_evenements){
        let contientTout = true;
        for(let filtre of filtres){
            if(! evenement.textContent.includes(filtre["champ"]+filtre["valeur"])){
                contientTout = false;
            }
        }
        if(contientTout){
            evenement.style.display = "block";
        }else{
            evenement.style.display = "none";
        }
    }
}