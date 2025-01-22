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

function choix_enregistrement(){
    parent.window.go.main.App.CreationNouveauProjet().then(resultat =>{
        document.getElementById("enregistrement").value = resultat;
        document.getElementById("archives").value = "";
        document.getElementById("valider").style.display = "none";
    })
}

function choix_orc(){
    let chemin_enreg = document.getElementById("enregistrement").value;
    if(chemin_enreg == ""){
        alert("Vous deviez d'abord choisir où vous voulez enregistrer votre ORC")
    } 
    else{
        document.getElementById("patientez").style.display = "inline";
        document.getElementById("formulaire").style.display = "none";
        parent.window.go.main.App.AjoutORCNouveauProjet().then(resultat =>{
            document.getElementById("archives").value = resultat;
            if(resultat == ""){
                document.getElementById("valider").style.display = "none";
            }
            else{
                let nom_auteur = document.getElementById("auteur").value;
                if(nom_auteur != ""){
                    document.getElementById("valider").style.display = "inline";
                }
            }
            document.getElementById("patientez").style.display = "none";
            document.getElementById("formulaire").style.display = "inline";
        })
    }
}

function change_auteur(){
    let nom_auteur = document.getElementById("auteur").value;
    let enregistrement = document.getElementById("enregistrement").value;
    let archives = document.getElementById("archives").value;
    if(nom_auteur == "" || enregistrement == "" || archives == ""){
        document.getElementById("valider").style.display = "none";
    }
    else{
        document.getElementById("valider").style.display = "inline";
    }
}

function validation(){
    let auteur = document.getElementById("auteur").value;
    let description = document.getElementById("description").value;
    parent.window.go.main.App.ValidationCreationProjet(auteur, description).then(resultat =>{
        if(resultat){ 
            window.location.replace("../html/extraction.html");
            parent.document.getElementsByTagName("header")[0].style.display = "inline";
            let onglet_courant = parent.document.getElementById("onglet_extraction");
            onglet_courant.style.backgroundColor = "#FCF5DC";
            onglet_courant.style.color = "#000";
        }else{
            window.location.replace("../html/accueil.html")
        }
    })
}