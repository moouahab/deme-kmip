# Transport TCP

## Rôle

Le transport TCP permet d'envoyer un message TTLV brut directement au serveur, sans passer par HTTP.

Dans un vrai contexte KMIP, les échanges réseau se font généralement sur TCP avec TLS ou mTLS. Ce projet ne va pas jusque-là, mais il ajoute une couche TCP simple pour montrer qu'une même logique métier peut être exposée par plusieurs transports.

## Adresse

Le serveur TCP écoute sur :

```txt
localhost:5696
```

Le port `5696` est le port habituellement associé à KMIP.

## Module Go

Le code se trouve dans :

```txt
internal/transport/tcpapi
```

Fichiers principaux :

| Fichier | Rôle |
|---|---|
| `server.go` | Serveur TCP et traitement d'un message TTLV |
| `server_test.go` | Tests du transport TCP |

## Flux

```txt
client TCP
  |
  | bytes TTLV
  v
tcpapi.Server
  |
  v
ttlv.DecodeBlocks
  |
  v
kmipwire.BlocksToOperationRequest
  |
  v
usecase.Dispatcher
  |
  v
handler métier KMS
  |
  v
réponse JSON
```

## Réponse

Le transport TCP retourne une ligne JSON.

Succès :

```json
{
  "ok": true,
  "data": {
    "key_id": "key-...",
    "status": "active"
  }
}
```

Erreur :

```json
{
  "ok": false,
  "error": "bad_request",
  "message": "invalid or missing operation"
}
```

## Pourquoi un package kmipwire

Le mapping des blocs TTLV vers une opération métier est partagé entre HTTP et TCP.

Il est donc placé dans :

```txt
internal/transport/kmipwire
```

Cela évite que `httpapi` et `tcpapi` dupliquent la même logique de validation.

## Limites

Le transport TCP actuel reste simple :

- pas de TLS ;
- pas de mTLS ;
- pas d'authentification ;
- pas de framing KMIP complet ;
- lecture d'un message par connexion ;
- réponse JSON au lieu d'une réponse TTLV KMIP complète.

Ces limites sont acceptables pour un prototype pédagogique.
