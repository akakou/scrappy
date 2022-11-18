package core

import (
	"fmt"
	"testing"
)

func TestCrypto(t *testing.T) {
	origin := "aaa"
	period := 1000

	err := CryptoSetup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	basename := fmt.Sprintf("%v_%v", origin, period)

	signatuere, err := CryptoSign(basename)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = CryptoVerify(signatuere, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}
}
