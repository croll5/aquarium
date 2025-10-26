# Extraction de journaux au format clé-valeur

Ce module permet d'extraire des données à partir de journaux au format clé-valeur.

## Paramètres complémentaires

Les paramètres complémentaires sont les suivants : 

| Clé | Format | Valeur à renseigner |
|-----|--------|---------------------|
| `encodage` | `utf16` ou `utf8` | L’encodage utilisé par le fichier (par défaut : utf8). |
| `separateur` | `separateur` | Le séparateur entre deux évènements, par exemple : \n |
| `separateur_champs` | `separateur` | Le séparateur entre deux champs d’un même évènement |
| `symbole_association` | `symbole` | Le symbole de séparation entre la clé et la valeur, par exemple `=` |

## Contenu des colonnes

Le contenu des colonnes peut être configuré avec les valeurs suivantes :

| Valeurs possibles | Résultat dans la table |
|-------------------|------------------------|
| `aqua_source` | Le chemin du fichier extrait |
| `nom_cle` | La valeur associée à la clé `nom_cle` dans chaque évènement |
| `nom_cle\|aqua_encodage:encodage` | La valeur de la clé `nom_cle` décodée grâce à l’encodage `encodage` | 


# Exemple

Si l’on souhaite extraire le fichier suivant :
```txt
date=20/08/2008 15:35,type=telechargement,programme=virus.exe,resultat=succes
date=25/08/2008 16:45,type=execution,programme=virus.exe,resultat=erreur
```
On peut utiliser la configuration suivante : 

```xml
<extraction>
    ...
    <table nom="evenements">
        <colonne type="DATETIME" contenu="date|aqua_encodage:date[aqua_sep]02/02/2006 15:04">horodatage</colonne>
        <colonne type="DATETIME" contenu="type">operation</colonne>
        <colonne type="TEXT" contenu="programme">programme</colonne>
        <colonne type="TEXT" contenu="resultat">resultat</colonne>
        <colonne type="TEXT" contenu="aqua_source">source</colonne>
    </table>
    <complements>
        <parametre cle="encodage">utf8</parametre>
        <parametre cle="separateur">\n</parametre>
        <parametre cle="separateur_champs">,</parametre>
        <parametre cle="symbole_association">=</parametre>
    </complements>
    ...
</extraction>
```
