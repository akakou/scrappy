package ecdaa_helper

import (
	"fmt"
	"miracl/core"
	"miracl/core/FP256BN"

	"github.com/akakou/ecdaa"
	"github.com/akakou/ecdaa/tpm_utils"
)

func GenSeed(rng *core.RAND) ([]byte, []byte, error) {
	seed, issuerB, err := ecdaa.GenJoinSeed(rng)

	if err != nil {
		panic(err)
	}

	seedBuf, err := seed.Encode()

	if err != nil {
		panic(err)
	}

	var encodedB [33]byte
	issuerB.ToBytes(encodedB[:], true)

	return seedBuf, encodedB[:], err
}

func GenJoinReq(seedBuf []byte, rng *core.RAND) ([]byte, []byte, error) {
	var seed ecdaa.JoinSeed
	err := seed.Decode(seedBuf)
	if err != nil {
		fmt.Print("cant decode seed")
		return nil, nil, err
	}

	joinReq, sk, err := ecdaa.GenJoinReq(&seed, rng)

	if err != nil {
		return nil, nil, err
	}

	joinReqBuf, err := joinReq.Encode()

	if err != nil {
		return nil, nil, err
	}

	var skBuf [32]byte
	sk.ToBytes(skBuf[:])

	return joinReqBuf, skBuf[:], err
}

func MakeCred(joinReqBuf []byte, issuerBBuf []byte, issuer *ecdaa.Issuer, rng *core.RAND) ([]byte, error) {
	var joinReq ecdaa.JoinRequest
	err := joinReq.Decode(joinReqBuf)

	if err != nil {
		return nil, err
	}

	issuerB := FP256BN.ECP_fromBytes(issuerBBuf[:])
	cred, err := issuer.MakeCred(&joinReq, issuerB, rng)

	if err != nil {
		return nil, err
	}

	credBin, err := cred.Encode()

	return credBin, err
}

func VerifyCred(credBuf []byte, ipkBuf []byte) error {
	var ipk ecdaa.IPK
	err := ipk.Decode(ipkBuf)

	if err != nil {
		return err
	}

	var cred ecdaa.Credential
	err = cred.Decode(credBuf)

	if err != nil {
		return err
	}

	return ecdaa.VerifyCred(&cred, &ipk)

}

func MakeCredWithTPM(joinReqBuf []byte, issuerBBuf []byte, rng *core.RAND) ([]byte, error) {
	issuer, err := IssuerFromConfigFile()

	if err != nil {
		return nil, err
	}

	return MakeCred(joinReqBuf, issuerBBuf, issuer, rng)
}

func Join(issuer *ecdaa.Issuer, rng *core.RAND) (*SignerConfig, error) {
	seed, issuerB, err := ecdaa.GenJoinSeed(rng)
	if err != nil {
		return nil, err
	}

	req, sk, err := ecdaa.GenJoinReq(seed, rng)
	if err != nil {
		return nil, err
	}

	cred, err := issuer.MakeCred(req, issuerB, rng)
	if err != nil {
		return nil, err
	}

	credBin, err := cred.Encode()
	if err != nil {
		return nil, err
	}

	ipkBin, err := issuer.Ipk.Encode()
	if err != nil {
		return nil, err
	}

	skBin := [32]byte{}
	sk.ToBytes(skBin[:])
	if err != nil {
		return nil, err
	}

	config := SignerConfig{
		Cred: credBin,
		IPK:  ipkBin,
		SK:   skBin[:],
	}

	return &config, nil
}

func JoinTPM(issuer *ecdaa.Issuer, rng *core.RAND) (*SignerConfigTPM, error) {
	tpm, err := tpm_utils.OpenTPM([]byte(PASSWORD), TPM_PATH)

	if err != nil {
		return nil, err
	}

	defer tpm.Close()

	seed, issuerB, err := ecdaa.GenJoinSeed(rng)
	if err != nil {
		return nil, err
	}

	req, handle, err := ecdaa.GenJoinReqWithTPM(seed, tpm, rng)
	if err != nil {
		return nil, err
	}

	cipherCred, _, err := issuer.MakeCredEncrypted(req, issuerB, rng)
	if err != nil {
		return nil, err
	}

	cred, err := ecdaa.ActivateCredential(cipherCred, issuerB, req.JoinReq.Q, &issuer.Ipk, handle, tpm)
	if err != nil {
		return nil, err
	}

	credBin, err := cred.Encode()
	if err != nil {
		return nil, err
	}

	ipkBin, err := issuer.Ipk.Encode()
	if err != nil {
		return nil, err
	}

	config := SignerConfigTPM{
		Cred:       credBin,
		IPK:        ipkBin,
		HandleNum:  handle.Handle.HandleValue(),
		HandleName: handle.Handle.Name.Buffer,
	}

	return &config, nil
}

func ExampleInitSigner(rng *core.RAND) (*SignerConfig, error) {
	issuer, err := IssuerFromConfigFile()

	if err != nil {
		return nil, err
	}

	tpmConfig, err := Join(issuer, rng)

	if err != nil {
		return nil, err
	}

	err = WriteConfig(tpmConfig, SIGNER_CONFIG_PATH)

	return tpmConfig, err
}

func InitSignerWithTPM(rng *core.RAND) (*SignerConfigTPM, error) {
	issuer, err := IssuerFromConfigFile()

	if err != nil {
		return nil, err
	}

	tpmConfig, err := JoinTPM(issuer, rng)

	if err != nil {
		return nil, err
	}

	err = WriteConfig(tpmConfig, SIGNER_TPM_CONFIG_PATH)

	return tpmConfig, err
}
