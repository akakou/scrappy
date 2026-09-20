package scrappy_test

import (
	"path/filepath"
	"testing"

	"github.com/akakou/ecdaa"
	amclutils "github.com/akakou/fp256bn-amcl-utils"
	"github.com/akakou/scrappy"
)

func TestAll(t *testing.T) {
	scrappy.SIGNER_LOG_DB_PATH = "./test_log_sign.db"
	scrappy.VERIFIER_LOG_DB_PATH = "./test_log_verify.db"
	scrappy.VERIFIER_RL_DB_PATH = "./test_rl_verify.db"

	origin := "aaa"
	period := scrappy.Now()

	rng := amclutils.InitRandom()
	issuer, signer, err := ecdaa.ExampleInitialize(rng)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signature, err := scrappy.Sign(origin, period, signer)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	sweepDB(scrappy.SIGNER_LOG_DB_CONF, scrappy.SIGNER_LOG_DB_PATH, scrappy.HashBasename(origin))
	sweepDB(scrappy.VERIFIER_LOG_DB_CONF, scrappy.VERIFIER_LOG_DB_PATH, scrappy.GetBasename(origin, period))

	err = scrappy.Verify(signature, origin, period, &issuer.Ipk, &ecdaa.RevocationList{})

	if err != nil {
		t.Fatalf("%v: ", err)
	}
}

func TestFailBecauseOfMultiSignature(t *testing.T) {
	scrappy.SIGNER_LOG_DB_PATH = "./test_log_sign.db"
	scrappy.VERIFIER_LOG_DB_PATH = "./test_log_verify.db"
	scrappy.VERIFIER_RL_DB_PATH = "./test_rl_verify.db"

	origin := "aaa"
	period := scrappy.Now()

	rng := amclutils.InitRandom()
	issuer, signer, err := ecdaa.ExampleInitialize(rng)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signature, err := scrappy.Sign(origin, period, signer)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = scrappy.Verify(signature, origin, period, &issuer.Ipk, &ecdaa.RevocationList{})

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = scrappy.Verify(signature, origin, period, &issuer.Ipk, &ecdaa.RevocationList{})

	if err == nil {
		t.Fatalf("%v: ", err)
	}

	sweepDB(scrappy.SIGNER_LOG_DB_CONF, scrappy.SIGNER_LOG_DB_PATH, scrappy.HashBasename(origin))
	sweepDB(scrappy.VERIFIER_LOG_DB_CONF, scrappy.VERIFIER_LOG_DB_PATH, scrappy.GetBasename(origin, period))
}

func TestRateLimitK(t *testing.T) {
	oldK := scrappy.RateLimitK
	oldSignerLog := scrappy.SIGNER_LOG_DB_PATH
	oldVerifierLog := scrappy.VERIFIER_LOG_DB_PATH
	oldVerifierRL := scrappy.VERIFIER_RL_DB_PATH
	defer func() {
		scrappy.RateLimitK = oldK
		scrappy.SIGNER_LOG_DB_PATH = oldSignerLog
		scrappy.VERIFIER_LOG_DB_PATH = oldVerifierLog
		scrappy.VERIFIER_RL_DB_PATH = oldVerifierRL
	}()

	scrappy.RateLimitK = 2
	tmp := t.TempDir()
	scrappy.SIGNER_LOG_DB_PATH = filepath.Join(tmp, "signer.db")
	scrappy.VERIFIER_LOG_DB_PATH = filepath.Join(tmp, "verifier.db")
	scrappy.VERIFIER_RL_DB_PATH = filepath.Join(tmp, "rl.db")

	origin := "aaa"
	period := scrappy.Now()
	rng := amclutils.InitRandom()
	issuer, signer, err := ecdaa.ExampleInitialize(rng)
	if err != nil {
		t.Fatalf("%v", err)
	}

	proof1, err := scrappy.Sign(origin, period, signer)
	if err != nil {
		t.Fatalf("first proof: %v", err)
	}
	proof2, err := scrappy.Sign(origin, period, signer)
	if err != nil {
		t.Fatalf("second proof: %v", err)
	}
	if _, err := scrappy.Sign(origin, period, signer); err == nil {
		t.Fatal("third proof should exceed k=2")
	}

	if err := scrappy.Verify(proof1, origin, period, &issuer.Ipk, &ecdaa.RevocationList{}); err != nil {
		t.Fatalf("verify first proof: %v", err)
	}
	if err := scrappy.Verify(proof2, origin, period, &issuer.Ipk, &ecdaa.RevocationList{}); err != nil {
		t.Fatalf("verify second proof: %v", err)
	}
}
