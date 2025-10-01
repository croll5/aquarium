# Extraction d'un fichier CSV

Ce module permet d'extraire des données à partir d'un fichier CSV.

## Paramètres complémentaires

Il n'y a aucun paramètre complémentaire à configurer.

## Contenu des colonnes

Le contenu des colonnes peut être configuré avec les valeurs suivantes :

|                   Valeurs possibles                    |                      Résultat dans la table                      |
|--------------------------------------------------------|------------------------------------------------------------------|
| `<nom_colonne>` _en-tête d'une colonne du fichier CSV_ | Valeurs contenues dans la colonne `<nom_colonne>` du fichier CSV |
|                        `source`                        |                     Chemin du fichier source                     |

## Exemple

Si l'on a un fichier CSV contenant les données suivantes : 

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
</extraction>
```