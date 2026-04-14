package domain

type Signer interface {
	Sign(data []byte) ([]byte, error)
}
