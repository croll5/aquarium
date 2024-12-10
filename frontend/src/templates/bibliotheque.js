if(parent.contrastes){
    contrastes_eleves();
}
if(parent.dyslexie){
    police_dyslexie();
}

function contrastes_eleves(){
    let couleurs = document.documentElement;
    couleurs.style.setProperty('--bleu-1', '#000');
    couleurs.style.setProperty('--bleu-2', '#000');
    couleurs.style.setProperty('--bleu-3', '#000');
    couleurs.style.setProperty('--bleu-4', '#FFF');
    couleurs.style.setProperty('--bleu-5', '#FFF');
    couleurs.style.setProperty('--bleu-6', '#FFF');
    couleurs.style.setProperty('--or-1', '#000');
    couleurs.style.setProperty('--or-2', '#000');
    couleurs.style.setProperty('--or-3', '#FFF');
    couleurs.style.setProperty('--or-4', '#FFF');
    couleurs.style.setProperty('--or-5', '#FFF');
    couleurs.style.setProperty('--or-6', '#FFF');
}

function police_dyslexie(){
    document.body.style.fontFamily = "Open Dyslexic";
}