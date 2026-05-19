# KMIP Demo

Prototype pédagogique d'un mini serveur KMS avec transports HTTP et TCP, basé sur des messages KMIP simplifiés encodés en TTLV.

Le projet sert à comprendre comment une requête binaire de type KMIP peut être reçue, décodée, transformée en opération métier, exécutée dans un KMS simulé, puis affichée dans une interface React.

## Fonctionnalités

- Décodage TTLV simplifié.
- Opérations KMIP simulées :
  - `Create`
  - `Get`
  - `Destroy`
  - `Activate`
  - `Revoke`
  - `Locate`
  - `GetAttributes`
- Repository KMS en mémoire.
- Audit log en mémoire.
- Métriques en mémoire.
- Dashboard React/Vite.
- Tests Go par module.
- CI GitHub Actions backend + frontend.

## Architecture

```txt
web/React ou client TCP
  |
  | HTTP POST /kmip
  | TCP localhost:5696
  v
cmd/api
  |
  v
internal/transport/httpapi ou tcpapi
  |
  v
internal/ttlv
  |
  v
internal/usecase
  |
  v
internal/kms
  |
  +--> internal/audit
  +--> internal/metrics
```

## Modules principaux

| Module | Rôle |
|---|---|
| `cmd/api` | Point d'entrée du serveur Go |
| `internal/transport/httpapi` | Endpoints HTTP et mapping des requêtes |
| `internal/transport/tcpapi` | Transport TCP brut pour messages TTLV |
| `internal/transport/kmipwire` | Mapping partagé TTLV vers opération métier |
| `internal/ttlv` | Types, constantes et décodage TTLV |
| `internal/usecase` | Dispatcher et handlers métier KMIP |
| `internal/kms` | Modèle de clé et repository mémoire |
| `internal/audit` | Journalisation des opérations |
| `internal/metrics` | Compteurs en mémoire |
| `web` | Interface React/Vite |
| `docs` | Documentation pédagogique KMIP/KMS/TTLV/PKCS#11 |

## Endpoints

| Endpoint | Méthode | Description |
|---|---|---|
| `/kmip` | `POST` | Exécute une opération KMIP simplifiée |
| `/keys` | `GET` | Liste les clés du repository |
| `/metrics` | `GET` | Retourne les métriques |
| `/audit` | `GET` | Retourne les événements d'audit |
| `/health` | `GET` | Vérifie l'état du backend |
| `/dashboard` | `GET` | Page HTML simple de redirection |

## Transport TCP

Le projet expose aussi un transport TCP brut sur :

```txt
localhost:5696
```

Le client envoie directement un message TTLV. Le serveur répond par une ligne JSON.

Ce transport est plus proche de l'idée d'un serveur KMIP réseau que l'API HTTP utilisée par le dashboard, mais il reste volontairement simplifié.

## Lancer le backend

```bash
go run ./cmd/api
```

Le backend HTTP écoute sur :

```txt
http://localhost:8080
```

Le transport TCP écoute sur :

```txt
localhost:5696
```

## Lancer le frontend

```bash
cd web
npm install
npm run dev
```

Le frontend Vite écoute généralement sur :

```txt
http://localhost:5173
```

Le proxy Vite redirige les appels `/kmip`, `/keys`, `/metrics`, `/audit` et `/health` vers le backend Go.

## Tests et vérifications

Tests Go :

```bash
go test ./...
```

Si le cache Go local est en lecture seule dans l'environnement courant :

```bash
GOCACHE=/tmp/kmipdemo-go-cache go test ./...
```

Frontend :

```bash
cd web
npm run lint
npm run build
```

## Documentation

La documentation pédagogique est dans [`docs/index.md`](./docs/index.md).

Elle explique :

- KMIP ;
- KMS ;
- TTLV ;
- PKCS#11 ;
- le transport HTTP ;
- le transport TCP ;
- le rôle des modules Go du projet.

## Limites

Ce projet est un laboratoire d'apprentissage, pas un KMS de production.

Il ne gère pas encore :

- TLS/mTLS ;
- authentification ;
- autorisation ;
- persistance ;
- vraie matière cryptographique ;
- intégration HSM ;
- encodage complet des réponses KMIP en TTLV.

Ces limites sont volontaires pour garder le code lisible et progressif.
