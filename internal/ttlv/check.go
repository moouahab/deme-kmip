package ttlv

import "fmt"

// IsValid vérifie si le type TTLV est connu.
func (t Type) IsValid() bool {
	switch t {
	case TypeStructure,
		TypeInteger,
		TypeLongInteger,
		TypeBigInteger,
		TypeEnumeration,
		TypeBoolean,
		TypeTextString,
		TypeByteString,
		TypeDateTime,
		TypeInterval:
		return true
	default:
		return false
	}
}

// IsValid vérifie si le tag TTLV est connu dans notre simulateur KMIP.
func (t Tag) IsValid() bool {
	switch t {
	case TagRequestMessage,
		TagRequestHeader,
		TagBatchItem,
		TagOperation,
		TagRequestPayload,
		TagResponseMessage,
		TagResponseHeader,
		TagResponsePayload,
		TagUniqueIdentifier,
		TagObjectType,
		TagCryptographicAlgorithm,
		TagCryptographicLength,
		TagName,
		TagErrorMessage:
		return true
	default:
		return false
	}
}

// IsValid vérifie si l'opération KMIP est connue.
func (op Operation) IsValid() bool {
	switch op {
	case OperationCreate,
		OperationGet,
		OperationDestroy,
		OperationLocate,
		OperationRevoke,
		OperationActivate,
		OperationGetAttributes,
		OperationCancel,
		OperationUpdate,
		OperationQuery,
		OperationNotify,
		OperationPoll,
		OperationPong,
		OperationArchive:
		return true
	default:
		return false
	}
}

// IsValid vérifie si le type d'objet KMIP est connu.
func (o ObjectType) IsValid() bool {
	switch o {
	case ObjectTypeCertificate,
		ObjectTypeSymmetricKey,
		ObjectTypePublicKey,
		ObjectTypePrivateKey,
		ObjectTypeSplitKey,
		ObjectTypeTemplate,
		ObjectTypeSecretData,
		ObjectTypeOpaqueObject:
		return true
	default:
		return false
	}
}

// Validate vérifie qu'un bloc TTLV est cohérent.
func (b Block) Validate() error {
	if !b.Tag.IsValid() {
		return fmt.Errorf("ttlv: unknown tag: 0x%06X", uint32(b.Tag))
	}
	if !b.Type.IsValid() {
		return fmt.Errorf("ttlv: unknown type: 0x%02X", uint8(b.Type))
	}
	if b.Length > MaxValueLength {
		return fmt.Errorf("ttlv: value too large: %d bytes", b.Length)
	}
	if uint32(len(b.Value)) != b.Length {
		return fmt.Errorf("ttlv: invalid length: expected %d bytes, got %d bytes", b.Length, len(b.Value))
	}
	return nil
}
