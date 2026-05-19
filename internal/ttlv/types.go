// Package ttlv contient les types et constantes pour encoder/décoder
// des blocs TTLV proches du format KMIP.
//
// Format TTLV simplifié :
// Tag    = 3 bytes en KMIP réel, stocké ici dans uint32
// Type   = 1 byte
// Length = 4 bytes
// Value  = Length bytes
package ttlv

// Type représente le type TTLV KMIP.
type Type uint8

const (
	TypeStructure   Type = 0x01
	TypeInteger     Type = 0x02
	TypeLongInteger Type = 0x03
	TypeBigInteger  Type = 0x04
	TypeEnumeration Type = 0x05
	TypeBoolean     Type = 0x06
	TypeTextString  Type = 0x07
	TypeByteString  Type = 0x08
	TypeDateTime    Type = 0x09
	TypeInterval    Type = 0x0A
)

// Tag représente un tag KMIP.
// En KMIP réel, le tag est encodé sur 3 bytes.
// On le stocke dans uint32 pour simplifier le code Go.
type Tag uint32

const (
	TagRequestMessage         Tag = 0x420078
	TagRequestHeader          Tag = 0x420077
	TagBatchItem              Tag = 0x42000F
	TagOperation              Tag = 0x42005C
	TagRequestPayload         Tag = 0x420079
	TagResponseMessage        Tag = 0x42007B
	TagResponseHeader         Tag = 0x42007A
	TagResponsePayload        Tag = 0x42007C
	TagUniqueIdentifier       Tag = 0x420094
	TagObjectType             Tag = 0x420057
	TagCryptographicAlgorithm Tag = 0x420028
	TagCryptographicLength    Tag = 0x42002A
	TagName                   Tag = 0x420053
	TagErrorMessage           Tag = 0x42003D
)

// Operation représente une opération KMIP.
type Operation uint32

const (
	OperationCreate   Operation = 0x00000001
	OperationGet      Operation = 0x0000000A
	OperationLocate   Operation = 0x00000008
	OperationActivate Operation = 0x00000012
	OperationRevoke   Operation = 0x00000013
	OperationDestroy  Operation = 0x00000014
	OperationArchive  Operation = 0x00000015
	OperationCancel   Operation = 0x00000016
	OperationUpdate   Operation = 0x00000017
	OperationQuery    Operation = 0x00000018
	OperationNotify   Operation = 0x00000019
	OperationPoll     Operation = 0x0000001A
	OperationPong     Operation = 0x0000001B

	// OperationGetAttributes est une operation KMIP courante. La valeur est
	// celle de KMIP pour GetAttributes dans ce prototype.
	OperationGetAttributes Operation = 0x0000000B
)

// ObjectType représente un type d'objet KMIP.
type ObjectType uint32

const (
	ObjectTypeCertificate  ObjectType = 0x00000001
	ObjectTypeSymmetricKey ObjectType = 0x00000002
	ObjectTypePublicKey    ObjectType = 0x00000003
	ObjectTypePrivateKey   ObjectType = 0x00000004
	ObjectTypeSplitKey     ObjectType = 0x00000005
	ObjectTypeTemplate     ObjectType = 0x00000006
	ObjectTypeSecretData   ObjectType = 0x00000007
	ObjectTypeOpaqueObject ObjectType = 0x00000008
)

// Tailles des champs du header TTLV en bytes.
const (
	TagSize    = 3
	TypeSize   = 1
	LengthSize = 4

	HeaderSize     = TagSize + TypeSize + LengthSize
	MaxValueLength = 1 << 20 // 1 MiB
)

// Block représente un bloc TTLV : Tag, Type, Length, Value.
type Block struct {
	Tag    Tag
	Type   Type
	Length uint32
	Value  []byte
}
