package crypto

import (
	"testing"

	"github.com/akakou/ecdaa"
)

func TestSW(t *testing.T) {
	rng := ecdaa.InitRandom()
	issuerConf, err := SetupIssuerAndSave(rng)

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

	signature, err := Sign([]byte("basename"), &singerConfig)

	if err != nil {
		t.Fatal(err)
	}

	verifierConfig := VerifierConfig{
		IPK: issuerConf.IPK,
	}

	err = Verify(signature, []byte("basename"), ecdaa.RevocationList{}, &verifierConfig)

	if err != nil {
		t.Fatal(err)
	}

}
