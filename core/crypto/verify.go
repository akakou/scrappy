package crypto

import "github.com/akakou/ecdaa"

func Verify(signatureBuf, basename []byte, rl ecdaa.RevocationList, config *VerifierConfig) error {
	var signature ecdaa.Signature
	err := signature.Decode(signatureBuf)
	if err != nil {
		return err
	}

	var ipk ecdaa.IPK
	err = ipk.Decode(config.IPK)
	if err != nil {
		return err
	}

	err = ecdaa.Verify(
		[]byte{},
		[]byte(basename),
		&signature,
		&ipk,
		rl,
	)

	return err
}

func VerifyWithConfig(signatureBuf, basename []byte, rl ecdaa.RevocationList) error {
	var config VerifierConfig
	err := ReadConfig(&config, VERIFIER_CONFIG_PATH)
	if err != nil {
		return err
	}

	return Verify(signatureBuf, basename, rl, &config)

}
