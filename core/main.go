package core

import "fmt"

func Setup() error {
	return CryptoSetup()
}

func Sign(origin string, period int) (string, error) {
	db, err := SetupDB(SIGNER_DB_PATH)

	if err != nil {
		return "", err
	}

	defer db.Close()

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	signature, err := CryptoSign(origin, period)

	if err != nil {
		return "", err
	}

	K, err := GetK(signature)

	if err != nil {
		return "", err
	}

	err = InsertIfItHasNotExist(db, K)

	if err != nil {
		return "", err
	}

	return signature, nil

}

func Verify(signature, origin string, period int) error {
	db, err := SetupDB(VERIFIER_DB_PATH)

	if err != nil {
		return err
	}

	defer db.Close()

	if !IsValidPeriod(period) {
		return fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	err = CryptoVerify(signature, origin, period)

	if err != nil {
		return err
	}

	K, err := GetK(signature)

	if err != nil {
		return err
	}

	err = InsertIfItHasNotExist(db, K)

	return err
}
