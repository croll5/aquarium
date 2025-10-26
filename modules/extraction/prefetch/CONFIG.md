# Extraction fichiers de préchargement Windows

Ce module permet d'extraire des données à partir des fichiers de préchargement Windows. Il s’agit de fichiers utilisés par Windows pour charger les ressources des pregrammes exécutés fréquemment avant leur exécution, afin de réduire le temps de chargement. Ces fichiers donnent des informations précieuses sur les programmes exécutés par la machine et les ressources utilisées par ceux-ci. 

## Paramètres complémentaires

Il n'y a aucun paramètre complémentaire à configurer.

## Contenu des colonnes

Le contenu des colonnes peut être configuré avec les valeurs suivantes :

| Valeurs possibles | Résultat dans la table |
|-------------------|------------------------|
| `aqua_source` | Le chemin du fichier qui a été extrait |
| `executable` | Le nom de l’exécutable |
| `taille_fichier` | La taille de l'exécutable |
| `empreinte` | L'empreinte (hash en anglais) de l’exécutable |
| `nb_executions` | Le nombre de fois où le pregramme a été exécuté |
| `version` | Le numéro de version de l’exécutable |
| `date_execution` | Les dates des dernières exécutions du programme |
| `ressource` | L'emplacement des ressources utilisées par le programme |


## Attribut « condition » des tables

L’attribut « condition » peut prendre les valeurs suivantes : 
| Valeur | Effet sur la table |
| `date_execution` | La table contiendra une ligne par date d’exécution de chaque programme |
| `ressource` | La table contiendra une ligne par ressource utilisée par chaque programme |

## Exemple

Si pour un programme `virus.exe` s’est exécuté le 15 avril 2005 à 14 h 23 et à 15 h 17 et a pour numéro de version 2.3, la configuration suivante : 

```xml
<?xml version="1.0" encoding="UTF-8"?>
<extraction>
    ...
    <table nom="prechargement" condition="dates_executions">
        <colonne type="TEXT" contenu="executable">executable</colonne>
        <colonne type="TEXT" contenu="date_execution">date_execution</colonne>
        <colonne type="DATETIME" contenu="version">version</colonne>
    </table>
    ...
</extraction>
```

Permettra d’obtenir la table `prechargement` suivante :


| executable | date_execution | version |
|------------|----------------|---------|
| virus.exe | 2005-04-15 14:23:00 +0200 CEST | 2.3 |
| virus.exe | 2005-04-15 15:17:00 +0200 CEST | 2.3 |