parent.window.go.main.App.ListeExtractionsPossibles().then(resultat =>{
    let div_possibilites = document.getElementById("possibilites_extractions");
    // Pour chaque extracteur possible, on l'ajoute sur la page
    for (let [key, value] of Object.entries(resultat)) {
        let paragraphe = document.createElement('p');
        paragraphe.innerText = "🫧" + value + "🫧";
        paragraphe.id = key;
        paragraphe.className = "liste_options"
        paragraphe.onclick = function() { extraire_elements(key) };
        div_possibilites.appendChild(paragraphe);
    }
})

function extraire_elements(module_id){
    let description = document.getElementById(module_id).value;
    parent.window.go.main.App.ExtraireElements(module_id, description);
}
