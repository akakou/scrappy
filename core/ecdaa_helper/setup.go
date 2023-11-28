package ecdaa_helper

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

func VerifierAndSave(isser *IssuerConfig) (*VerifierConfig, error) {
	verifierConfig := VerifierConfig{
		IPK: isser.IPK,
		RL:  [][]byte{},
	}

	err := WriteConfig(verifierConfig, VERIFIER_CONFIG_PATH)

	return &verifierConfig, err
}
