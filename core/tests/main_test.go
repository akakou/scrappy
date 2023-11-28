package scrappy_test

import (
	"testing"

	"github.com/akakou/ecdaa"
	"github.com/akakou/mcl_utils"
	"github.com/akakou/scrappy"
)

func TestAll(t *testing.T) {
	scrappy.SIGNER_LOG_DB_PATH = "./test_log_sign.db"
	scrappy.VERIFIER_LOG_DB_PATH = "./test_log_verify.db"
	scrappy.VERIFIER_RL_DB_PATH = "./test_rl_verify.db"

	origin := "aaa"
	period := scrappy.Now()

	rng := mcl_utils.InitRandom()
	issuer, signer, err := ecdaa.ExampleInitialize(rng)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signature, err := scrappy.Sign(origin, period, signer)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	sweepDB(scrappy.SIGNER_LOG_DB_CONF, scrappy.SIGNER_LOG_DB_PATH, scrappy.GetBasename(origin, period))
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

	rng := mcl_utils.InitRandom()
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

	sweepDB(scrappy.SIGNER_LOG_DB_CONF, scrappy.SIGNER_LOG_DB_PATH, scrappy.GetBasename(origin, period))
	sweepDB(scrappy.VERIFIER_LOG_DB_CONF, scrappy.VERIFIER_LOG_DB_PATH, scrappy.GetBasename(origin, period))
}
