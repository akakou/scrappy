package crypto

import "github.com/akakou/ecdaa"

func DecodeSignature(signatureString string) (*ecdaa.Signature, error) {
	signatureBuf, err := decodeBase64(signatureString)
	if err != nil {
		return nil, err
	}

	var signature ecdaa.Signature
	err = signature.Decode(signatureBuf)
	if err != nil {
		return nil, err
	}

	return &signature, nil
}

func VerifyWithConfig(signature *ecdaa.Signature, basename []byte) error {
	var config VerifierConfig
	err := ReadConfig(&config, VERIFIER_CONFIG_PATH)
	if err != nil {
		return err
	}

	var ipk ecdaa.IPK
	err = ipk.Decode(config.IPK)
	if err != nil {
		return err
	}

	rl := ecdaa.DecodeRevocationList(config.RL)

	return ecdaa.Verify([]byte{}, basename, signature, &ipk, rl)

}
