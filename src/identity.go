// © 2026 Polypodomy — AGPLv3

package network

import (
	"cloaq/src/config"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

type Identity struct {
	PrivateKey *ecdh.PrivateKey
	PublicKey  *ecdh.PublicKey
}

func (i *Identity) String() string {
	return hex.EncodeToString(i.PublicKey.Bytes())
}

func identityPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".cloaq")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, "identity.key"), nil
}

func saveIdentity(path string, key []byte) error {
	return os.WriteFile(path, key, 0600)
}

func loadIdentity(path string) (*ecdh.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ecdh.X25519().NewPrivateKey(data)
}

func (i *Identity) Generate() error {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	i.PrivateKey = priv
	i.PublicKey = priv.Public().(*ecdh.PublicKey)

	return nil
}

func CreateOrLoadIdentity() (*Identity, error) {
	store, err := config.LoadStore()
	if err != nil {
		return nil, err
	}

	if len(store.Keys) > 0 {
		priv, err := ecdh.X25519().NewPrivateKey(store.Keys[0])
		if err != nil {
			return nil, err
		}

		return &Identity{
			PrivateKey: priv,
			PublicKey:  priv.Public().(*ecdh.PublicKey),
		}, nil
	}

	id := &Identity{}
	if err := id.Generate(); err != nil {
		return nil, err
	}

	// save identity
	store.Keys = append(store.Keys, id.PrivateKey.Bytes())
	if err := config.SaveStore(store); err != nil {
		return nil, err
	}

	return id, nil
}

// DO NOT REMOVE THIS
// the new change made identity creation revolve around checking for a yaml file
// this blocks me out of creating two identities on a single machine which i need for creating tests
func GenerateTestIdentity() (*Identity, error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Identity{
		PrivateKey: priv,
		PublicKey:  priv.Public().(*ecdh.PublicKey),
	}, nil
}

func ParsePublicKey(data []byte) (*ecdh.PublicKey, error) {
	curve := ecdh.X25519()

	pub, err := curve.NewPublicKey(data)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func (i *Identity) DeriveSharedKey(peerPub *ecdh.PublicKey) ([]byte, error) {
	// Perform ECDH key exchange
	secret, err := i.PrivateKey.ECDH(peerPub)
	if err != nil {
		return nil, err
	}
	// Hash the shared secret to derive symmetric key
	hash := sha256.Sum256(secret)
	return hash[:], nil
}
