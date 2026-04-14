package domain

type Algorithm string

const (
	RSA Algorithm = "RSA"
	ECC Algorithm = "ECC"
)

type Device struct {
	ID               string
	Label            string
	Algorithm        Algorithm
	PrivateKey       interface{}
	PublicKey        interface{}
	SignatureCounter uint64
	LastSignature    string
}
