package network

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"testing"
)

func TestDeriveSharedSecret(t *testing.T) {
	curve := ecdh.X25519()

	alicePriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate alice private key: %v", err)
	}

	bobPriv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate bob private key: %v", err)
	}

	aliceShared, err := DeriveSharedSecret(alicePriv.Bytes(), bobPriv.PublicKey().Bytes())
	if err != nil {
		t.Fatalf("alice derive failed: %v", err)
	}

	bobShared, err := DeriveSharedSecret(bobPriv.Bytes(), alicePriv.PublicKey().Bytes())
	if err != nil {
		t.Fatalf("bob derive failed: %v", err)
	}

	if !bytes.Equal(aliceShared, bobShared) {
		t.Fatal("shared secrets do not match")
	}

	if len(aliceShared) != 32 {
		t.Fatalf("expected 32-byte secret, got %d", len(aliceShared))
	}
}
