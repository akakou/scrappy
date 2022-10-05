package core

import "fmt"

func Setup() error {
	return CryptoSetup()
}

func Sign(origin string, period int) (string, error) {
	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	db, err := SetupDB(CLIENT_DB_PATH)

	if err != nil {
		return "", err
	}

	defer db.Close()

	err = InsertIfItDoesNotExist(db, origin, period)

	if err != nil {
		return "", err
	}

	return CryptoSign(origin, period)
}

func Verify(signature, origin string, period int) error {
	if !IsValidPeriod(period) {
		return fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	db, err := SetupDB(SERVER_DB_PATH)

	if err != nil {
		return err
	}

	defer db.Close()

	err = InsertIfItDoesNotExist(db, origin, period)

	if err != nil {
		return err
	}

	return CryptoVerify(signature, origin, period)
}
