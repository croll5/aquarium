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


let params = new URLSearchParams(document.location.search);
let machine_a_afficher = params.get("machine");

document.getElementById("nom_machine").textContent = params.get("nom_machine");

if(machine_a_afficher == null || machine_a_afficher == ""){
    parent.signaler_erreur("", "Aucune machine à afficher...")
}else{
    affichage_si_extraction_en_cours(machine_a_afficher);
}



function affichage_si_extraction_en_cours(idMachine){
    parent.window.go.main.App.ExtractionEnCours().then(resultat =>{
        if(resultat){
            document.getElementById("selection_arborescence").style.display = "none";
            let patientez = document.getElementById("patientez");
            afficher_salle_d_attente(patientez);
            let verifPasExtrait = setInterval(function(){
                parent.window.go.main.App.ExtractionEnCours().then(reponse =>{
                    if(!reponse){
                        clearInterval(verifPasExtrait);
                        extraire_arborescence(idMachine);
                    }
                });
            }, 1000);
        } else{
            extraire_arborescence(idMachine);
        }
    })
}

async function extraire_arborescence(idMachine){
    // On affiche au besoin la salle d’attente
    let patientez = document.getElementById("patientez");
    parent.window.go.main.App.ArborescenceEnCache().then(resultat =>{
        if(resultat != idMachine){
            console.log("oki");
            afficher_salle_d_attente(patientez);
        }
    });
    // On supprime tout ce qu’il y a dans l’arborescence actuelle
    let div_arborescence = document.getElementById("arborescence");
    div_arborescence.innerHTML = "";
    // On affiche la racine de l’arborescence
    await ajouter_contenu_dossier(div_arborescence, [], idMachine);
    patientez.style.display = "none";
}

async function ajouter_contenu_dossier(emplacement, chemin, idMachine){
    return new Promise(fini =>{
        let copieEmplacement = emplacement.cloneNode(true);
        emplacement.replaceWith(copieEmplacement);
        emplacement = copieEmplacement;
        emplacement.open = true;
        parent.window.go.main.App.ArborescenceMachineAnalysee(chemin,idMachine).then(resultat => {
            console.log(idMachine);
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
                    ajouter_infos_fichier(affichage_fichier, fichier, idMachine);
                    emplacement.appendChild(affichage_fichier);
                }
            }
            fini();
        })
    })
    
}

function ajouter_infos_fichier(emplacement, donnees_fichier, idMachine) {
    let details = document.createElement("button");
    details.textContent = "🪪";
    details.classList.add("bouton_invisible");
    details.addEventListener("click", function(event){afficher_metadonnees_fichier(donnees_fichier["IdSource"], idMachine)});
    emplacement.appendChild(details);
    if(donnees_fichier["IdCopie"] != 0){
        let visualiser = document.createElement("button");
        visualiser.textContent = "🔍";
        visualiser.classList.add("bouton_invisible");
        visualiser.addEventListener("click", function(event){alert("⌛ La fonctionalité de visualisation de fichiers n’a pas encore été implémentée. \nVous pouvez télécharger la mise à jour sur https://github.com/croll5/aquarium. Peut-être y sera-t-elle implémentée ? 🙃")});
        emplacement.appendChild(visualiser);
    }
}

function afficher_metadonnees_fichier(idFichier, idMachine) {
    parent.window.go.main.App.DetailsFichierArborescence(idFichier, idMachine).then(resultat => {
        console.log(resultat);
        if(resultat.length < 1){
            return
        }
        // On remplit la popup avec les informations sur le fichier
        let table = document.getElementById("table_metadonnees")
        table.innerHTML = "";
        for(let [nom_donnee, donnee] of Object.entries(resultat[0])){
            let ligneInfo = document.createElement("tr");
            let nom = document.createElement("td");
            nom.textContent = nom_donnee;
            ligneInfo.appendChild(nom);
            let valeur = document.createElement("td");
            valeur.textContent = donnee;
            ligneInfo.appendChild(valeur);
            table.appendChild(ligneInfo);
            console.log(nom_donnee + " : " + donnee);
        }
        document.getElementById("fond_popup").style.display = "block"
        document.getElementById("popup_infos_fichier").style.display = "block";
    })
}