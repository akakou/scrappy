package crypto

import (
	"miracl/core"

	"github.com/akakou/ecdaa"
)

const PASSWORD = "password"
const TPM_PATH = "/dev/tpm0"

func SetupIssuer(rng *core.RAND) (*IssuerConfig, error) {
	issuer := ecdaa.RandomIssuer(rng)

	err := ecdaa.VerifyIPK(&issuer.Ipk)
	if err != nil {
		return nil, err
	}

	isk, err := issuer.Isk.Encode()
	if err != nil {
		return nil, err
	}

	ipk, err := issuer.Ipk.Encode()
	if err != nil {
		return nil, err
	}

	issuerConfig := IssuerConfig{
		ISK: isk,
		IPK: ipk,
	}

	return &issuerConfig, nil
}

func SetupIssuerAndSave(rng *core.RAND) (*IssuerConfig, error) {
	issuerConfig, err := SetupIssuer(rng)

	if err != nil {
		return nil, err
	}

	err = WriteConfig(issuerConfig, ISSUER_CONFIG_PATH)

	return issuerConfig, err
}

// var SignerSetup = Join

// func SignerSetupAndSave(rng *core.RAND, issuer *ecdaa.Issuer) error {
// 	signerConfig, err := SignerSetup(rng, issuer)

// 	if err != nil {
// 		return err
// 	}

// 	err = writeConfig(signerConfig, SIGNER_CONFIG_PATH)

// 	return err
// }
