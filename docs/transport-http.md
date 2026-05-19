# Protocole de transport HTTP

## Role du transport

Le transport est la couche qui expose le projet au monde exterieur.

Dans ce document, on parle du transport HTTP. Le projet expose aussi un transport TCP brut documenté dans [`transport-tcp.md`](./transport-tcp.md).

Son role :

- recevoir les requetes ;
- verifier la methode HTTP ;
- lire le body ;
- limiter la taille des entrees ;
- decoder le TTLV ;
- transformer les blocs en requete metier ;
- appeler le dispatcher ;
- retourner une reponse JSON ;
- exposer les endpoints de supervision.

## Module dans le projet

Le transport HTTP se trouve dans :

```txt
internal/transport/httpapi
```

Fichiers principaux :

| Fichier | Role |
|---|---|
| `handler.go` | Endpoint principal `POST /kmip` |
| `keys_handler.go` | Endpoint `GET /keys` |
| `metrics_handler.go` | Endpoint `GET /metrics` |
| `audit_handler.go` | Endpoint `GET /audit` |
| `dashboard.go` | Page HTML simple `/dashboard` |
| `handler_*_test.go` | Tests HTTP par comportement |

Le serveur est assemble dans :

```txt
cmd/api/main.go
```

## Endpoints disponibles

| Endpoint | Methode | Role |
|---|---|---|
| `/kmip` | `POST` | Executer une operation KMIP simplifiee |
| `/keys` | `GET` | Lister les cles connues du repository |
| `/metrics` | `GET` | Lire les compteurs en memoire |
| `/audit` | `GET` | Lire les evenements d'audit |
| `/dashboard` | `GET` | Page HTML de lien vers le dashboard React |
| `/health` | `GET` | Verification simple de sante |

## Endpoint POST /kmip

`POST /kmip` attend un body binaire encode en TTLV simplifie.

Exemple logique pour creer une cle :

```txt
Operation  = Create
ObjectType = SymmetricKey
```

Exemple logique pour lire une cle :

```txt
Operation        = Get
UniqueIdentifier = key-...
```

Exemple logique pour detruire une cle :

```txt
Operation        = Destroy
UniqueIdentifier = key-...
```

## Flux de traitement

```txt
HandleKMIP
  |
  | verifie POST
  | applique une limite de taille
  | lit le body
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
operation handler
  |
  v
writeJSON
```

## Limite de taille

Le body de `/kmip` est limite avec :

```go
http.MaxBytesReader
```

La limite est alignee sur :

```go
ttlv.MaxValueLength
```

Actuellement :

```txt
1 MiB
```

Si la requete depasse la limite, le serveur retourne :

```http
413 Request Entity Too Large
```

avec une reponse JSON :

```json
{
  "error": "request_too_large",
  "message": "request body too large"
}
```

## Validation des operations

Le transport valide maintenant les champs obligatoires :

| Operation | Champs requis |
|---|---|
| `Create` | `ObjectType` valide |
| `Get` | `UniqueIdentifier` non vide |
| `Destroy` | `UniqueIdentifier` non vide |
| `Activate` | `UniqueIdentifier` non vide |
| `Revoke` | `UniqueIdentifier` non vide |
| `Locate` | Aucun champ obligatoire |
| `GetAttributes` | `UniqueIdentifier` non vide |

Si la validation echoue, le serveur retourne :

```http
400 Bad Request
```

## Format des reponses

Les reponses du prototype sont en JSON.

Succes :

```json
{
  "key_id": "key-...",
  "status": "active"
}
```

Erreur :

```json
{
  "error": "bad_request",
  "message": "invalid or missing operation"
}
```

Cle introuvable :

```json
{
  "error": "not_found",
  "message": "kms: key not found"
}
```

## Endpoint GET /keys

`GET /keys` retourne la liste des cles connues :

```json
[
  {
    "id": "key-...",
    "created_at": "2026-05-19T...",
    "object_type": 2,
    "status": "active",
    "updated_at": "2026-05-19T..."
  }
]
```

Cet endpoint est utilise par le dashboard React pour afficher la page Keys et les statistiques du tableau de bord.

## Endpoint GET /metrics

`GET /metrics` retourne les compteurs :

```json
{
  "http_requests_total": 1,
  "http_errors_total": 0,
  "create_key_total": 1,
  "get_key_total": 0,
  "destroy_key_total": 0,
  "success_total": 1,
  "not_found_total": 0
}
```

## Endpoint GET /audit

`GET /audit` retourne les evenements :

```json
[
  {
    "time": "2026-05-19T...",
    "operation": "create_key",
    "key_id": "key-...",
    "status": "active",
    "result": "success"
  }
]
```

## Frontend et proxy Vite

Le frontend Vite utilise un proxy vers le backend Go.

Configuration :

```txt
web/vite.config.ts
```

Les routes proxifiees sont :

- `/kmip` ;
- `/keys` ;
- `/metrics` ;
- `/audit` ;
- `/health`.

## Limites du transport actuel

Le transport actuel ne gere pas encore :

- TLS ;
- mTLS ;
- authentification ;
- CORS explicite pour un deploiement separe ;
- rate limiting ;
- logs d'acces HTTP ;
- correlation ID ;
- timeout serveur configure ;
- reponses KMIP encodees en TTLV.

Pour un vrai serveur expose en reseau, ces points doivent etre ajoutes.
