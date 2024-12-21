parent.window.go.main.App.ListeExtractionsPossibles().then(resultat =>{
    console.log(resultat);
    let div_possibilites = document.getElementById("possibilites_extractions");
    // Pour chaque extracteur possible, on l'ajoute sur la page
    for (let [cle, valeur] of Object.entries(resultat)) {
        let paragraphe = document.createElement('p');
        paragraphe.innerText = "🫧" + valeur["Description"] + "🫧";
        if(valeur["Progression"] >= 100){
            paragraphe.innerText += " ✅";
            paragraphe.className = "non_cliquable"
        }else if(valeur["Progression"] >= 0){
            ajouter_chargement(paragraphe, valeur["Progression"], cle)
        }else{
            
            paragraphe.className = "liste_options"
            paragraphe.onclick = function() { extraire_elements(cle) };
        }
        paragraphe.id = cle;
        
        div_possibilites.appendChild(paragraphe);
    }
})

function extraire_elements(module_id){
    let paragraphe = document.getElementById(module_id);
    parent.window.go.main.App.ExtraireElements(module_id, paragraphe.value);
    paragraphe.onclick = "";
    ajouter_chargement(paragraphe, 0, module_id)
}

function ajouter_chargement(paragraphe, valeur_initiale, module_id){
    let progression = document.createElement("progress");
    progression.max = 100;
    progression.value = valeur_initiale;
    paragraphe.textContent += " - chargement... ";
    paragraphe.appendChild(progression);
    let annuler = document.createElement("button");
    annuler.innerText = "❌";
    annuler.className = "bouton_invisible";
    annuler.onclick = function() { annuler_extraction(module_id) };
    let maj = setInterval(function(){
        parent.window.go.main.App.ProgressionExtraction(module_id).then(pourcentageExtraction =>{
        progression.value = pourcentageExtraction;
        if (progression.value >= 100){
            paragraphe.removeChild(progression);
            paragraphe.removeChild(annuler);
            paragraphe.textContent = paragraphe.textContent.replace("- chargement... ", "✅");
            clearInterval(maj);
            paragraphe.className = "non_cliquable";
        }
        })
    },50);
    paragraphe.appendChild(annuler);
}

function annuler_extraction(idExtracteur){
    if(confirm("Voulez-vous vraiment annuler l'extraction de " + idExtracteur + " ?")){
        parent.window.go.main.App.AnnulerExtraction(idExtracteur).then(succes =>{
            if(succes) {
                alert("L'extraction a bien été annulée 🥲");
                let paragraphe = document.getElementById(idExtracteur);
                paragraphe.onclick = function() { extraire_elements(idExtracteur) };
                let enfant = paragraphe.lastElementChild;
                while (enfant) {
                    paragraphe.removeChild(enfant);
                    enfant = paragraphe.lastElementChild;
                }
                paragraphe.innerHTML = paragraphe.innerText.replace("- chargement...", "");
            }else{
                alert("L'extraction n'a pas pu s'arrêter correctement. Réessayez.")
            }
        })
    }
}