# KMIP

## Definition

KMIP signifie :

```txt
Key Management Interoperability Protocol
```

C'est un protocole standard qui permet a des clients de communiquer avec un serveur de gestion de cles.

Le but de KMIP est l'interoperabilite :

```txt
clients differents
  |
  | protocole commun
  v
KMS
```

Sans KMIP, chaque fournisseur pourrait exposer sa propre API. Avec KMIP, un client peut parler a plusieurs KMS compatibles.

## Ce que KMIP definit

KMIP definit notamment :

- les operations possibles ;
- les objets cryptographiques ;
- les attributs ;
- le format des messages ;
- les resultats ;
- les erreurs ;
- le cycle de vie des objets ;
- l'encodage TTLV.

## Operations courantes

| Operation | Role |
|---|---|
| `Create` | Creer un objet cryptographique |
| `Get` | Recuperer un objet ou ses informations |
| `Destroy` | Detruire un objet |
| `Activate` | Activer un objet |
| `Revoke` | Revoquer un objet |
| `Locate` | Rechercher des objets |
| `Query` | Demander les capacites du serveur |
| `Register` | Enregistrer un objet existant |
| `GetAttributes` | Lire les attributs |

Le projet implemente actuellement :

- `Create` ;
- `Get` ;
- `Destroy` ;
- `Activate` ;
- `Revoke` ;
- `Locate` ;
- `GetAttributes`.

## Objets cryptographiques

KMIP peut representer plusieurs types d'objets :

- cle symetrique ;
- cle publique ;
- cle privee ;
- certificat ;
- secret data ;
- opaque object ;
- split key ;
- template.

Dans le projet, l'interface web construit surtout des requetes pour :

```txt
ObjectType = SymmetricKey
```

## Structure logique d'une requete KMIP

Dans KMIP complet, une requete ressemble plutot a ceci :

```txt
RequestMessage
  RequestHeader
    ProtocolVersion
    BatchCount
  BatchItem
    Operation
    RequestPayload
      ObjectType
      Attributes
```

Le projet simplifie ce format.

Au lieu de decoder toute la structure KMIP, il accepte directement une suite de blocs TTLV utiles :

```txt
Operation
ObjectType ou UniqueIdentifier
RequestPayload optionnel
```

## Flux dans le projet

### Create

```txt
Client
  |
  | POST /kmip
  | Operation = Create
  | ObjectType = SymmetricKey
  v
transport/httpapi
  |
  v
ttlv.DecodeBlocks
  |
  v
kmipwire.BlocksToOperationRequest
  |
  v
dispatcher.Dispatch
  |
  v
CreateKey
  |
  v
repository.Create
  |
  v
audit + metrics
  |
  v
JSON { key_id, status }
```

### Get

```txt
Operation = Get
UniqueIdentifier = key-id
```

Le serveur cherche la cle. Si elle n'existe pas ou si elle est detruite, il retourne une erreur `not_found`.

### Destroy

```txt
Operation = Destroy
UniqueIdentifier = key-id
```

Le serveur passe la cle a l'etat `destroyed`.

### Activate

```txt
Operation = Activate
UniqueIdentifier = key-id
```

Le serveur passe la cle a l'etat `active`, sauf si elle est deja detruite.

### Revoke

```txt
Operation = Revoke
UniqueIdentifier = key-id
```

Le serveur passe la cle a l'etat `revoked`, sauf si elle est deja detruite.

### Locate

```txt
Operation = Locate
```

Le serveur retourne les cles non detruites connues du repository.

### GetAttributes

```txt
Operation = GetAttributes
UniqueIdentifier = key-id
```

Le serveur retourne les attributs visibles de la cle : identifiant, type d'objet, etat, date de creation et date de mise a jour.

## Mapping KMIP vers Go

La couche de mapping se trouve dans :

```txt
internal/transport/kmipwire/mapper.go
```

Elle transforme :

```txt
[]ttlv.Block
```

en :

```go
models.OperationRequest
```

Exemple :

```go
type OperationRequest struct {
    Operation  ttlv.Operation
    KeyID      string
    ObjectType ttlv.ObjectType
    Payload    []byte
}
```

## Validation actuelle

Le projet verifie maintenant :

- operation presente et valide ;
- `Create` avec `ObjectType` valide ;
- `Get` avec `KeyID` non vide ;
- `Destroy` avec `KeyID` non vide ;
- `Activate` avec `KeyID` non vide ;
- `Revoke` avec `KeyID` non vide ;
- `GetAttributes` avec `KeyID` non vide ;
- `Locate` sans `KeyID` obligatoire ;
- tags et types TTLV connus ;
- longueur TTLV coherente.

## Reponses

Les reponses du prototype sont en JSON :

```json
{
  "key_id": "key-...",
  "status": "active"
}
```

Un vrai serveur KMIP renverrait une reponse encodee en TTLV, avec une structure `ResponseMessage`.

## Limites KMIP du projet

Le projet n'est pas encore un serveur KMIP complet :

- pas de `RequestHeader` complet ;
- pas de `BatchItem` complet ;
- pas de version de protocole ;
- pas de batch multiple ;
- pas d'encodage TTLV des reponses ;
- pas de result status KMIP complet ;
- pas de result reason complet ;
- pas de gestion complete des attributs ;
- pas de TLS/mTLS ;
- pas d'authentification client.

Ces limites sont acceptables pour un prototype d'apprentissage.

## Prochaines evolutions utiles

Pour rapprocher le projet de KMIP reel :

1. Ajouter `RequestMessage`, `RequestHeader`, `BatchItem`.
2. Encoder aussi les reponses en TTLV.
3. Ajouter `Activate`, `Revoke`, `Locate`, `Query`.
4. Ajouter les attributs KMIP.
5. Ajouter TLS/mTLS.
6. Ajouter des tests binaires avec messages KMIP plus proches de la norme.
