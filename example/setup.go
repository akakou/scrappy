package main

import (
	amclutils "github.com/akakou/fp256bn-amcl-utils"
	"github.com/akakou/scrappy/ecdaa_helper"
)

func setup() error {
	rng := amclutils.InitRandom()

	issuer, err := ecdaa_helper.SetupIssuerAndSave(rng)
	if err != nil {
		return err
	}

	_, err = ecdaa_helper.InitSignerWithTPM(rng)
	if err != nil {
		return err
	}

	ecdaa_helper.VerifierAndSave(issuer)

	signer, err := ecdaa_helper.ExampleInitSigner(rng)

	if err != nil {
		return err
	}

	err = ecdaa_helper.WriteConfig(&signer, ecdaa_helper.SIGNER_CONFIG_PATH)

	if err != nil {
		return err
	}

	return nil
}
