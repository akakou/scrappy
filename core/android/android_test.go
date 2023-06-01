package android_scrappy

import (
	"testing"
	"time"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy"
	"github.com/akakou/scrappy/crypto"
	scrappy_gin "github.com/akakou/scrappy/gin"
	_ "github.com/mattn/go-sqlite3"
)

const WAIT_TIME = 2

func TestAndroidAndGin(t *testing.T) {
	now := scrappy.Now()
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

	go scrappy_gin.RunIssuer("../templates/*.html")

	time.Sleep(WAIT_TIME * time.Second)

	signer, err := androidJoin("http://127.0.0.1:8080", issuer.IPK)

	if err != nil {
		t.Fatal(err)
	}

	signature, err := androidSign("http://localhost:8080", now, signer)

	if err != nil {
		t.Fatal(err)
	}
	err = scrappy.Verify(signature, "http://localhost:8080", now)

	if err != nil {
		t.Fatal(err)
	}
}
