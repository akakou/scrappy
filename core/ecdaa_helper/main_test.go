package ecdaa_helper

import (
	"testing"

	"github.com/akakou/ecdaa"
	amclutils "github.com/akakou/fp256bn-amcl-utils"
)

func TestSW(t *testing.T) {
	rng := amclutils.InitRandom()
	issuerConf, err := SetupIssuerAndSave(rng)

	if err != nil {
		t.Fatal(err)
	}

	_, err = VerifierAndSave(issuerConf)

	if err != nil {
		t.Fatal(err)
	}

	issuer, err := IssuerFromConfigFile()

	if err != nil {
		t.Fatal(err)
	}

	seed, issuerB, err := GenSeed(rng)

	if err != nil {
		t.Fatal(err)
	}

	joinReq, sk, err := GenJoinReq(seed, rng)

	if err != nil {
		t.Fatal(err)
	}

	cred, err := MakeCred(joinReq, issuerB, issuer, rng)

	if err != nil {
		t.Fatal(err)
	}

	err = VerifyCred(cred, issuerConf.IPK)

	if err != nil {
		t.Fatal(err)
	}

	singerConfig := SignerConfig{
		IPK:  issuerConf.IPK,
		Cred: cred,
		SK:   sk,
	}

	signer, err := PrepareSWSigner(&singerConfig)

	if err != nil {
		t.Fatal(err)
	}

	sigantureStr, err := SignWithEncoding("basename", signer)

	if err != nil {
		t.Fatal(err)
	}

	signature, err := DecodeSignature(sigantureStr)

	if err != nil {
		t.Fatal(err)
	}

	ipk, rl, err := PrepareVerifierConfig()

	if err != nil {
		t.Fatal(err)
	}

	err = ecdaa.Verify([]byte{}, []byte("basename"), signature, ipk, *rl)

	if err != nil {
		t.Fatal(err)
	}

}
