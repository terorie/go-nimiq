package ed25519

// #cgo CFLAGS: -Inative
// #cgo LDFLAGS: -L${SRCDIR}/native -led25519
// #include <ed25519.h>
// #include <stdlib.h>
import "C"
import "unsafe"

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
	cSignature := unsafe.Pointer(signature)
	cMessage := unsafe.Pointer(&message[0])
	cMessageLen := C.size_t(len(message))
	cPubKey := unsafe.Pointer(publicKey)

	result := C.ed25519_verify(
		(*C.uchar)(cSignature),
		(*C.uchar)(cMessage),
		cMessageLen,
		(*C.uchar)(cPubKey),
	)

	switch result {
		case 0: return false
		case 1: return true
		default: panic("Undefined behaviour in ed25519_verify")
	}
}

// Single signature functions

func PublicKeyDerive(privateKey *PrivateKey) PublicKey {
	var outPublicKey PublicKey
	cPrivateKey := unsafe.Pointer(privateKey)
	cOutPublicKey := unsafe.Pointer(&outPublicKey)

	C.ed25519_public_key_derive(
		(*C.uchar)(cOutPublicKey),
		(*C.uchar)(cPrivateKey),
	)

	return outPublicKey
}

func Sign(message []byte, publicKey *PublicKey, privateKey *PrivateKey) Signature {
	var outSignature Signature
	cOutSignature := unsafe.Pointer(&outSignature)
	cMessage := unsafe.Pointer(&message[0])
	cMessageLen := C.size_t(len(message))
	cPubKey := unsafe.Pointer(publicKey)
	cPrivKey := unsafe.Pointer(privateKey)

	C.ed25519_sign(
		(*C.uchar)(cOutSignature),
		(*C.uchar)(cMessage),
		cMessageLen,
		(*C.uchar)(cPubKey),
		(*C.uchar)(cPrivKey),
	)

	return outSignature
}
