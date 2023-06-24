package scrappy

import (
	"fmt"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/ecdaa_helper"
)

func Sign(origin string, period int, signer ecdaa.Signer) (string, error) {
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

	signature, err := ecdaa_helper.SignWithEncoding(basename, signer)

	if err != nil {
		return "", err
	}

	err = Insert(db, basename)

	return signature, err
}

func Verify(signatureString, origin string, period int, ipk *ecdaa.IPK, rl *ecdaa.RevocationList) error {
	signature, err := ecdaa_helper.DecodeSignature(signatureString)

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

	err = ecdaa.Verify([]byte{}, []byte(basename), signature, ipk, *rl)

	if err != nil {
		return err
	}

	err = Insert(logDB, K)

	return err
}
