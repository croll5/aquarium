/*
Copyright ou © ou Copr. Cécile Rolland, (21 janvier 2025) 

aquarium[@]mailo[.]com

Ce logiciel est un programme informatique servant à l'analyse des collectes
traçologiques effectuées avec le logiciel DFIR-ORC. 

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

/** Fonction appelée lorsque l'utilisateur appuie sur le bouton "nouvelle analyse"
 */
function nouvelle_analyse(){
    window.location.replace("../html/nouvelle_analyse.html");
}

function analyse_existante(){
    parent.window.go.main.App.OuvrirAnalyseExistante().then(resultat=>{
        if(resultat){
            window.location.replace("../html/synthese.html")
            parent.document.getElementsByTagName("header")[0].style.display = "flex";
        }
    })
}

function nouveau_modele(){
    window.location.replace("../html/nouveau_modele.html");
}

function parametres(){
    window.location.replace("../html/parametres.html")
}

parent.window.go.main.App.VerifierPrerequisDemarrage().then(manquants => {
    if (!Array.isArray(manquants) || manquants.length === 0) {
        return;
    }
    const cadre = document.getElementById("alerte-prerequis");
    const liste = document.getElementById("liste-prerequis-manquants");
    if (!cadre || !liste) {
        return;
    }
    liste.innerHTML = "";
    for (const manquant of manquants) {
        const item = document.createElement("li");
        item.textContent = manquant;
        liste.appendChild(item);
    }
    cadre.style.display = "block";
    const ouvrir_analyse_existante = document.getElementById("ouvrir_analyse_existante");
    if (!ouvrir_analyse_existante) {
        return;
    }
    ouvrir_analyse_existante.setAttribute("disabled", "true")
});
