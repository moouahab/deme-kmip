# Comprendre KMIP

## Présentation

**KMIP** signifie **Key Management Interoperability Protocol**.

C’est un protocole standard utilisé pour permettre à différents systèmes de communiquer avec un serveur de gestion de clés, appelé **KMS** pour **Key Management Service**.

L’objectif de KMIP est simple :

> Permettre à des clients différents de gérer des objets cryptographiques de manière normalisée, sans dépendre d’un fournisseur unique.

---

## À quoi sert KMIP ?

Dans une infrastructure, plusieurs applications peuvent avoir besoin de clés cryptographiques :

- chiffrement de données ;
- chiffrement de disques ;
- gestion de certificats ;
- sécurisation de bases de données ;
- stockage sécurisé ;
- services cloud ;
- HSM ;
- sauvegardes chiffrées.

Sans standard, chaque fournisseur pourrait avoir sa propre API, son propre format et sa propre manière de gérer les clés.

KMIP apporte une façon commune de demander au serveur :

- de créer une clé ;
- de récupérer une clé ;
- de détruire une clé ;
- d’activer une clé ;
- de révoquer une clé ;
- de localiser une clé ;
- de gérer les attributs d’un objet cryptographique.

---

## Différence entre KMS et KMIP

Il faut bien séparer les deux notions.

### KMS

Le **KMS** est le serveur ou le service qui gère les objets cryptographiques.

Il est responsable de :

- créer les clés ;
- stocker les clés ;
- protéger les clés ;
- contrôler les accès ;
- gérer le cycle de vie des clés ;
- journaliser les opérations ;
- exposer des métriques ;
- appliquer les règles de sécurité.

```txt
KMS = le coffre-fort / le serveur de gestion des clés
```

### KMIP

**KMIP** est le protocole de communication utilisé pour parler avec le KMS.

Il définit :

- le format des messages ;
- les opérations possibles ;
- les objets manipulés ;
- les attributs ;
- les réponses ;
- les erreurs.

```txt
KMIP = le langage normalisé utilisé pour parler au coffre-fort
```

---

## Image simple

On peut comparer le KMS à une chambre noire sécurisée.

```txt
Client
  ↓
Requête KMIP
  ↓
KMS
  ↓
Réponse KMIP
```

Le client demande une action, par exemple :

```txt
Crée une clé symétrique
```

Le KMS exécute l’action et retourne une réponse, par exemple :

```txt
Clé créée avec l’identifiant key-123
```

Le client n’a pas besoin de connaître toute l’implémentation interne du KMS.

---

## Les objets cryptographiques

KMIP ne gère pas uniquement des clés simples.

Il peut représenter plusieurs types d’objets cryptographiques :

- clé symétrique ;
- clé publique ;
- clé privée ;
- certificat ;
- secret data ;
- opaque object ;
- split key ;
- template.

Exemple :

```txt
ObjectType = SymmetricKey
ObjectType = PublicKey
ObjectType = PrivateKey
ObjectType = Certificate
```

---

## Le cycle de vie d’une clé

Une clé cryptographique ne reste pas toujours dans le même état.

Elle suit un cycle de vie.

Exemple simplifié :

```txt
Create
  ↓
Pre-Active
  ↓
Activate
  ↓
Active
  ↓
Revoke
  ↓
Revoked
  ↓
Destroy
  ↓
Destroyed
```

Dans un prototype simple, on peut commencer avec :

```txt
Create
  ↓
Active
  ↓
Destroy
  ↓
Destroyed
```

---

## Opérations KMIP courantes

KMIP définit plusieurs opérations.

Voici les plus importantes pour comprendre le protocole :

| Opération | Rôle |
|---|---|
| `Create` | Créer un nouvel objet cryptographique |
| `Get` | Récupérer un objet ou ses informations |
| `Destroy` | Détruire un objet cryptographique |
| `Activate` | Activer une clé |
| `Revoke` | Révoquer une clé |
| `Locate` | Rechercher des objets |
| `Query` | Demander les capacités du serveur |
| `Register` | Enregistrer un objet déjà existant |
| `GetAttributes` | Lire les attributs d’un objet |
| `AddAttribute` | Ajouter un attribut à un objet |

---

## Qu’est-ce que TTLV ?

KMIP utilise un format binaire appelé **TTLV**.

**TTLV** signifie :

```txt
Tag - Type - Length - Value
```

Chaque morceau du message est représenté avec ces quatre éléments.

---

## Structure TTLV

Un bloc TTLV contient :

```txt
Tag    = identifie le champ
Type   = indique le type de donnée
Length = indique la taille de la valeur
Value  = contient la valeur
```

Exemple simplifié :

```txt
Tag    : Operation
Type   : Enumeration
Length : 4
Value  : Create
```

Cela signifie :

```txt
Ce bloc indique que l’opération demandée est Create.
```

---

## Exemple de requête simplifiée

Un client veut créer une clé symétrique.

La requête peut contenir :

```txt
Operation  = Create
ObjectType = SymmetricKey
```

Flux simplifié :

```txt
Client
  ↓
Message TTLV
  ↓
Serveur KMIP/KMS
  ↓
Décodage du message
  ↓
Création de la clé
  ↓
Audit log
  ↓
Réponse
```

Réponse simplifiée :

```json
{
  "key_id": "key-123",
  "status": "active"
}
```

---

## Exemple de décodage logique

Un serveur KMIP peut recevoir un message binaire.

Il doit ensuite :

1. lire le `Tag` ;
2. lire le `Type` ;
3. lire la `Length` ;
4. lire la `Value` ;
5. interpréter la valeur ;
6. transformer le message en opération métier.

Exemple :

```txt
TagOperation + TypeEnumeration + ValueCreate
        ↓
OperationRequest{ Operation: Create }
```

---

## Pourquoi KMIP est important ?

KMIP est important parce qu’il apporte de l’interopérabilité.

Cela permet à plusieurs clients et produits de communiquer avec le même KMS de façon standard.

Sans KMIP :

```txt
Client A → API spécifique fournisseur A
Client B → API spécifique fournisseur B
Client C → API spécifique fournisseur C
```

Avec KMIP :

```txt
Client A
Client B
Client C
   ↓
Protocole KMIP
   ↓
KMS
```

---

## Sécurité autour de KMIP

KMIP est lié à des environnements sensibles, car il manipule des objets cryptographiques.

Un vrai serveur KMIP/KMS doit prendre en compte :

- TLS ou mTLS ;
- authentification forte ;
- contrôle d’accès ;
- audit log ;
- séparation des rôles ;
- stockage sécurisé ;
- chiffrement des secrets ;
- rotation des clés ;
- révocation ;
- supervision ;
- gestion des erreurs ;
- conformité.

---

## KMIP et HSM

Un KMS peut parfois s’appuyer sur un **HSM**.

Un **HSM** est un module matériel sécurisé utilisé pour protéger les clés cryptographiques.

Schéma possible :

```txt
Client KMIP
   ↓
Serveur KMS
   ↓
HSM
```

Dans ce cas :

- KMIP sert à communiquer avec le KMS ;
- le KMS orchestre les opérations ;
- le HSM protège réellement les clés.

---

## KMIP et PKCS#11

Il ne faut pas confondre KMIP et PKCS#11.

| Élément | Rôle |
|---|---|
| KMIP | Protocole réseau pour gérer des objets cryptographiques à distance |
| PKCS#11 | API locale souvent utilisée pour interagir avec un HSM ou un token cryptographique |

Formulation simple :

```txt
KMIP = parler à un serveur de gestion de clés
PKCS#11 = utiliser une interface cryptographique locale
```

---

## Exemple d’architecture simple

```txt
Interface Web
    ↓
API Backend
    ↓
Décodeur TTLV
    ↓
Dispatcher KMIP
    ↓
KMS Core
    ↓
Repository sécurisé
    ↓
Audit + métriques
```

---

## Exemple dans un prototype Go

Dans un prototype en Go, on peut découper le projet ainsi :

```txt
ttlv/
  decoder.go
  types.go
  check.go

kms/
  key.go
  repository.go

usecase/
  dispatcher.go
  handlers/

audit/
  logger.go

metrics/
  collector.go

transport/httpapi/
  handler.go
```

Flux :

```txt
POST /kmip
  ↓
DecodeBlocks()
  ↓
blocksToOperationRequest()
  ↓
dispatcher.Dispatch()
  ↓
CreateKey / GetKey / DestroyKey
  ↓
Audit + metrics
```

---

## Ce qu’un prototype peut simuler

Un prototype peut simuler :

- le décodage TTLV ;
- les opérations `Create`, `Get`, `Destroy` ;
- un repository mémoire ;
- des statuts de clés ;
- des logs d’audit ;
- des métriques ;
- une interface web de démonstration.

Mais il ne doit pas prétendre être un vrai serveur KMIP complet.

---

## Limites d’un prototype pédagogique

Un prototype simple ne gère généralement pas encore :

- la norme KMIP complète ;
- la structure complète `RequestMessage`, `RequestHeader`, `BatchItem`, `RequestPayload` ;
- l’encodage TTLV des réponses ;
- TLS/mTLS ;
- l’authentification ;
- l’autorisation ;
- le stockage sécurisé ;
- les vraies clés cryptographiques ;
- les politiques d’accès ;
- la rotation de clés ;
- l’intégration HSM.

---

## Formulation courte pour entretien

> KMIP est un protocole standard qui permet à différents clients de communiquer avec un KMS pour gérer des objets cryptographiques.  
> Le KMS est le service qui gère le cycle de vie des clés, tandis que KMIP définit le langage et le format des messages utilisés pour demander des opérations comme Create, Get, Destroy, Activate ou Revoke.  
> Les messages KMIP sont encodés dans un format binaire appelé TTLV, pour Tag, Type, Length, Value.

---

## Formulation très simple

```txt
KMS = le coffre-fort qui gère les clés
KMIP = le langage pour parler au coffre-fort
TTLV = le format binaire des messages
```

---

## Résumé

KMIP permet de standardiser la gestion des objets cryptographiques.

Il aide à construire des systèmes plus interopérables, plus propres et plus faciles à intégrer dans des environnements complexes.

Dans une architecture moderne, KMIP peut servir de pont entre :

- des applications clientes ;
- un KMS ;
- un HSM ;
- des outils de sécurité ;
- des services cloud ;
- des systèmes de supervision.
