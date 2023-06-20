package main

import (
	"github.com/akakou/mcl_utils"
	"github.com/akakou/scrappy/ecdaa_helper"
)

func setup() error {
	rng := mcl_utils.InitRandom()

	issuer, err := ecdaa_helper.SetupIssuerAndSave(rng)
	if err != nil {
		return err
	}

	_, err = ecdaa_helper.InitSignerWithTPM(rng)
	if err != nil {
		return err
	}

	ecdaa_helper.VerifierAndSave(issuer)
	return nil
}
