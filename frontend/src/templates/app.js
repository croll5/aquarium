/* ***********************************************/
/* ************ LOAD TEMPLATE FONCTIONS **********/
/* ***********************************************/

function loadScripts(scripts) {
    scripts.forEach(scriptSrc => {
        const script = document.createElement('script');
        script.src = scriptSrc;
        document.head.appendChild(script);
    });
}

fetch('../templates/template.html')
    .then(response => response.text())
    .then(data => {
        let parser = new DOMParser();
        let doc = parser.parseFromString(data, "text/html");

        // Insérer les éléments du head
        let templateHeadElements = doc.querySelectorAll('head > *');
        templateHeadElements.forEach(element => {
            document.head.prepend(element);
        });

        // Insérer le contenu des balises header et footer
        document.querySelector('header').insertAdjacentHTML('afterbegin', doc.querySelector('header').innerHTML);
        document.querySelector('footer').insertAdjacentHTML('afterbegin', doc.querySelector('footer').innerHTML);

        // Charger les scripts spécifiques après l'insertion des templates
        loadScripts(['/wails/ipc.js', '/wails/runtime.js']);
    })
    .catch(error => console.error('Erreur lors du chargement du template :', error));




/* **************************************/
/* *********** HEADER FONCTIONS *********/
/* **************************************/

function accueil(){
    window.location.replace("/")
}

function change_onglet(destination){
    window.location.replace(destination);
    // On met tous les autres onglets à la couleur standard
    /*let onglets = document.getElementsByClassName("onglet");
    for (const onglet of onglets) {
        onglet.style.backgroundColor = "#856B0D";
        onglet.style.color = "#fff";
    }
    // On met l'onglet sur lequel on va à la couleur de la page
    let onglet_courant = parent.document.getElementById(id_onglet);
    onglet_courant.style.backgroundColor = "#FCF5DC";
    onglet_courant.style.color = "#000";*/
}


/* *******************************************/
/* *********** TEST NEWPAGE FONCTION *********/
/* *******************************************/

// Test for newPage.html
window.blancPageFunction = function () {
    let nameElement = document.getElementById("name");
    let resultElement = document.getElementById("result");
    let name = nameElement.value;
    if (name === "") return;
    try {
        parent.window.go.main.App.BlancPageFunction(name)
            .then((result) => {
                // Update result with data back from App.Greet()
                resultElement.innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};


/* **************************************/
/* ************ UTILS FONCTION **********/
/* **************************************/

function executeWhenReady(callback) {
    let checkInterval = setInterval(function() {
        try {
            if (parent.window.go && parent.window.go.main && parent.window.go.main.App) {
                callback(); // Exécute la fonction passée en paramètre
                clearInterval(checkInterval); // Arrête l'intervalle une fois les éléments prêts
            }
        } catch (err) {
            console.error("Error checking parent window:", err);
        }
    }, 200); // Vérifie toutes les Xms
}


