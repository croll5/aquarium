# Comment configurer aquarium ?

Aquarium peut être configuré à l'aide d'un fichier `config/config.xml`, contenu par défaut dans le dossier de l'exécutable. Ce fichier a la structure suivante : 

```xml
<?xml version="1.0" encoding="UTF-8"?>
<extractions>
    <collecte><!--Nom du dossier dans lequel sera enregistrée la collecte--></collecte>
    <!--Définition de la structure de la table rassemblant tous les évènements-->
    <chronologie nom="..."> <!--nom : le nom qu'aura la table SQL-->
        <colonne type="..."><!--Nom de l'attribut (ie de la colonne SQL --></colonne> <!-- type : le type de l'attribut (par exemple INT, TEXT, DATETIME, ...). Il doit s'agit d'un type SQL valide -->
    </chronologie>
    <extraction><!--Nom du fichier XML de définition de l'extraction, par exemple « getthis » pour le fichier extraactions/getthis.xml --></extraction>
</extractions>
```

> 💡 Par défaut, les fichiers de configuration sont cherchés dans le dossier d'analyse. S'ils ne s'y trouvent pas, ils sont cherchés dans le dossier de l'exécutable

Les fichiers de configuration des extractions ont la structure suivante : 

```xml
<extraction>
    <id><!--Identifiant unique de l'extraction--></id>
    <extracteur><!--Identifiant de l'extracteur. La liste des identifiants utilisables et de leurs particularités est décrite dans la section 2--></extracteur>
    <nom><!--Nom de l'extraction, qui affiché dans l'interface graphique--></nom>
    <description><!--Description de l'extraction--></description>
    <table nom="..."> <!--nom : le nom de la table dans laquelle seront enregistrées les données extraites-->
        <colonne type="..." contenu="..."><!--nom de l'attribut--></colonne> 
        <!--type : type de l'attribut (par exemple INT, TEXT, DATETIME, ...). Il doit s'agit d'un type SQL valide-->
        <!--contenu : définition de la donnée à extraire, propre à chaque extracteur (cf : section 2)-->
    </table>
    <chemin>
        <dossier><!--nom du dossier qui doit être parcouru (cette balise peut être présente plusieurs fois si les fichiers recherchés sont dans des dossiers imbriqués--></dossier>
        <archive><!--optionnel : nom de l'archive dans laquelle sont enregistrés les fichiers à extraire--></archive>
        <fichier><!--Chemin du ou des fichiers dans l'archive ou nom du ou des fichiers dans le dossier--></fichier>
    </chemin>
    <complements>
        <parametre cle="..."><!--Valeur d'un ou plusieurs éventuels paramètres complémentaires--></parametre> 
        <!--cle : cle des éventuels paramètres complémentaires-->
    </complements>
    <sql_chronologie><!--Requete SQL permettant de récupérer les données à ajouter dans la table Chronologie--></sql_chronologie>

</extraction>
```
En fonction de l’extracteur utilisé, des paramètres spécifiques sont attendus dans la configuration des extractions. 

Voici la liste des extracteurs utilisables et leur documentation
- [sqlite](../extraction/bdd_sqlite/CONFIG.md) : extraction des bases de données au format SQLite
- [csv](../extraction/csv/CONFIG.md) : extraction des fichiers au format CSV
- [evtx](../extraction/evtx/CONFIG.md) : extraction des fichiers d'évènements Windows (evtx)
- [registre](../extraction/registre/CONFIG.md) : extraction des bases de clés de registre Windows
- [prefetch](../extraction/prefetch/CONFIG.md) : extraction des données des fichiers de préchargement