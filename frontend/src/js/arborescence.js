
window.onload = function() {
    executeWhenReady(function() {
        construireArborescence("arborescence", []);
    });
};

/** Fonction permettant d'afficher le contenu d'un dossier
 *
 * @param {*} id_racine identifiant du dossier duquel on veut afficher le contenu
 * @param {*} chemin_num chemin dans l'arborescence (json) sur laquelle on se base
 * @returns
 */
function construireArborescence(id_racine, chemin_num){
    // On récupère un pointeur vers le dossier duquel on veut afficher le contenu
    let racine = document.getElementById(id_racine);
    // Si le dossier a déjà un contenu, inutile d'en ré-extraire le contenu. On s'arrête là
    if(racine == undefined || racine.children.length > 1){
        return
    }
    try {
        // On indique à l'utilisateur qu'il faut patienter
        document.body.style.cursor = "wait";
    } catch (error) {
    }
    // On interroge une fonction go qui renvoie une liste contenant les métadonnées des fichiers
    // et dossiers contenus dans le dossier concerné
    parent.window.go.main.App.ArborescenceMachineAnalysee(chemin_num).then(resultat =>{
        // Si l'on a aucun résultat, cela signifie que l'arborescence n'a pas encore été extraite.
        // On affiche donc un menu permettant à l'utilisateur de lancer l'extraction
        if(resultat.length == 0){
            document.getElementById("extraction_arborescence").style.display = "inline";
            document.getElementById("patientez").style.display = "none";
            return
        }
        // On cherche de quels indicateurs de légitimité l'utilisateur demande l'affichage
        let afficher_inconnu, afficher_ok, afficher_suspect;
        try{
            afficher_inconnu = document.getElementById("affiche_inconnu").checked;
            afficher_ok = document.getElementById("affiche_ok").checked;
            afficher_suspect = document.getElementById("affiche_suspect").checked;
        }catch(error){
            // Par défault s'il y a une erreur, on affiche tout
            alert(document.getElementById("affiche_inconnu"));
            afficher_inconnu = true;
            afficher_ok = true;
            afficher_suspect = true;
        }
        // On ajoute les fichiers et dossiers contenus dans le dossier concerné
        for(let i=0; i < resultat.length; i++){
            // Si le fichier a des enfants, on l'affiche comme un dossier
            if(resultat[i]["ADesEnfants"]){
                let enfant = document.createElement("details");
                enfant.id = String.prototype.concat(id_racine, "_", i);
                enfant.className = "dossier_arborescence";
                let chemin_enfant = chemin_num.concat([i]);
                // L'élément titre_enfant contient le nom du dossier
                let titre_enfant = document.createElement("summary");
                titre_enfant.textContent = resultat[i]["Nom"];
                // Lorsque l'on cliquera sur ce dossier, ses enfants seront affichés grâce à cette même fonction
                titre_enfant.onclick = function(ev){return construireArborescence(enfant.id, chemin_enfant)};
                enfant.appendChild(titre_enfant);
                // On ajoute le sous-dossier dans le dossier
                racine.appendChild(enfant);
            }else{
                // Sinon, il s'agit d'un fichier
                let enfant = document.createElement("p");
                // On affiche la légitimité (si le fichier est présent dans le modèle)
                let legitimite = document.createElement("strong");
                legitimite.textContent = "😇";
                legitimite.className = "legitimite_ok";
                legitimite.style.display = afficher_ok ? "inline" : "none";
                if(resultat[i]["EnfantsSuspects"] > 0){
                    // Si le fichier a une empreinte différente de celle dans le modèle
                    legitimite.textContent = "🥴";
                    legitimite.className = "legitimite_suspect";
                    legitimite.style.display = afficher_suspect ? "inline" : "none";
                }else if(resultat[i]["EnfantsInconnus"] > 0){
                    // Si le fichier n'existe pas dans le modèle
                    legitimite.textContent = "😵";
                    legitimite.className = "legitimite_aucune";
                    legitimite.style.display = afficher_inconnu ? "inline" : "none";
                }
                enfant.id = String.prototype.concat(id_racine, "_", i);
                enfant.className = "fichier_arborescence";
                // On ajoute le nom du fichier
                enfant.textContent = resultat[i]["Nom"];
                enfant.appendChild(legitimite);
                // On ajoute le fichier dans le dossier
                racine.appendChild(enfant);
            }
        }
        // On remet le curseur standard pour indiquer que les calculs sont achevés
        document.body.style.cursor = "default";
        // On affiche l'arborescence et sa légende
        document.getElementById("affichage_arbo").style.display = "inline";
    })

}
/** Fonction permettant d'extraire l'arborescence d'un ORC en faisant
 *  appel à la fonction Go ExtraireArborescence
 */
function extraire_arborescence(){
    // On masque le menu permettant d'extraire l'arborescence
    document.getElementById("extraction_arborescence").style.display = "none";
    // On affiche la ligne demandant de patienter
    document.getElementById("patientez").style.display = "inline";
    // On regarde si l'arborescence doit être extraite avec un modèle
    let avec_modele = document.getElementById("avec_modele").checked;
    parent.window.go.main.App.ExtraireArborescence(avec_modele).then(resultat =>{
        // On fois que l'arborescence a été extraite, on l'affiche
        document.getElementById("patientez").style.display = "none";
        construireArborescence("arborescence", []);
    })
}

/** Fonction permettant d'afficher ou de masquer les indicateurs de légitimité
 * demandés par l'utilisateur
 * Cette fonction se déclenche lorsque l'utilisateur coche ou décoche une case de la secion "légende"
 * @param {*} id_checkbox : l'identifiant de la "checkbox" que l'utilisateur a changée
 * @param {*} nom_classe : non de la classe des éléments à afficher ou masquer
 * (par exemple "legitimite_aucune" pour les fichiers n'étant pas présents dans le modèle)
 */
function affichage_legitimite(id_checkbox, nom_classe){
    let choix = document.getElementById(id_checkbox);
    // On récupère la liste des indicateurs
    let smileys = document.getElementsByClassName(nom_classe);
    // On regarde si la case a été cochée ou décochée
    if (choix.checked){
        for(const element of smileys){
            // Si elle a été cochée, on affiche les indicateurs
            element.style.display = "inline";
        }
    } else{
        for(const element of smileys){
            // Si elle a été décochée, on masque les indicateurs
            element.style.display = "none";
        }
    }
}

