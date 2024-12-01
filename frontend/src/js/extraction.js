
// Utilisation de la fonction avec un callback
window.onload = function() {
    executeWhenReady(function() {
        parent.window.go.main.App.ListeExtractionsPossibles().then(resultat => {
            let div_possibilites = document.getElementById("possibilites_extractions");
            checkBox = "NO - "
            // Pour chaque extracteur possible, on l'ajoute sur la page
            for (let [key, value] of Object.entries(resultat)) {
                let paragraphe = document.createElement('p');
                paragraphe.innerText = "" + checkBox + value + "";
                paragraphe.id = key;
                paragraphe.className = "liste_options";
                paragraphe.onclick = function() { extraire_elements(key) };
                div_possibilites.appendChild(paragraphe);
            }
        }).catch(err => console.error("Error fetching extractions:", err));
    });
};



function extraire_elements(module_id){
    let description = document.getElementById(module_id).value;
    parent.window.go.main.App.ExtraireElements(module_id, description);
}


