affichage_table();

function affichage_table(){
    parent.window.go.main.App.ValeursTableChronologie(0, 5).then(resultat =>{
        let emplacement_resultat = document.getElementById("emplacement_table");
        creer_tableau_depuis_dico(resultat, emplacement_resultat);
    })
}