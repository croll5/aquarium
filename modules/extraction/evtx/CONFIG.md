# Extraction des fichiers évènements Windows

Ce module permet d'extraire des données à partir de fichiers évènements Windows (`.evtx`). 

## Paramètres complémentaires

Il n'y a aucun paramètre complémentaire à configurer.

## Contenu des colonnes

Le contenu des colonnes peut être configuré avec les valeurs suivantes :

| Valeurs possibles | Résultat dans la table |
|-------------------|------------------------|
| <i>`chemin_valeur`</i> _chemin vers une valeur de l’évènement Windows_ | Valeur renseignée dans l’évènement au chemin <i>`chemin_valeur`</i> |
| `source` | Chemin du fichier source |
| `message` | Données complémentaires de l’évènement, au format Json |
| `horodatage` | Heure de l’évènement |

## Exemple

Si l’on a un journal d’évènements Windows contenant l’évènement suivant : 

```xml
<Event xmlns="http://schemas.microsoft.com/win/2004/08/events/event">
    <System>
        <Provider Name="Microsoft-Windows-Security-Auditing" Guid="{56487569-1457-9853-5f78-4568a135b688}" /> 
        <EventID>5061</EventID> 
        ...
        <TimeCreated SystemTime="2025-09-23T15:51:40.9603200Z" /> 
        ...
        <Channel>Security</Channel> 
        <Computer>MON-ORDINATEUR</Computer> 
    </System>
    <EventData>
        <Data Name="SubjectUserName">PirateDangereux</Data> 
        <Data Name="SubjectDomainName">MON-ORDINATEUR</Data> 
        ...
    </EventData>
  </Event>
```
On pourra utiliser la configuration suivante : 

```xml
<extraction>
    ...
    <table nom="executions_programmes">
        <colonne type="DATETIME" contenu="horodatage">horodatage</colonne>
        <colonne type="TEXT" contenu="Event/System/EventID">CodeEvenement</colonne>
        <colonne type="TEXT" contenu="Event/System/Provider/Name">ProviderName</colonne>
        <!-- ⚠️ La ligne suivante ne peut pas encore être utilisée, la fonctionnalité n’étant pas encore implémentée.
        Si vous être frustrés de cette limitation, n’hésitez pas à la développer 😉.-->
        <colonne type="TEXT" contenu="Event/EventData/Data|SubjectUserName">
</table>
    ...
    <!-- Pas de complément -->
    ...
</extraction>
```