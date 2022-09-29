package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/akakou/ecdaa"
)

const TPM_PATH = "/dev/tpm0"

type Config struct {
	Cred   *ecdaa.MiddleEncodedCredential
	Isk    *ecdaa.MiddleEncodedISK
	Ipk    *ecdaa.MiddleEncodedIPK
	Handle []byte
}

func main() {
	password := []byte("piyo")

	rng := ecdaa.InitRandom()

	tpm, err := ecdaa.OpenTPM(password, TPM_PATH)
	if err != nil {
		log.Fatalf("%v", err)
	}

	issuer := ecdaa.RandomIssuer(rng)

	err = ecdaa.VerifyIPK(&issuer.Ipk)
	if err != nil {
		log.Fatalf("%v", err)
	}

	seed, issuerSession, err := issuer.GenSeedForJoin(rng)
	if err != nil {
		log.Fatalf("%v", err)
	}

	member := ecdaa.NewMember(tpm)
	req, memberSession, err := member.GenReqForJoin(seed, rng)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tpm.Close()

	cipherCred, err := issuer.MakeCred(req, issuerSession, rng)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tpm, err = ecdaa.OpenTPM(password, TPM_PATH)
	if err != nil {
		log.Fatalf("%v", err)
	}

	tmp := member
	member = ecdaa.NewMember(tpm)
	member.KeyHandles = tmp.KeyHandles

	cred, err := member.ActivateCredential(cipherCred, memberSession, &issuer.Ipk)

	if err != nil {
		log.Fatalf("%v", err)
	}

	credBin := cred.Encode()

	iskBin := issuer.Isk.Encode()
	ipkBin := issuer.Ipk.Encode()
	handle := member.KeyHandles.Handle.Name.Buffer

	config := Config{
		Cred:   credBin,
		Isk:    iskBin,
		Ipk:    ipkBin,
		Handle: handle,
	}

	buf, err := json.Marshal(&config)

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	fmt.Printf("%s", buf)

}
