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

if(parent.contrastes){
    document.getElementById("contrastes").checked = true;
}
function changer_contrastes(){
    if(document.getElementById("contrastes").checked){
        parent.contrastes = true;
        contrastes_eleves();
    }else{
        parent.contrastes = false;
        contrastes_normaux();
    }
}

if(parent.dyslexie){
    document.getElementById("dyslexie").checked = true;
}
function changer_dyslexie(){
    if(document.getElementById("dyslexie").checked){
        parent.dyslexie = true;
        police_dyslexie()
    }
    else{
        parent.dyslexie = false;
        document.body.style.fontFamily = "Comic Sans Ms"
    }
}

if(parent.non_aux_bubulles){
    document.getElementById("non_aux_bubulles").checked = true;
}
function enlever_bubulles(){
    if(document.getElementById("non_aux_bubulles").checked){
        parent.non_aux_bubulles = true;
    }else{
        parent.non_aux_bubulles = false;
    }
}

function quitter_parametres(){
    window.location.replace("accueil.html");
}

function contrastes_normaux(){
    let couleurs = document.documentElement;
    couleurs.style.setProperty('--bleu-1', '#001D64');
    couleurs.style.setProperty('--bleu-2', '#002C9A');
    couleurs.style.setProperty('--bleu-3', '#0045F2');
    couleurs.style.setProperty('--bleu-4', '#5787FF');
    couleurs.style.setProperty('--bleu-5', '#ABC3FF');
    couleurs.style.setProperty('--bleu-6', '#E5ECFF');
    couleurs.style.setProperty('--or-1', '#856B0D');
    couleurs.style.setProperty('--or-2', '#CAA314');
    couleurs.style.setProperty('--or-3', '#EBC331');
    couleurs.style.setProperty('--or-4', '#F1D46B');
    couleurs.style.setProperty('--or-5', '#F6E4A4');
    couleurs.style.setProperty('--or-6', '#FCF5DC');
}
