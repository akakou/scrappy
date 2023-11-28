package ecdaa_helper

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

func PrepareVerifierConfig() (*ecdaa.IPK, *ecdaa.RevocationList, error) {
	var config VerifierConfig
	err := ReadConfig(&config, VERIFIER_CONFIG_PATH)
	if err != nil {
		return nil, nil, err
	}

	var ipk ecdaa.IPK
	err = ipk.Decode(config.IPK)
	if err != nil {
		return nil, nil, err
	}

	rl := ecdaa.DecodeRevocationList(config.RL)

	return &ipk, &rl, nil

}
