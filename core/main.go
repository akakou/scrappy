package core

import "fmt"

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

	return signature, nil

}

func Verify(signature, origin string, period int) error {
	db, err := SetupDB(VERIFIER_DB_CONF, VERIFIER_DB_PATH)

	if err != nil {
		return err
	}

	defer db.DB.Close()

	if !IsValidPeriod(period) {
		return fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	hasExist, err := HasExist(db, signature)

	if err != nil {
		return fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return fmt.Errorf(HAS_EXIST_ERROR, signature)
	}

	err = CryptoVerify(signature, origin, period)

	if err != nil {
		return err
	}

	K, err := GetK(signature)

	if err != nil {
		return err
	}

	err = Insert(db, K)

	return err
}
