package crypto

import (
	"miracl/core/FP256BN"

	"github.com/akakou/ecdaa"
	"github.com/google/go-tpm/tpm2"
)

func Sign(basename []byte, config *SignerConfig) (string, error) {
	rng := ecdaa.InitRandom()

	var cred ecdaa.Credential
	err := cred.Decode(config.Cred)

	if err != nil {
		return "", err
	}

	var ipk ecdaa.IPK
	err = ipk.Decode(config.IPK)

	if err != nil {
		return "", err
	}

	sk := FP256BN.FromBytes(config.SK)

	signature, err := ecdaa.Sign(
		[]byte{},
		basename,
		sk,
		&cred,
		rng,
	)

	if err != nil {
		return "", err
	}

	encodedSignature, err := signature.Encode()

	if err != nil {
		return "", err
	}

	return encodeBase64(encodedSignature), err
}

func SignTPM(basename string, config *SignerConfigTPM) (string, error) {
	rng := ecdaa.InitRandom()

	tpm, err := ecdaa.OpenTPM([]byte(PASSWORD), TPM_PATH)
	if err != nil {
		return "", err
	}

	defer tpm.Close()

	authHandle := tpm2.AuthHandle{
		Handle: tpm2.TPMHandle(config.HandleNum),
		Name: tpm2.TPM2BName{
			Buffer: config.HandleName,
		},
		Auth: tpm2.PasswordAuth([]byte(PASSWORD)),
	}

	handle := ecdaa.KeyHandles{
		EkHandle:  &tpm2.AuthHandle{},
		SrkHandle: &tpm2.NamedHandle{},
		Handle:    &authHandle,
	}

	var cred ecdaa.Credential
	cred.Decode(config.Cred)

	signature, err := ecdaa.SignTPM(
		[]byte{},
		[]byte(basename),
		&cred,
		&handle,
		tpm,
		rng,
	)

	if err != nil {
		return "", err
	}

	encodedSignature, err := signature.Encode()

	if err != nil {
		return "", err
	}

	return encodeBase64(encodedSignature), err
}

func SignTPMWithConfig(basename string) (string, error) {
	var config SignerConfigTPM
	err := ReadConfig(&config, SIGNER_TPM_CONFIG_PATH)

	if err != nil {
		return "", err
	}

	return SignTPM(basename, &config)
}
