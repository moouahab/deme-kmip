# TTLV

## Definition

TTLV signifie :

```txt
Tag - Type - Length - Value
```

C'est le format binaire utilise par KMIP pour representer les champs d'un message.

Chaque bloc TTLV contient quatre parties :

| Partie | Role |
|---|---|
| `Tag` | Identifie le champ, par exemple `Operation` ou `UniqueIdentifier` |
| `Type` | Indique le type de la valeur, par exemple `Enumeration` ou `TextString` |
| `Length` | Indique la taille de la valeur en octets |
| `Value` | Contient la valeur brute |

## Exemple simple

Une operation `Create` peut etre representee comme ceci :

```txt
Tag    = Operation
Type   = Enumeration
Length = 4
Value  = 0x00000001
```

Dans le projet, cela correspond a :

```txt
TagOperation    = 0x42005C
TypeEnumeration = 0x05
OperationCreate = 0x00000001
```

## Encodage simplifie du projet

Le projet utilise un TTLV volontairement simplifie :

```txt
Tag    = 3 bytes
Type   = 1 byte
Length = 4 bytes, big endian
Value  = Length bytes
```

La taille du header est donc :

```txt
3 + 1 + 4 = 8 bytes
```

Exemple de bloc `Operation = Create` :

```txt
42 00 5C 05 00 00 00 04 00 00 00 01
```

Lecture :

| Octets | Signification |
|---|---|
| `42 00 5C` | Tag `Operation` |
| `05` | Type `Enumeration` |
| `00 00 00 04` | Longueur 4 |
| `00 00 00 01` | Operation `Create` |

## Module dans le projet

Le module TTLV se trouve dans :

```txt
internal/ttlv
```

Fichiers principaux :

| Fichier | Role |
|---|---|
| `types.go` | Declare les types, tags, operations et constantes |
| `decoder.go` | Decode une suite d'octets en blocs TTLV |
| `check.go` | Valide les tags, types, operations et blocs |
| `decoder_test.go` | Teste le decodage TTLV |

## Fonctionnement du decodeur

Le decodeur lit le message bloc par bloc :

```txt
DecodeBlocks(data)
  |
  v
decodeBlock(data[offset:])
  |
  | lit Tag, Type, Length
  | verifie la taille
  | extrait Value
  | valide le bloc
  v
[]Block
```

La structure Go produite est :

```go
type Block struct {
    Tag    Tag
    Type   Type
    Length uint32
    Value  []byte
}
```

## Validation

Le module verifie :

- que le message contient au moins un header complet ;
- que le tag est connu ;
- que le type est connu ;
- que la longueur annoncee correspond a la taille de `Value` ;
- que la valeur ne depasse pas `MaxValueLength`.

La constante importante est :

```go
MaxValueLength = 1 << 20 // 1 MiB
```

## Role dans le flux KMIP

Le module TTLV ne connait pas le metier KMS. Il ne cree pas de cle et ne choisit pas l'operation a executer.

Son role est uniquement de transformer :

```txt
bytes bruts
```

en :

```txt
blocs TTLV valides
```

Ensuite, le module transport transforme ces blocs en `OperationRequest`.

## Limites du prototype

Le projet ne gere pas encore tout TTLV/KMIP :

- pas de padding KMIP complet ;
- pas de structures imbriquees completes ;
- pas d'encodage TTLV des reponses ;
- pas de `RequestMessage` complet ;
- pas de `ResponseMessage` complet ;
- pas de version KMIP dans l'en-tete.

Ces choix simplifient le projet pour se concentrer sur le principe du protocole.
