parent.window.go.main.App.ListeExtractionsPossibles().then(resultat => {
    console.log(resultat);
    let listeExtractions = document.getElementById("liste_extractions");
    // Création d’un affichage d’extraction
    for(let [idMachine, configMachine] of Object.entries(resultat)){
        for(let [idExtraction, infosExtraction] of Object.entries(configMachine["ListeExtractions"])){
            console.log(idExtraction);
            console.log(infosExtraction);
            let extraction = document.createElement("div");
            extraction.draggable = true;
            extraction.setAttribute("idExtraction", idExtraction);
            extraction.setAttribute("idMachine", idMachine)
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
            texte.textContent = configMachine["ConfigMachine"]["nom"] + " - " + infosExtraction["InfosExtraction"]["Nom"];
            extraction.appendChild(texte);
            // Ajout de la description
            let description = document.createElement("i");
            description.classList.add("texte_extractions", "description_extraction");
            description.textContent = infosExtraction["InfosExtraction"]["Description"];
            extraction.appendChild(description)
        }
    }
    parent.window.go.main.App.ProgressionExtraction().then(extractionEnCours =>{
        if(extractionEnCours["chargement"] != undefined){
            document.getElementById("lancer_extraction").style.display = "none";
            miseAJourChargement();
        }
    });
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
    if(element_dessous != null && !element_dessous.classList.contains("element_deplace")){
        element_dessous.before(elementSelectionne);
    } else if(element_dessous == null){
        contenantElements.appendChild(elementSelectionne);
    }
}

function lancer_extraction(){
    document.getElementById("lancer_extraction").style.display = "none";
    let liste_extractions = document.getElementsByClassName("extraction");
    let ordreExtraction = [];
    for(let extraction of liste_extractions){
        extraction.draggable = false;
        let extractionSuivante = {
            idMachine : extraction.getAttribute("idMachine"), 
            idExtraction: extraction.getAttribute("idExtraction")
        }
        ordreExtraction.push(extractionSuivante);
    }
    parent.window.go.main.App.LancerExtraction(ordreExtraction);
    miseAJourChargement()
    console.log(ordreExtraction);
}

function miseAJourChargement(){
    let liste_extractions = document.getElementsByClassName("extraction");
    let maj = setInterval(function(){
        parent.window.go.main.App.ProgressionExtraction().then(extractionEnCours =>{
            if(liste_extractions.length == 0){
                console.log("FIN");
                clearInterval(maj);
            }
            for(let i = 0; i < liste_extractions.length; i++){
                let extraction = liste_extractions[i]
                if(extractionEnCours["idExtraction"] == extraction.getAttribute("idExtraction") && extraction.getAttribute("idMachine") == extractionEnCours["idMachine"]){
                    if(extractionEnCours["chargement"] != undefined){
                        let barre_chgt = extraction.getElementsByClassName("pct_chargement");
                        if(barre_chgt.length > 0){
                            barre_chgt[0].style.width = extractionEnCours["chargement"] + "%";
                        }
                    }
                    break
                }else{
                    extraction.remove();
                    i--;
                }
            }
    })
    }, 50)
}