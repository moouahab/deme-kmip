# PKCS#11

## Definition

PKCS#11 est une API standard pour interagir avec des objets cryptographiques stockes dans un token ou un HSM.

Il sert souvent a utiliser :

- un HSM ;
- une carte a puce ;
- un token USB ;
- un module logiciel compatible ;
- un coffre cryptographique local.

## Difference entre KMIP et PKCS#11

KMIP et PKCS#11 ne resolvent pas exactement le meme probleme.

| Element | Type | Role |
|---|---|---|
| KMIP | Protocole reseau | Gerer des objets cryptographiques a distance via un KMS |
| PKCS#11 | API locale | Utiliser des objets cryptographiques dans un token ou HSM |

Formulation simple :

```txt
KMIP    = parler a un serveur KMS
PKCS#11 = appeler une librairie locale pour utiliser un token/HSM
```

## Exemple d'utilisation PKCS#11

Une application peut charger une librairie fournie par un HSM :

```txt
application
  |
  | appels PKCS#11
  v
librairie fournisseur .so/.dll
  |
  v
HSM ou token
```

Les appels PKCS#11 peuvent demander :

- ouvrir une session ;
- se connecter avec un PIN ;
- generer une cle ;
- chercher un objet ;
- signer ;
- verifier ;
- chiffrer ;
- dechiffrer ;
- detruire un objet.

## Architecture avec KMIP et PKCS#11

Un KMS peut exposer KMIP cote reseau et utiliser PKCS#11 cote backend pour parler a un HSM.

```txt
Client KMIP
  |
  | reseau, KMIP
  v
Serveur KMS
  |
  | API locale, PKCS#11
  v
HSM
```

Dans cette architecture :

- le client ne connait pas PKCS#11 ;
- le client parle uniquement KMIP ;
- le KMS applique les regles metier ;
- le KMS utilise PKCS#11 pour deleguer les operations sensibles au HSM.

## Relation avec ce projet

Le projet actuel n'implemente pas PKCS#11.

Il simule seulement la couche KMS avec un repository memoire :

```txt
internal/kms/memory_repository.go
```

Pour ajouter PKCS#11 plus tard, il faudrait probablement creer une nouvelle implementation de repository ou de backend cryptographique :

```txt
internal/kms/pkcs11_repository.go
```

ou une couche separee :

```txt
internal/hsm/pkcs11_client.go
```

## Ce qu'un backend PKCS#11 devrait gerer

Une integration PKCS#11 reelle devrait gerer :

- chargement de la librairie PKCS#11 ;
- initialisation du module ;
- slots ;
- tokens ;
- sessions ;
- login utilisateur ou security officer ;
- generation d'objets ;
- recherche par label ou ID ;
- attributs PKCS#11 ;
- erreurs fournisseur ;
- concurrence ;
- fermeture propre des sessions.

## Comparaison des objets

KMIP et PKCS#11 ont chacun leur vocabulaire.

| Concept KMIP | Concept PKCS#11 proche |
|---|---|
| ObjectType | Object class |
| UniqueIdentifier | Object handle, label ou ID |
| Attributes | Attributes PKCS#11 |
| Create | GenerateKey / CreateObject |
| Destroy | DestroyObject |
| GetAttributes | GetAttributeValue |

Ce mapping n'est pas toujours direct. Un vrai connecteur doit traduire proprement les concepts.

## Pourquoi garder KMIP devant PKCS#11

Mettre KMIP devant PKCS#11 peut etre utile parce que :

- les clients parlent un protocole reseau standard ;
- le KMS peut centraliser les politiques ;
- le KMS peut journaliser toutes les operations ;
- le KMS peut cacher les details fournisseur du HSM ;
- le backend HSM peut changer sans changer les clients.

## Limites et risques

PKCS#11 est puissant, mais plus bas niveau que KMIP.

Points sensibles :

- mauvaise gestion des sessions ;
- erreurs de droits sur les objets ;
- attributs mal configures ;
- cles exportables par erreur ;
- concurrence mal controlee ;
- dependance aux specificites du fournisseur HSM.

Pour un projet pedagogique, il est raisonnable de commencer sans PKCS#11.
