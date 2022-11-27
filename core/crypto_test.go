package core

import (
	"testing"

	"github.com/akakou/ecdaa"
)

func TestCrypto(t *testing.T) {
	origin := "aaa"
	period := 1000

	err := CryptoSetup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	basename := getBasename(origin, period)

	signatuere, err := CryptoSign(basename)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = CryptoVerify(signatuere, basename, ecdaa.RevocationList{})

	if err != nil {
		t.Fatalf("%v: ", err)
	}
}
