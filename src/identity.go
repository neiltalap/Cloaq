// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package network

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

type Identity struct {
	PrivateKey *ecdh.PrivateKey
	PublicKey  *ecdh.PublicKey
}

func (i *Identity) String() string {
	return hex.EncodeToString(i.PublicKey.Bytes())
}

func GenerateIdentity() (*Identity, error) {
	identity := &Identity{}
	pKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	identity.PrivateKey = pKey
	identity.PublicKey = pKey.Public().(*ecdh.PublicKey)

	return identity, nil
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
