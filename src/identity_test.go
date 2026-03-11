package network

import (
	"bytes"
	"testing"
)

func TestSharedKeyDerivation(t *testing.T) {

	alice, err := GenerateIdentity()
	if err != nil {
		t.Fatal(err)
	}

	bob, err := GenerateIdentity()
	if err != nil {
		t.Fatal(err)
	}

	keyA, err := alice.DeriveSharedKey(bob.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	keyB, err := bob.DeriveSharedKey(alice.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(keyA, keyB) {
		t.Fatal("shared keys do not match")
	}
}
