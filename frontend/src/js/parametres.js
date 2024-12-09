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
