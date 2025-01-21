/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025) 

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à [rappeler les
caractéristiques techniques de votre logiciel]. 

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
            if(etape["RequeteSQL"] == ""){
                continue;
            }
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

function exporter_rapport() {
    let details = document.getElementsByTagName("details");
    for(let piste of details){
        piste.setAttribute("open", true);
    }
    document.body.style.cursor = "wait";
    setTimeout(function(){
        window.print();
        document.body.style.cursor = "default";
    }, 5000);
}