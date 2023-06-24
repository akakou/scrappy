//go:build tpm
// +build tpm

package scrappy

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/ecdaa_helper"
)

func SetupAllTPM() {
	rng := ecdaa.InitRandom()
	issuerConfig, err := ecdaa_helper.SetupIssuerAndSave(rng)

	if err != nil {
		log.Fatalf("setup issuer: %v\n", err)
	}

	issuer, err := ecdaa_helper.IssuerFromConfig(issuerConfig)

	if err != nil {
		log.Fatalf("setup issuer: %v\n", err)
	}

	err = ecdaa_helper.SignerTPMSetupAndSave(rng, issuer)

	if err != nil {
		log.Fatalf("setup issuer: %v\n", err)
	}

}

func SignTPM(origin string, period int) (string, error) {
	db, err := SetupDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH)

	if err != nil {
		return "", err
	}

	defer db.DB.Close()

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	basename := getBasename(origin, period)

	hasExist, err := HasExist(db, basename)

	if err != nil {
		return "", fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return "", fmt.Errorf(HAS_EXIST_ERROR, basename)
	}

	signature, err := ecdaa_helper.SignTPMWithConfig(basename)

	if err != nil {
		return "", err
	}

	err = Insert(db, basename)

	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(signature[:])

	return encoded, err

}
