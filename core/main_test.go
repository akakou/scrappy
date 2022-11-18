package core

import (
	"testing"
)

func TestAll(t *testing.T) {
	SIGNER_DB_PATH = "./test_sign.db"
	VERIFIER_DB_PATH = "./test_verify.db"

	origin := "aaa"
	period := Now()

	err := Setup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signature, err := Sign(origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = Verify(signature, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}
}

func TestFailBecauseOfMultiSignature(t *testing.T) {
	SIGNER_DB_PATH = "./test_sign.db"
	VERIFIER_DB_PATH = "./test_verify.db"

	origin := "aaa"
	period := Now()

	err := Setup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signature, err := Sign(origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = Verify(signature, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = Verify(signature, origin, period)

	if err == nil {
		t.Fatalf("%v: ", err)
	}
}
