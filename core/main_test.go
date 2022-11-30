package core

import (
	"testing"
)

func TestAll(t *testing.T) {
	SIGNER_LOG_DB_PATH = "./test_log_sign.db"
	VERIFIER_LOG_DB_PATH = "./test_log_verify.db"
	VERIFIER_RL_DB_PATH = "./test_rl_verify.db"

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

	sweepDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, getBasename(origin, period))
	sweepDB(VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH, getBasename(origin, period))

	err = Verify(signature, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}
}

func TestFailBecauseOfMultiSignature(t *testing.T) {
	SIGNER_LOG_DB_PATH = "./test_log_sign.db"
	VERIFIER_LOG_DB_PATH = "./test_log_verify.db"
	VERIFIER_RL_DB_PATH = "./test_rl_verify.db"

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

	sweepDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, getBasename(origin, period))
	sweepDB(VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH, getBasename(origin, period))
}
