package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	keyPair *RSAKeyPair
}

type ECCSigner struct {
	keyPair *ECCKeyPair
}

// TODO: implement RSA and ECDSA signing ...
func NewRSASigner(keyPair *RSAKeyPair) *RSASigner {
	return &RSASigner{keyPair}
}

func NewECCSigner(keyPair *ECCKeyPair) *ECCSigner {
	return &ECCSigner{keyPair}
}

func (r *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	// Hash the data with SHA256 before signing
	hashed := sha256.Sum256(dataToBeSigned)
	signature, err := rsa.SignPSS(rand.Reader, r.keyPair.Private, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (e *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	signature, err := ecdsa.SignASN1(rand.Reader, e.keyPair.Private, dataToBeSigned)
	if err != nil {
		return nil, err
	}
	return signature, nil
}