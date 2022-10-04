package core

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/akakou/ecdaa"
	"github.com/google/go-tpm/tpm2"
)

const PASSWORD = "password"
const TPM_PATH = "/dev/tpm0"
const CONFIG_PATH = "../config.json"

type Config struct {
	Cred       *ecdaa.MiddleEncodedCredential
	Isk        *ecdaa.MiddleEncodedISK
	Ipk        *ecdaa.MiddleEncodedIPK
	HandleName []byte
	HandleNum  uint32
}

func readConfig() (*Config, error) {
	var config Config
	configBuf, err := ioutil.ReadFile(CONFIG_PATH)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configBuf, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func Setup() error {
	rng := ecdaa.InitRandom()

	issuer := ecdaa.RandomIssuer(rng)

	err := ecdaa.VerifyIPK(&issuer.Ipk)
	if err != nil {
		return err
	}

	seed, issuerSession, err := issuer.GenSeedForJoin(rng)
	if err != nil {
		return err
	}

	tpm, err := ecdaa.OpenTPM([]byte(PASSWORD), TPM_PATH)

	if err != nil {
		return err
	}

	defer tpm.Close()

	member := ecdaa.NewMember(tpm)
	req, memberSession, err := member.GenReqForJoin(seed, rng)
	if err != nil {
		return err
	}

	cipherCred, err := issuer.MakeCred(req, issuerSession, rng)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	tmp := member
	member = ecdaa.NewMember(tpm)
	member.KeyHandles = tmp.KeyHandles

	cred, err := member.ActivateCredential(cipherCred, memberSession, &issuer.Ipk)

	if err != nil {
		return err
	}

	credBin := cred.Encode()

	iskBin := issuer.Isk.Encode()
	ipkBin := issuer.Ipk.Encode()
	handle := member.KeyHandles.Handle

	config := Config{
		Cred:       credBin,
		Isk:        iskBin,
		Ipk:        ipkBin,
		HandleNum:  handle.HandleValue(),
		HandleName: handle.Name.Buffer,
	}

	buf, err := json.Marshal(&config)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(CONFIG_PATH, buf, 0644)

	if err != nil {
		return err
	}

	return nil
}

func Sign(period string) (string, error) {
	rng := ecdaa.InitRandom()

	tpm, err := ecdaa.OpenTPM([]byte(PASSWORD), TPM_PATH)
	if err != nil {
		return "", err
	}

	defer tpm.Close()

	config, err := readConfig()

	if err != nil {
		return "", err
	}

	handle := tpm2.AuthHandle{
		Handle: tpm2.TPMHandle(config.HandleNum),
		Name: tpm2.TPM2BName{
			Buffer: config.HandleName,
		},
		Auth: tpm2.PasswordAuth([]byte(PASSWORD)),
	}

	member := ecdaa.NewMember(tpm)

	member.KeyHandles = &ecdaa.KeyHandles{
		Handle: &handle,
	}

	signature, err := member.Sign(
		[]byte{},
		[]byte(period),
		config.Cred.Decode(),
		rng,
	)

	if err != nil {
		return "", err
	}

	encodedSignature := signature.Encode()

	jsonSignature, _ := json.Marshal(encodedSignature)

	result := base64.StdEncoding.EncodeToString([]byte(jsonSignature))

	return result, nil

}

func Verify(base64Signature, period string) error {
	var signature ecdaa.MiddleEncodedSignature

	config, err := readConfig()

	if err != nil {
		return err
	}

	attestBuf, err := base64.StdEncoding.DecodeString(base64Signature)

	if err != nil {
		return err
	}
	err = json.Unmarshal(attestBuf, &signature)

	if err != nil {
		return err
	}

	err = ecdaa.Verify(
		[]byte{},
		[]byte(period),
		signature.Decode(),
		config.Ipk.Decode(),
	)

	if err != nil {
		return err
	}

	return nil
}
