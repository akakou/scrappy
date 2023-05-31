package scrappy

import (
	"testing"
	"time"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/crypto"
)

const WAIT_TIME = 2

func TestAndroidAndGin(t *testing.T) {
	now := Now()
	rng := ecdaa.InitRandom()
	issuer, err := crypto.SetupIssuer(rng)

	if err != nil {
		t.Fatal(err)
	}

	err = crypto.WriteConfig(issuer, crypto.ISSUER_CONFIG_PATH)

	if err != nil {
		t.Fatal(err)
	}

	verifier := crypto.VerifierConfig{
		IPK: issuer.IPK,
	}

	err = crypto.WriteConfig(&verifier, crypto.VERIFIER_CONFIG_PATH)

	if err != nil {
		t.Fatal(err)
	}

	go ServIssuer()

	time.Sleep(WAIT_TIME * time.Second)

	signer, err := androidJoin("http://127.0.0.1:8080", issuer.IPK)

	if err != nil {
		t.Fatal(err)
	}

	signature, err := AndroidSign("http://localhost:8080", now, signer)

	if err != nil {
		t.Fatal(err)
	}
	err = Verify(signature, "http://localhost:8080", now)

	if err != nil {
		t.Fatal(err)
	}
}
