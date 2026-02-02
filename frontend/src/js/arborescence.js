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

remplir_select_machines();

function remplir_select_machines(){
    let select_machine = document.getElementById("choix_arborescence")
    parent.window.go.main.App.ListeMachinesAnalysees().then(resultat =>{
        for(let [idMachine, configMachine] of Object.entries(resultat)){
            let nvelle_option = document.createElement("option");
            nvelle_option.value = idMachine;
            nvelle_option.textContent = configMachine["nom"];
            select_machine.appendChild(nvelle_option);
        }
    })
}

function extraire_arborescence(){
    // On commence par regarder quelle arborescence il faut afficher
    let idMachine = document.getElementById("choix_arborescence").value;
    // On supprime tout ce qu’il y a dans l’arborescence actuelle
    let div_arborescence = document.getElementById("arborescence");
    div_arborescence.innerHTML = "";
    // On affiche la racine de l’arborescence
    ajouter_contenu_dossier(div_arborescence, [], idMachine).then(() => {
        document.getElementById("document_pour_patienter").style.display = "none";
    })
}

async function ajouter_contenu_dossier(emplacement, chemin, idMachine){
    let copieEmplacement = emplacement.cloneNode(true);
    emplacement.replaceWith(copieEmplacement);
    emplacement = copieEmplacement;
    emplacement.open = true;
    parent.window.go.main.App.ArborescenceMachineAnalysee(chemin,idMachine).then(resultat => {
        for(let fichier of resultat){
            if(fichier["ADesEnfants"]){
                let dossier = document.createElement("details");
                dossier.classList.add("dossier_arborescence");
                let nomDossier = document.createElement("summary");
                nomDossier.textContent = fichier["Nom"];
                dossier.appendChild(nomDossier);
                const cheminFichier = chemin.concat([fichier["Nom"]]);
                dossier.addEventListener("click", function(ev) {ajouter_contenu_dossier(dossier, cheminFichier, idMachine)});
                emplacement.appendChild(dossier);
            }else{
                let affichage_fichier = document.createElement("p");
                affichage_fichier.textContent = fichier["Nom"];
                affichage_fichier.classList.add("fichier_arborescence");
                ajouter_infos_fichier(affichage_fichier, fichier);
                emplacement.appendChild(affichage_fichier);
            }
        }
    })
}

function ajouter_infos_fichier(emplacement, donnees_fichier) {
    let details = document.createElement("button");
    details.textContent = "🪪";
    details.classList.add("bouton_invisible");
    details.addEventListener("click", function(event){alert("Métadonnées de " + donnees_fichier["IdSource"])});
    emplacement.appendChild(details);
    if(donnees_fichier["IdCopie"] != 0){
        let visualiser = document.createElement("button");
        visualiser.textContent = "🔍";
        visualiser.classList.add("bouton_invisible");
        visualiser.addEventListener("click", function(event){alert("Visualisation de " + donnees_fichier["IdCopie"])});
        emplacement.appendChild(visualiser);
    }
}