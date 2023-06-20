package crypto

import (
	"miracl/core/FP256BN"

	"github.com/akakou/ecdaa"
	"github.com/akakou/ecdaa/tpm_utils"
	"github.com/akakou/mcl_utils"
	"github.com/google/go-tpm/tpm2"
)

func SWSign(basename string, config *SignerConfig) (string, error) {
	sk := FP256BN.FromBytes(config.SK)

	var cred ecdaa.Credential
	err := cred.Decode(config.Cred)
	if err != nil {
		return "", err
	}

	signer := ecdaa.NewSWSigner(&cred, sk)
	return coreSign(basename, signer)
}

func TPMSign(basename string, config *SignerConfigTPM) (string, error) {
	tpm, err := tpm_utils.OpenTPM([]byte(PASSWORD), TPM_PATH)
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
	err = cred.Decode(config.Cred)
	if err != nil {
		return "", err
	}

	signer := ecdaa.NewTPMSigner(&cred, &handle, tpm)

	return coreSign(basename, &signer)
}

func TPMSignWithConfig(basename string) (string, error) {
	var config SignerConfigTPM
	err := ReadConfig(&config, SIGNER_TPM_CONFIG_PATH)

	if err != nil {
		return "", err
	}

	return TPMSign(basename, &config)
}

func coreSign(basename string, signer ecdaa.Signer) (string, error) {
	rng := mcl_utils.InitRandom()

	signature, err := signer.Sign(
		[]byte{},
		[]byte(basename),
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
