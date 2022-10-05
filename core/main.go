package core

func Setup() {
	CryptoSetup()
}

func Sign(origin string, period int) (string, error) {
	return CryptoSign(origin, period)
}

func Verify(signature, origin string, period int) error {
	return CryptoVerify(signature, origin, period)
}
