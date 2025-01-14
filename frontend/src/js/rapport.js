afficher_pistes()

function afficher_pistes(){
    parent.window.go.main.App.ListePistesRapport().then(resultat =>{
        let contenant_pistes = document.getElementById("liste_pistes");
        for(let ligne of resultat){
            let detail = document.createElement("details");
            detail.id = ligne["id"];
            detail.className = "details_piste";
            detail.addEventListener("toggle", (event) => {
                if(detail.childElementCount < 3){
                    afficher_etapes_piste(detail, ligne["id"]);
                }
              });
            // On ajoute le titre
            let titre = document.createElement("summary");
            titre.textContent = ligne["titre"];
            detail.appendChild(titre);
            // On ajoute la description
            let description = document.createElement("p");
            description.textContent = ligne["description"];
            detail.appendChild(description);
            contenant_pistes.appendChild(detail);
        }
    })
}

function afficher_etapes_piste(div_piste, idLigne){
    parent.window.go.main.App.ListeEtapesRapport(idLigne).then(async resultat =>{
        for(let etape of resultat){
            // On affiche une introduction à la requete SQL
            let introRequeteSQL = document.createElement("p");
            introRequeteSQL.textContent = "Requête SQL :";
            div_piste.appendChild(introRequeteSQL);
            // On affiche la requête SQL
            let requeteSQL = document.createElement("p");
            requeteSQL.className = "code";
            requeteSQL.textContent = etape["RequeteSQL"];
            div_piste.appendChild(requeteSQL);
            console.log(etape);
            // On affiche une introduction au résultat de la requête
            let introTable = document.createElement("p");
            introTable.textContent = "Lignes significatives du résultat : ";
            div_piste.appendChild(introTable);
            // On affiche la table
            let lignesTable = await parent.window.go.main.App.DonneesTableRapport(etape["NomTable"]);
            creer_tableau_depuis_dico(lignesTable, div_piste, false);
            // On prépare l'affichage du commentaire de l'analyste
            let introCommentaire = document.createElement("p");
            introCommentaire.textContent = "Commentaire de l'analyste :";
            div_piste.appendChild(introCommentaire);
            // On affiche le commentaire de l'analyste
            let commentaireAnalyste = document.createElement("p");
            commentaireAnalyste.textContent = etape["Commentaire"];
            div_piste.appendChild(commentaireAnalyste);
            console.log(lignesTable);
        }
    })
}

function nouvelle_piste(){
    let titre = document.getElementById("titre_nvelle_piste").value;
    let description = document.getElementById("description_nvelle_piste").value;
    parent.window.go.main.App.AjouterPisteDansRapport(titre, description).then(() =>{
        window.location.reload();
    })
}

function exporter_rapport(){
    window.print();
}

