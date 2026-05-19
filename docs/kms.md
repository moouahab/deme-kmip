# KMS

## Definition

KMS signifie :

```txt
Key Management Service
```

Un KMS est un service responsable de la gestion des objets cryptographiques :

- creation de cles ;
- stockage ;
- lecture ;
- destruction ;
- cycle de vie ;
- audit ;
- controle d'acces ;
- politiques de securite.

Dans ce projet, le KMS est volontairement simple : il gere des cles en memoire, sans stocker de vraie matiere secrete.

## Difference entre KMS et KMIP

Il faut separer les deux notions :

| Element | Role |
|---|---|
| KMS | Le service qui gere les cles |
| KMIP | Le protocole utilise pour parler au KMS |

Formulation simple :

```txt
KMS  = le coffre-fort
KMIP = le langage pour parler au coffre-fort
```

## Module dans le projet

Le module KMS se trouve dans :

```txt
internal/kms
```

Fichiers principaux :

| Fichier | Role |
|---|---|
| `key.go` | Definit la structure d'une cle et ses statuts |
| `repository.go` | Definit l'interface de stockage |
| `memory_repository.go` | Implante un repository en memoire |
| `memory_repository_test.go` | Teste le repository |

## Modele de cle

Une cle est representee par :

```go
type Key struct {
    ID         string
    CreatedAt  time.Time
    ObjectType ttlv.ObjectType
    Status     KeyStatus
    UpdatedAt  time.Time
}
```

Le projet ne stocke pas de vraie valeur de cle.

C'est un point important : le prototype simule la gestion d'une cle, mais il ne manipule pas de secret cryptographique.

## Statuts de cle

Les statuts definis sont :

| Statut | Signification |
|---|---|
| `pre_active` | La cle existe mais n'est pas encore active |
| `active` | La cle est utilisable |
| `revoked` | La cle est revoquee |
| `destroyed` | La cle est detruite logiquement |

Actuellement, le flux principal est simplifie :

```txt
Create
  |
  v
active
  |
  v
Destroy
  |
  v
destroyed
```

## Repository

L'interface de repository est :

```go
type Repository interface {
    Create(ctx context.Context, key Key) (Key, error)
    Get(ctx context.Context, id string) (Key, error)
    Update(ctx context.Context, key Key) (Key, error)
    List(ctx context.Context) ([]Key, error)
}
```

Cette interface permet de remplacer plus tard le stockage memoire par :

- SQLite ;
- PostgreSQL ;
- BoltDB ;
- un HSM ;
- un service cloud ;
- un repository chiffre.

## Repository memoire

`MemoryRepository` utilise :

```go
map[string]Key
```

protege par :

```go
sync.RWMutex
```

Cela permet d'eviter les races quand plusieurs requetes arrivent en meme temps.

## Comportement des operations

### Create

Le handler `CreateKey` :

- genere un ID avec UUID ;
- cree une cle active ;
- l'enregistre dans le repository ;
- ecrit un evenement d'audit ;
- retourne `key_id` et `status`.

### Get

Le handler `GetKey` :

- cherche la cle par ID ;
- retourne une erreur si elle n'existe pas ;
- cache les cles detruites en les traitant comme introuvables ;
- ecrit un evenement d'audit.

### Destroy

Le handler `DestroyKey` :

- cherche la cle par ID ;
- refuse une cle deja detruite ;
- passe son statut a `destroyed` ;
- met a jour le repository ;
- ecrit un evenement d'audit.

### Activate

Le handler `ActivateKey` :

- cherche la cle par ID ;
- refuse une cle detruite ;
- passe son statut a `active` ;
- ecrit un evenement d'audit.

### Revoke

Le handler `RevokeKey` :

- cherche la cle par ID ;
- refuse une cle detruite ;
- passe son statut a `revoked` ;
- ecrit un evenement d'audit.

### Locate

Le handler `LocateKeys` :

- liste les cles du repository ;
- masque les cles detruites ;
- retourne un resume des cles visibles.

### GetAttributes

Le handler `GetKeyAttributes` :

- cherche la cle par ID ;
- refuse une cle detruite ;
- retourne les attributs visibles : identifiant, type, etat, dates ;
- ecrit un evenement d'audit.

## Audit et metriques

Le KMS ne travaille pas seul.

Chaque operation importante produit :

- un evenement d'audit dans `internal/audit` ;
- un compteur dans `internal/metrics`, via la couche HTTP.

Cela permet au frontend d'afficher :

- le nombre de requetes ;
- le nombre de creations ;
- le nombre de lectures ;
- le nombre de destructions ;
- les erreurs ;
- les cles introuvables.

## Limites du KMS actuel

Le module actuel ne gere pas encore :

- persistance ;
- chiffrement au repos ;
- authentification ;
- autorisation par utilisateur ou role ;
- rotation de cles ;
- expiration ;
- activation separee ;
- revocation reelle ;
- generation cryptographique ;
- stockage de matiere secrete ;
- integration HSM.

Pour un vrai KMS, ces points deviennent obligatoires.
