# Extraction des bases de données au format SQLite

Ce module permet d'extraire des données à partir de bases de données au format `sqlite`. 

## Paramètres complémentaires

Les paramètres complémentaires sont les suivants : 

| Clé | Valeur à renseigner |
|-----|---------------------|
| `requete_sql` | Requête SQL dont on veut extraire le résultat, se basant sur le fichier à extraire. |

## Contenu des colonnes

Le contenu des colonnes peut être configuré avec les valeurs suivantes :

|                   Valeurs possibles                    |                      Résultat dans la table                      |
|--------------------------------------------------------|------------------------------------------------------------------|
| <i>`nom_colonne`</i> _nom d’une colonne du résultat de la <br>requête définie dans les paramètres complémentaires_ | Valeurs contenues dans la colonne <i>`nom_colonne`</i> <br>du résultat de la requête SQL |
|                      `aqua_source`                      |                     Chemin du fichier source                     |

## Exemple

Si l'on a une base de données contenant la table `executions` suivante : 

| Programme    | DateExecution |
| ------------ | ------------- |
| word.exe     | 20/04/2022    |
| code.exe     | 25/05/2022    |
| virus.bat    | 12/06/2022    |

On pourra utiliser la configuration suivante : 

```xml
<extraction>
    ...
    <table nom="executions_programmes">
        <colonne type="DATETIME" contenu="DateExecution">horodatage</colonne>
        <colonne type="TEXT" contenu="Programme">executable</colonne>
    </table>
    ...
    <complements>
        <parametre cle="requete_sql">SELECT * FROM executions</parametre>
    </complements>
</extraction>
```