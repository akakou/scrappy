package scrappy

import (
	"encoding/base64"
	"fmt"

	"github.com/akakou/scrappy/crypto"
)

func SignWithoutDB(origin string, period int, sk, cred, ipk string) (string, error) {
	SK, err := base64.StdEncoding.DecodeString(sk)
	if err != nil {
		return "", err
	}

	Cred, err := base64.StdEncoding.DecodeString(cred)
	if err != nil {
		return "", err
	}

	IPK, err := base64.StdEncoding.DecodeString(ipk)
	if err != nil {
		return "", err
	}

	var config = crypto.SignerConfig{
		IPK: IPK, Cred: Cred, SK: SK,
	}

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	basename := GetBasename(origin, period)

	signature, err := crypto.SWSign(basename, &config)

	return signature, err
}

func Sign(origin string, period int, skString, credString, ipkString string) (string, error) {
	ipk, err := base64.RawStdEncoding.DecodeString(skString)
	if err != nil {
		return "", err
	}

	cred, err := base64.RawStdEncoding.DecodeString(credString)
	if err != nil {
		return "", err
	}

	sk, err := base64.RawStdEncoding.DecodeString(ipkString)
	if err != nil {
		return "", err
	}

	config := crypto.SignerConfig{
		IPK:  ipk,
		Cred: cred,
		SK:   sk,
	}

	db, err := SetupDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH)

	if err != nil {
		return "", err
	}

	defer db.DB.Close()

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	basename := GetBasename(origin, period)

	hasExist, err := HasExist(db, basename)

	if err != nil {
		return "", fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return "", fmt.Errorf(HAS_EXIST_ERROR, basename)
	}

	signature, err := crypto.SWSign(basename, &config)

	if err != nil {
		return "", err
	}

	err = Insert(db, basename)

	return signature, err
}

func Verify(signatureString, origin string, period int) error {
	signature, err := decodeBase64(signatureString)

	if err != nil {
		return err
	}

	basename := GetBasename(origin, period)

	K, err := GetK(signature)

	if err != nil {
		return err
	}

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

	rl, err := SelectAllRL(rlDB)

	if err != nil {
		return fmt.Errorf("can't get RL: %v", err)
	}

	err = crypto.VerifyWithConfig(signatureString, []byte(basename), rl)

	if err != nil {
		return err
	}

	err = Insert(logDB, K)

	return err
}
