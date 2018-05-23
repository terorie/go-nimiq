package ed25519

// #cgo CFLAGS: -Inative
// #cgo LDFLAGS: -L${SRCDIR}/native -led25519
// #include <ed25519.h>
// #include <stdlib.h>
import "C"

const (
	SeedSize = uintptr(32)
	SignatureSize = uintptr(64)
	PublicKeySize = uintptr(32)
	PrivateKeySize = uintptr(64)
	ScalarSize = uintptr(32)
	SharedSecretSize = uintptr(32)
)

type Seed [SeedSize]byte
type Signature [SignatureSize]byte
type PublicKey [PublicKeySize]byte
type PrivateKey [PrivateKeySize]byte
type Scalar [ScalarSize]byte
type SharedSecret [SharedSecretSize]byte

// Common functions

// func PrivateKeyDecompress(…)

func Verify(signature *Signature, message []byte, publicKey *PublicKey) bool {
	cSignature := C.CBytes(signature[:])
	cMessage := C.CBytes(message)
	cMessageLen := C.size_t(len(message))
	cPubKey := C.CBytes(publicKey[:])

	result := C.ed25519_verify(
		(*C.uchar)(cSignature),
		(*C.uchar)(cMessage),
		cMessageLen,
		(*C.uchar)(cPubKey),
	)

	C.free(cSignature)
	C.free(cMessage)
	C.free(cPubKey)

	switch result {
		case 0: return false
		case 1: return true
		default: panic("Undefined behaviour in ed25519_verify")
	}
}

// Single signature functions

func PublicKeyDerive(privateKey *PrivateKey) PublicKey {
	cPrivateKey := C.CBytes(privateKey[:])
	cOutPublicKey := C.malloc(C.size_t(PublicKeySize))

	C.ed25519_public_key_derive(
		(*C.uchar)(cOutPublicKey),
		(*C.uchar)(cPrivateKey),
	)

	var goPublicKey PublicKey
	pubKeyBytes := C.GoBytes(cOutPublicKey, C.int(PublicKeySize))
	copy(goPublicKey[:], pubKeyBytes[:PublicKeySize])

	C.free(cPrivateKey)
	C.free(cOutPublicKey)

	return goPublicKey
}

func Sign(message []byte, publicKey *PublicKey, privateKey *PrivateKey) Signature {
	cOutSignature := C.malloc(C.size_t(SignatureSize))
	cMessage := C.CBytes(message)
	cMessageLen := C.size_t(len(message))
	cPubKey := C.CBytes(publicKey[:])
	cPrivKey := C.CBytes(privateKey[:])

	C.ed25519_sign(
		(*C.uchar)(cOutSignature),
		(*C.uchar)(cMessage),
		cMessageLen,
		(*C.uchar)(cPubKey),
		(*C.uchar)(cPrivKey),
	)

	var goSignature Signature
	signatureBytes := C.GoBytes(cOutSignature, C.int(SignatureSize))
	copy(goSignature[:], signatureBytes[:SignatureSize])

	C.free(cOutSignature)
	C.free(cMessage)
	C.free(cPubKey)
	C.free(cPrivKey)

	return goSignature
}
