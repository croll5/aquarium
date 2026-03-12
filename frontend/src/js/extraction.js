/*
Copyright ou © ou Copr. Cécile Rolland , (21 janvier 2025) 

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

parent.window.go.main.App.ListeExtractionsPossibles().then(resultat =>{
    console.log(resultat);
    let div_possibilites = document.getElementById("possibilites_extractions");
    // Pour chaque extracteur possible, on l'ajoute sur la page
    for (let [cle, valeur] of Object.entries(resultat)){
        let details_machine = document.createElement("details");
        let nom_machine = document.createElement("summary");
        nom_machine.textContent = valeur["ConfigMachine"]["nom"];
        div_possibilites.appendChild(details_machine);
        details_machine.appendChild(nom_machine);
        ajouter_extractions_machine(details_machine, valeur["ListeExtractions"], cle);
    }
})

function ajouter_extractions_machine(div_machine, resultat, idMachine){
    for (let [cle, valeur] of Object.entries(resultat)) {
        let paragraphe = document.createElement('p');
        paragraphe.innerText = "🫧" + valeur["Description"] + "🫧";
        if(valeur["Progression"] >= 100){
            paragraphe.innerText += " ✅";
            paragraphe.className = "non_cliquable";
        }else if(valeur["Progression"] >= 0){
            ajouter_chargement(paragraphe, valeur["Progression"], cle, idMachine);
        }else{
            paragraphe.className = "liste_options";
            paragraphe.onclick = function() { extraire_elements(cle, idMachine) };
        }
        paragraphe.id = cle;
        div_machine.appendChild(paragraphe);
    }
}

function extraire_elements(module_id, id_machine){
    let paragraphe = document.getElementById(module_id);
    parent.window.go.main.App.ExtraireElements(module_id, paragraphe.value, id_machine);
    paragraphe.onclick = "";
    ajouter_chargement(paragraphe, 0, module_id, id_machine);
}

function ajouter_chargement(paragraphe, valeur_initiale, module_id, id_machine){
    let progression = document.createElement("progress");
    progression.max = 100;
    progression.value = valeur_initiale;
    paragraphe.textContent += " - chargement... ";
    paragraphe.appendChild(progression);
    let annuler = document.createElement("button");
    annuler.innerText = "❌";
    annuler.className = "bouton_invisible";
    annuler.onclick = function() { annuler_extraction(module_id, id_machine) };
    let maj = setInterval(function(){
        parent.window.go.main.App.ProgressionExtraction(id_machine, module_id).then(pourcentageExtraction =>{
            progression.value = pourcentageExtraction;
            console.log(pourcentageExtraction)
            if (progression.value >= 100){
                progression.remove();
                annuler.remove();
                paragraphe.textContent = paragraphe.textContent.replace("- chargement... ", "✅");
                clearInterval(maj);
                paragraphe.className = "non_cliquable";
            }
        })
    },50);
    paragraphe.appendChild(annuler);
}

function annuler_extraction(idExtracteur, idMachine){
    if(confirm("Voulez-vous vraiment annuler l'extraction de " + idExtracteur + " ?")){
        parent.window.go.main.App.AnnulerExtraction(idMachine, idExtracteur).then(succes =>{
            if(succes) {
                alert("L'extraction a bien été annulée 🥲");
                let paragraphe = document.getElementById(idExtracteur);
                paragraphe.onclick = function() { extraire_elements(idExtracteur) };
                let enfant = paragraphe.lastElementChild;
                while (enfant) {
                    enfant.remove();
                    enfant = paragraphe.lastElementChild;
                }
                paragraphe.innerHTML = paragraphe.innerText.replace("- chargement...", "");
            }else{
                alert("L'extraction n'a pas pu s'arrêter correctement. Réessayez.")
            }
        })
    }
}

function extraire_chronologie(){
    document.body.style.cursor = "wait";
    document.getElementById("bouton_ext_chrono").style.display = "none";
    document.getElementById("possibilites_extractions").style.display = "none";
    document.getElementById("patience_chronologie").style.display = "inline";
    document.getElementById("document_patientez").src = "../assets/documents/MONTAUBAN_Albane_Memoire_M2.pdf";
    parent.window.go.main.App.ExtractionChronologie().then(resultat =>{
        alert("L'extraction de la chronologie s'est terminée avec succès 🥳");
        document.getElementById("bouton_ext_chrono").style.display = "inline";
        document.getElementById("possibilites_extractions").style.display = "inline";
        document.getElementById("patience_chronologie").style.display = "none";
        document.body.style.cursor = "default";
    })
}