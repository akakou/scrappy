package core

import "testing"

func TestCrypto(t *testing.T) {
	origin := "aaa"
	period := 1000

	err := CryptoSetup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signatuere, err := CryptoSign(origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = CryptoVerify(signatuere, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

}
