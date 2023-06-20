package scrappy

import (
	"fmt"

	"github.com/akakou/scrappy/crypto"
)

func Sign(origin string, period int, config *crypto.SignerConfig) (string, error) {
	db, err := SetupDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH)

	if err != nil {
		return "", err
	}

	defer db.DB.Close()

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	basename := GetBasename(origin, period)

	err = CheckBasenameExists(basename, db)
	if err != nil {
		return "", err
	}

	signer, err := crypto.PrepareSWSigner(config)

	if err != nil {
		return "", err
	}

	signature, err := crypto.SignWithEncoding(basename, signer)

	if err != nil {
		return "", err
	}

	err = Insert(db, basename)

	return signature, err
}

func Verify(signatureString, origin string, period int) error {
	signature, err := crypto.DecodeSignature(signatureString)

	if err != nil {
		return err
	}

	basename := GetBasename(origin, period)

	K := GetKBytes(signature)

	logDB, err := SetupDB(VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH)

	if err != nil {
		return err
	}

	defer logDB.DB.Close()

	rlDB, err := SetupDB(VERIFIER_RL_DB_CONF, VERIFIER_RL_DB_PATH)

	if err != nil {
		return err
	}

	defer rlDB.DB.Close()

	if !IsValidPeriod(period) {
		return fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	hasExist, err := HasExist(logDB, K)

	if err != nil {
		return fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return fmt.Errorf(HAS_EXIST_ERROR, signature)
	}

	err = crypto.VerifyWithConfig(signature, []byte(basename))

	if err != nil {
		return err
	}

	err = Insert(logDB, K)

	return err
}
