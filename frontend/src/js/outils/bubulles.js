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

let lastBubbleTime = 0;

/** Code permettant d'afficher de magnifique bulles donnant la sensation de nager parmi les ORCs **/
if(!parent.non_aux_bubulles){
    document.addEventListener('mousemove', function(e) {
        const currentTime = Date.now();
        if (currentTime - lastBubbleTime > 100) { // Délai de 100ms entre chaque bulle
            const bubble = document.createElement('div');
            bubble.classList.add('bubble');
            document.body.appendChild(bubble);

            const size = Math.random() * 20 + 10; // Taille aléatoire entre 10px et 30px
            bubble.style.width = size + 'px';
            bubble.style.height = size + 'px';
            bubble.style.left = e.pageX + 'px';
            bubble.style.top = e.pageY + 'px';

            bubble.addEventListener('animationend', function() {
                bubble.remove();
            });

            lastBubbleTime = currentTime;
        }
    });
    // On supprime l'affichage du pointeur
    document.body.style.cursor = "none";
}