package core

import (
	"fmt"

	"github.com/akakou/ecdaa"
)

func Setup() error {
	return CryptoSetup()
}

func Sign(origin string, period int) (string, error) {
	db, err := SetupDB(SIGNER_DB_CONF, SIGNER_DB_PATH)

	if err != nil {
		return "", err
	}

	defer db.DB.Close()

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	basename := fmt.Sprintf("%v_%v", origin, period)

	hasExist, err := HasExist(db, basename)

	if err != nil {
		return "", fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return "", fmt.Errorf(HAS_EXIST_ERROR, basename)
	}

	signature, err := CryptoSign(basename)

	if err != nil {
		return "", err
	}

	err = Insert(db, basename)

	return signature, err

}

func Verify(signature, origin string, period int, rl ecdaa.RevocationList) error {
	K, err := GetK(signature)

	if err != nil {
		return err
	}

	db, err := SetupDB(VERIFIER_DB_CONF, VERIFIER_DB_PATH)

	if err != nil {
		return err
	}

	defer db.DB.Close()

	if !IsValidPeriod(period) {
		return fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	hasExist, err := HasExist(db, K)

	if err != nil {
		return fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return fmt.Errorf(HAS_EXIST_ERROR, signature)
	}

	err = CryptoVerify(signature, origin, period, rl)

	if err != nil {
		return err
	}

	err = Insert(db, K)

	return err
}
