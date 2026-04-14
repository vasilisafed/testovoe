package crypto

import "crypto/ecdsa"
import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"fmt"
	"math/big"
)

type ECDSASigner struct {
	privateKey *ecdsa.PrivateKey
}

type ecdsaSignature struct {
	R *big.Int
	S *big.Int
}

func NewECDSASigner(privateKey interface{}) (*ECDSASigner, error) {
	key, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok || key == nil {
		return nil, fmt.Errorf("invalid ECDSA private key")
	}

	return &ECDSASigner{privateKey: key}, nil
}

func (s *ECDSASigner) Sign(data []byte) ([]byte, error) {
	sum := sha256.Sum256(data)
	r, ss, err := ecdsa.Sign(rand.Reader, s.privateKey, sum[:])
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ecdsaSignature{R: r, S: ss})
}
