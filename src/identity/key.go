package key

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"os"
)

var WrongByteSize = errors.New("size of the key is not 32 bytes")

// ed25519 based cryptographic private key
type Key struct {
	Seed       [32]byte
	PrivateKey ed25519.PrivateKey
}

// generates a new ed25519 key and a seed along side it
func (privateKey *Key) Generate() error {
	_, err := rand.Read(privateKey.Seed[:])
	if err != nil {
		return err
	}

	privateKey.PrivateKey = ed25519.NewKeyFromSeed(privateKey.Seed[:])

	return nil
}

func (key *Key) Save(filePath string) error {
	if err := os.WriteFile(filePath, key.Seed[:], 0600); err != nil {
		return err
	}

	return nil
}

// rederives the private ed25519 key on load of the stored hash
func Load(filePath string) (*Key, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(data) != 32 {
		return nil, WrongByteSize
	}

	tempPrivateKey := new(Key)
	// copying the data read from the file to the instance
	copy(tempPrivateKey.Seed[:], data)
	// re-generating the ed25519 key from the seed
	tempPrivateKey.PrivateKey = ed25519.NewKeyFromSeed(data)

	return tempPrivateKey, nil
}
