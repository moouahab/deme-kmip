export const TagOperation = 0x42005c;
export const TagUniqueIdentifier = 0x420094;
export const TagObjectType = 0x420057;

export const TypeEnumeration = 0x05;
export const TypeTextString = 0x07;

export const OperationCreate = 0x00000001;
export const OperationGet = 0x0000000a;
export const OperationDestroy = 0x00000014;

export const ObjectTypeSymmetricKey = 0x00000002;

function uint32ToBytes(value: number): number[] {
  return [
    (value >>> 24) & 0xff,
    (value >>> 16) & 0xff,
    (value >>> 8) & 0xff,
    value & 0xff,
  ];
}

function tagToBytes(tag: number): number[] {
  return [(tag >>> 16) & 0xff, (tag >>> 8) & 0xff, tag & 0xff];
}

function textToBytes(value: string): number[] {
  return Array.from(new TextEncoder().encode(value));
}

function encodeBlock(tag: number, type: number, value: number[]): number[] {
  return [...tagToBytes(tag), type, ...uint32ToBytes(value.length), ...value];
}

function encodeEnumeration(tag: number, value: number): number[] {
  return encodeBlock(tag, TypeEnumeration, uint32ToBytes(value));
}

function encodeTextString(tag: number, value: string): number[] {
  return encodeBlock(tag, TypeTextString, textToBytes(value));
}

export function buildCreateKeyRequest(): Uint8Array {
  const bytes = [
    ...encodeEnumeration(TagOperation, OperationCreate),
    ...encodeEnumeration(TagObjectType, ObjectTypeSymmetricKey),
  ];

  return new Uint8Array(bytes);
}

export function buildGetKeyRequest(keyId: string): Uint8Array {
  const bytes = [
    ...encodeEnumeration(TagOperation, OperationGet),
    ...encodeTextString(TagUniqueIdentifier, keyId),
  ];

  return new Uint8Array(bytes);
}

export function buildDestroyKeyRequest(keyId: string): Uint8Array {
  const bytes = [
    ...encodeEnumeration(TagOperation, OperationDestroy),
    ...encodeTextString(TagUniqueIdentifier, keyId),
  ];

  return new Uint8Array(bytes);
}

export function toHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map((byte) => byte.toString(16).padStart(2, "0").toUpperCase())
    .join(" ");
}