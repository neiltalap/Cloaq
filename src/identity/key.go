package key

import (
	"crypto/rand"
)

type PrivateKey struct {
	Data [32]byte
}

func (privateKey *PrivateKey) Generate() error {
	_, err := rand.Read(privateKey.Data[:])
	if err != nil {
		return err
	}

	return nil
}
