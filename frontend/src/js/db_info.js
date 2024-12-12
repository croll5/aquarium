var sqlTable_list = {}

window.onload = function() {
        print_db_infos();
};


function print_db_infos() {
    parent.window.go.main.App.Get_db_info().then(resultat => {
        //console.log(resultat);
        let table_name = document.getElementById("table_selection")
        let div_db_infos = document.getElementById("db_infos");
        sqlTable_list = {};
        for (let [key, value] of Object.entries(resultat)) {
            sqlTable_list[key] = value;

            // Add new element in the selectBox
            let option = document.createElement('option');
            option.innerText = key;
            option.value = key;
            table_name.appendChild(option)
        }
        div_db_infos.innerHTML = "Number of tables: " + Object.keys(sqlTable_list).length ;
        print_table_head()
    }).catch(err => console.error("Error db_info:", err));

}
function print_table_head() {
    let table_name = document.getElementById("table_selection").value;
    let rows_limit = document.getElementById("rows_limit").value;
    console.log(sqlTable_list[table_name])

    document.getElementById("table_info").innerHTML = sqlTable_list[table_name];
    parent.window.go.main.App.Get_header_table(table_name, rows_limit).then(resultat => {
        //console.log(resultat);
        console.table(resultat);
        let div_db_infos = document.getElementById("table_values");
        div_db_infos.innerHTML = ''
        // Créer un conteneur avec une barre de défilement horizontal
        let scrollContainer = document.createElement('div');
        scrollContainer.style.maxHeight = '450px'
        scrollContainer.style.overflowX = 'auto';
        scrollContainer.style.overflowY = 'auto';


        // Créer un tableau Bootstrap
        let table = document.createElement('table');
        table.className = 'table table-striped table-bordered';
        table.style.fontSize = 'smaller'; // Réduire la taille du texte

        // Créer l'en-tête du tableau
        let thead = document.createElement('thead');
        let headerRow = document.createElement('tr');

        // Assumer que le premier objet dans résultat contient les colonnes
        let firstRow = Object.values(resultat)[0];
        for (let key in firstRow) {
            let th = document.createElement('th');
            th.innerText = key;
            headerRow.appendChild(th);
        }
        thead.appendChild(headerRow);
        table.appendChild(thead);

        // Créer le corps du tableau
        let tbody = document.createElement('tbody');
        for (let [idRow, valueRows] of Object.entries(resultat)) {
            let tr = document.createElement('tr');
            for (let [key, value] of Object.entries(valueRows)) {
                let td = document.createElement('td');
                td.innerText = value;
                tr.appendChild(td);
            }
            tbody.appendChild(tr);
        }
        table.appendChild(tbody);

        scrollContainer.appendChild(table);
        div_db_infos.appendChild(scrollContainer);
    }).catch(err => console.error("Error db_info:", err));
}
