# Extraction des bases de clés de registres Windows

Ce module permet d'extraire des données à partir de bases de données au format clés de registres Windows. 

## Paramètres complémentaires

Les paramètres complémentaires sont les suivants : 

| Clé | Format | Valeur à renseigner |
|-----|--------|---------------------|
| `registre` | `chemin/vers/la/clé` | Le chemin des clés de registre desquelles on souhaite extraire des valeurs. |
| `parcourir_enfants` | `oui` ou `non` | Si l’on souhaite extraire les enfants de la clé ou la valeur de la clé directement |
| `exclusions` | `clé_1;clé_2;...;clé_n` | Une liste de clés à ne pas extraire |

## Contenu des colonnes

Le contenu des colonnes peut être configuré avec les valeurs suivantes :

| Valeurs possibles | Résultat dans la table |
|-------------------|------------------------|
| `<index>:fixe:<index début>:<taille>:<type>` <br> avec :<br>-`index` : l’index de la valeur de clé à extraire<br>-`index début` : l’index de début de la valeur à extraire dans la clé<br>-`taille` : la taille de la valeur à extraire<br>-`type` : le type d’encodage, parmi `string`, `utf16`, `littleEndian64` et `filetime`. | La valeur extraite dans la partie `index` de la clé de registre, entre `index début` et `index début` + `taille`, dans l’encodage `type` |
| `<index>:variable:<index début>:<index taille>:<type>` <br> avec :<br>-`index` : l’index de la valeur de clé à extraire<br>-`index début` : l’index auquel est renseigné l’index de début de la valeur à extraire dans la clé<br>-`taille` : l’index auquel est renseignée la taille de la valeur à extraire<br>-`type` : le type d’encodage, parmi `string`, `utf16`, `littleEndian64` et `filetime`. | La valeur extraite dans la partie `index` de la clé de registre, qui commence à l’index renseigné à `index début` et a une taille définie à l’index `index taille`, dans l’encodage `type` |

## Exemple

Si l’on souhaite extraire de la base SAM, qui contient des informations sur les comptes utilisateurs, l’identifiant et le nom des comptes avec leurs dates de création et de dernière connexion, on pourra utiliser la configuration suivante : 

```xml
<extraction>
    ...
    <complements>
        <parametre cle="registre">SAM/Domains/Account/Users</parametre>
        <parametre cle="parcourir_enfants">oui</parametre>
        <parametre cle="exclusions">Names</parametre>
    </complements>
    <table nom="sam">
        <colonne type="DATETIME" contenu="0:fixe:0x8:0x8:filetime">derniereConnexion</colonne>
        <colonne type="DATETIME" contenu="0:fixe:0x18:0x8:filetime">creationCompte</colonne>
        <colonne type="TEXT" contenu="nomCle">idCompte</colonne>
        <colonne type="TEXT" contenu="1:variable:0xC:0x10:utf16">nomCompte</colonne>
        <colonne type="TEXT" contenu="source">source</colonne>
    </table>
    ...
</extraction>
```

En se basant sur cette configuration, le programme va traiter chaque clé enfant de la clé SAM/Domains/Account/Users, à l’exception de la clé « Names ». 

Voici un exemple de contenu d'une de ces clés : 
```
[F]:
03-00-01-00-00-00-00-00-96-64-C2-FE-2B-1E-D6-01-00-00-00-00-00-00-00-00-5B-A2-44-AA-A5-17-DA-01-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-00-F4-01-00-00-01-02-00-00-15-02-00-00-00-00-00-00-00-00-0A-00-01-00-00-00-00-00-63-00-00-00-06-00

[V]:
00-00-00-00-E8-00-00-00-03-00-01-00-E8-00-00-00-0C-00-00-00-00-00-00-00-F4-00-00-00-00-00-00-00-00-00-00-00-F4-00-00-00-36-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-00-00-00-00-00-00-00-00-2C-01-00-00-08-00-00-00-01-00-00-00-34-01-00-00-18-00-00-00-00-00-00-00-4C-01-00-00-18-00-00-00-00-00-00-00-64-01-00-00-18-00-00-00-00-00-00-00-7C-01-00-00-18-00-00-00-00-00-00-00-01-00-14-80-C8-00-00-00-D8-00-00-00-14-00-00-00-44-00-00-00-02-00-30-00-02-00-00-00-02-C0-14-00-44-00-05-01-01-01-00-00-00-00-00-01-00-00-00-00-02-C0-14-00-FF-FF-1F-00-01-01-00-00-00-00-00-05-07-00-00-00-02-00-84-00-04-00-00-00-00-00-14-00-1B-03-02-00-01-01-00-00-00-00-00-01-00-00-00-00-00-00-18-00-FF-07-0F-00-01-02-00-00-00-00-00-05-20-00-00-00-20-02-00-00-00-00-18-00-FF-07-0F-00-01-02-00-00-00-00-00-05-20-00-00-00-24-02-00-00-00-00-38-00-1B-03-02-00-01-0A-00-00-00-00-00-0F-03-00-00-00-00-04-00-00-DE-A2-28-67-21-3E-D2-AF-19-AD-5D-79-B0-C1-07-29-27-56-FC-20-D8-AD-66-F6-10-F2-68-FA-DF-2A-F8-0F-01-02-00-00-00-00-00-05-20-00-00-00-20-02-00-00-01-02-00-00-00-00-00-05-20-00-00-00-20-02-00-00-49-00-6E-00-76-00-69-00-74-00-E9-00-43-00-6F-00-6D-00-70-00-74-00-65-00-20-00-64-00-19-20-75-00-74-00-69-00-6C-00-69-00-73-00-61-00-74-00-65-00-75-00-72-00-20-00-69-00-6E-00-76-00-69-00-74-00-E9-00-61-00-01-02-00-00-07-00-00-00-05-00-02-00-00-00-00-00-83-13-E5-28-AA-FE-EF-86-2F-5C-A1-4E-75-05-21-39-05-00-02-00-00-00-00-00-81-48-F7-22-E9-D8-F5-A9-A7-30-FB-38-22-4D-DF-44-05-00-02-00-00-00-00-00-E2-45-7C-DA-FB-C6-50-3C-82-10-63-43-F8-08-50-44-05-00-02-00-00-00-00-00-E3-8A-1E-F7-0C-BE-C7-CA-E4-91-96-63-95-D4-FF-3B
```

On retrouve dans la valeur « F », qui correspond à l’index 0, la date de création du compte et de dernière connexion. Pour la date de création, la configuration nous indique qu’elle est renseignée au format « filetime » entre l’octet 0x18 (le 24<sup>ème</sup>) et l’octet 0x18 + 0x08 = 0x20 (le 32<sup>ème</sup>). La date en hexadécimal a donc la valeur suivante : `5B-A2-44-AA-A5-17-DA-01`, ce qui correspond au 15 novembre 2023 à 9 h 25 min 25 s au format filetime.

> 💡 Dans les clés de registre Windows, les nombres sont encodés en Little Endian. `5B-A2-44-AA-A5-17-DA-01` correspondra donc à `01DA17A5AA44A25B` en notation hexadécimale standard.
