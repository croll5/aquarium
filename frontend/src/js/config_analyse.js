parent.window.go.main.App.ListeExtractionsPossibles().then(resultat => {
    console.log(resultat);
    let listeExtractions = document.getElementById("liste_extractions");
    // Création d’un affichage d’extraction
    for(let donneesExtraction of resultat){
        let extraction = document.createElement("div");
        extraction.draggable = true;
        extraction.classList.add("extraction", "bouton_mi_tons", "barre_progression");
        listeExtractions.appendChild(extraction);
        // Ajout de la barre de progression
        let progression = document.createElement("div");
        progression.classList.add("pct_chargement");
        progression.style.width = "0%";
        extraction.appendChild(progression); 
        // Ajout du texte
        let texte = document.createElement("div");
        texte.classList.add("texte_extractions");
        texte.textContent = donneesExtraction["NomMachine"] + " - " + donneesExtraction["InfosExtraction"]["Nom"];
        extraction.appendChild(texte);
        // Ajout de la description
        let description = document.createElement("i");
        description.classList.add("texte_extractions", "description_extraction");
        description.textContent = donneesExtraction["InfosExtraction"]["Description"];
        extraction.appendChild(description)
    }
    
})

let elementSelectionne = null;

function debut_deplacement(e){
    elementSelectionne = e.target;
    e.target.classList.add('element_deplace');
}

function fin_deplacement(e){
    e.target.classList.remove('element_deplace');
    document.querySelectorAll('.extraction')
        .forEach(item => item.classList.remove('over'));
    elementSelectionne = null;
}

function deplacement(e){
    e.preventDefault();
    let contenantElements = document.getElementById("liste_extractions");
    let listePositionsPotentielles = document.getElementsByClassName("extraction");
    let position_element_dessous = window.screen.height;
    let element_dessous;
    for(let element of listePositionsPotentielles){
        let positionElement = element.getBoundingClientRect().y;
        if(positionElement > e.y && positionElement < position_element_dessous){
            element_dessous = element;
            position_element_dessous = positionElement;
        }
    }
    console.log(element_dessous);
    if(element_dessous != null && !element_dessous.classList.contains("element_deplace")){
        element_dessous.before(elementSelectionne);
    } else if(element_dessous == null){
        contenantElements.appendChild(elementSelectionne);
    }
}

function lancer_extraction(){
    document.getElementById("lancer_extraction").style.display = "none";
    let liste_extractions = document.getElementsByClassName("extraction")
    for(let extraction of liste_extractions){
        extraction.draggable = false;
    }
}