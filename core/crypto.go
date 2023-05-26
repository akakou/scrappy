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

var CONFIG_PATH = "../config.json"

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

func CryptoSetup() error {
	rng := ecdaa.InitRandom()

	issuer := ecdaa.RandomIssuer(rng)

	err := ecdaa.VerifyIPK(&issuer.Ipk)
	if err != nil {
		return err
	}

	seed, issuerB, err := ecdaa.GenJoinSeed(rng)
	if err != nil {
		return err
	}

	tpm, err := ecdaa.OpenTPM([]byte(PASSWORD), TPM_PATH)

	if err != nil {
		return err
	}

	defer tpm.Close()

	req, handle, err := ecdaa.GenJoinReqWithTPM(seed, tpm, rng)
	if err != nil {
		return err
	}

	cipherCred, _, err := issuer.MakeCredEncrypted(req, issuerB, rng)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	cred, err := ecdaa.ActivateCredential(cipherCred, req.JoinReq.Proof.B, req.JoinReq.Q, &issuer.Ipk, handle, tpm)

	if err != nil {
		return err
	}

	credBin := cred.Encode()

	iskBin := issuer.Isk.Encode()
	ipkBin := issuer.Ipk.Encode()

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

func CryptoSign(basename string) (string, error) {
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

	authHandle := tpm2.AuthHandle{
		Handle: tpm2.TPMHandle(config.HandleNum),
		Name: tpm2.TPM2BName{
			Buffer: config.HandleName,
		},
		Auth: tpm2.PasswordAuth([]byte(PASSWORD)),
	}
	handle := ecdaa.KeyHandles{
		Handle: authHandle,
	}

	signature, err := ecdaa.SignTPM(
		[]byte{},
		[]byte(basename),
		config.Cred.Decode(),
		handle,
		tpm,
		rng,
	)

	if err != nil {
		return "", err
	}

	encodedSignature := signature.Encode()

	result, _ := json.Marshal(encodedSignature)

	return string(result), nil

}

func CryptoVerify(signatureBuf, basename string, rl ecdaa.RevocationList) error {
	var signature ecdaa.MiddleEncodedSignature

	config, err := readConfig()

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(signatureBuf), &signature)

	if err != nil {
		return err
	}

	err = ecdaa.Verify(
		[]byte{},
		[]byte(basename),
		signature.Decode(),
		config.Ipk.Decode(),
		rl,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetK(signatureBuf string) (string, error) {
	var signature ecdaa.MiddleEncodedSignature

	err := json.Unmarshal([]byte(signatureBuf), &signature)

	if err != nil {
		return "nil", err
	}

	result := base64.StdEncoding.EncodeToString(signature.Proof.K)
	return result, nil
}
