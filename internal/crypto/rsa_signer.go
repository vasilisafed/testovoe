package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type RSASigner struct {
	privateKey *rsa.PrivateKey
}

func NewRSASigner(privateKey interface{}) (*RSASigner, error) {
	key, ok := privateKey.(*rsa.PrivateKey)
	if !ok || key == nil {
		return nil, fmt.Errorf("invalid RSA private key")
	}

	return &RSASigner{privateKey: key}, nil
}

func (s *RSASigner) Sign(data []byte) ([]byte, error) {
	sum := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, sum[:])
}
