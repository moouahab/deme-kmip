# Documentation pédagogique KMIP Demo

Ce dossier contient la documentation pédagogique du projet.

La racine du dépôt explique le projet, comment le lancer et comment tester le code. Ici, le but est différent : expliquer les concepts KMIP/KMS/TTLV/PKCS#11 et les relier aux modules Go du prototype.

Le projet ne cherche pas encore à implémenter toute la norme KMIP, mais il montre les blocs essentiels d'un serveur de gestion de clés :

- réception d'une requête binaire ;
- décodage TTLV ;
- transformation en opération métier ;
- exécution dans un KMS simulé ;
- audit ;
- métriques ;
- visualisation via une interface web.

## Guides disponibles

| Document | Sujet |
|---|---|
| [TTLV](./ttlv.md) | Format binaire `Tag - Type - Length - Value` utilisé par KMIP |
| [KMS](./kms.md) | Rôle du Key Management Service et module `internal/kms` |
| [KMIP](./kmip.md) | Protocole de gestion de clés et flux du projet |
| [PKCS#11](./pkcs11.md) | Différence entre KMIP et PKCS#11, usage avec HSM |
| [Transport HTTP](./transport-http.md) | Endpoints, requêtes, réponses, validation et limites |
| [Transport TCP](./transport-tcp.md) | Transport brut TTLV sur TCP, plus proche de KMIP réseau |
| [Notes d'entretien](./entretien.md) | Support pour présenter le projet en entretien |

## Vue d'ensemble du flux

```txt
Navigateur / client TCP
  |
  | POST /kmip
  | TCP localhost:5696
  | body binaire TTLV
  v
transport/httpapi ou transport/tcpapi
  |
  | DecodeBlocks()
  v
ttlv
  |
  | OperationRequest
  v
usecase.Dispatcher
  |
  | Create / Get / Destroy
  v
kms.Repository
  |
  | audit + metrics
  v
JSON response
```

## Modules Go principaux

| Module | Chemin | Responsabilite |
|---|---|---|
| TTLV | `internal/ttlv` | Decoder et valider les blocs TTLV |
| KMS | `internal/kms` | Representer les cles et les stocker en memoire |
| Usecase | `internal/usecase` | Router les operations KMIP vers les handlers metier |
| Transport HTTP | `internal/transport/httpapi` | Exposer les endpoints HTTP |
| Transport TCP | `internal/transport/tcpapi` | Exposer un transport TCP brut |
| KMIP wire | `internal/transport/kmipwire` | Mapper les blocs TTLV vers le métier |
| Audit | `internal/audit` | Journaliser les operations |
| Metrics | `internal/metrics` | Compter les requetes et resultats |
| API | `cmd/api` | Assembler les modules et demarrer le serveur |

## Limites importantes

Ce depot est un prototype :

- pas de TLS/mTLS ;
- pas d'authentification ;
- pas d'autorisation ;
- pas de stockage persistant ;
- pas de vraie matiere cryptographique ;
- pas d'encodage KMIP complet des reponses ;
- pas d'integration HSM ou PKCS#11 reelle.

Ces limites sont volontaires pour garder le projet comprehensible.
