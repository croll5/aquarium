/** Fonction appelée lorsque l'utilisateur appuie sur le bouton "nouvelle analyse"
 */
function nouvelle_analyse(){
    window.location.replace("../html/nouvelle_analyse.html");
}

function analyse_existante(){
    parent.window.go.main.App.OuvrirAnalyseExistante().then(resultat=>{
        if(resultat){
            window.location.replace("../html/extraction.html")
            //parent.document.getElementsByTagName("header")[0].style.display = "inline";
            //let onglet_courant = parent.document.getElementById("onglet_extraction");
            //onglet_courant.style.backgroundColor = "#FCF5DC";
            //onglet_courant.style.color = "#000";
        }
    })
}

function nouveau_modele(){
    window.location.replace("../html/nouveau_modele.html");
}