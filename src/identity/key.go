package key

import (
	"crypto/rand"
	"errors"
	"os"
)

var WrongByteSize = errors.New("size of the key is not 32 bytes")

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

func (privateKey *PrivateKey) Save(filePath string) error {
	if err := os.WriteFile(filePath, privateKey.Data[:], 0600); err != nil {
		return err
	}

	return nil
}

func Load(filePath string) (*PrivateKey, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(data) != 32 {
		return nil, WrongByteSize
	}

	tempPrivateKey := new(PrivateKey)
	// copying the data read from the file to the instance
	copy(tempPrivateKey.Data[:], data)

	return tempPrivateKey, nil
}
