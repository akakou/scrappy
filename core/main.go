package core

func Setup() error {
	return CryptoSetup()
}

func Sign(origin string, period int) (string, error) {
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
