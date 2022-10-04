package core

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/akakou/ecdaa"
)

const PASSWORD = "password"
const TPM_PATH = "/dev/tpm0"
const CONFIG_PATH = "../init/config.json"

type Config struct {
	Cred   *ecdaa.MiddleEncodedCredential
	Isk    *ecdaa.MiddleEncodedISK
	Ipk    *ecdaa.MiddleEncodedIPK
	Handle []byte
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

func Sign(period string) (string, error) {
	rng := ecdaa.InitRandom()

	config, err := readConfig()

	if err != nil {
		return "", err
	}

	tpm, err := ecdaa.OpenTPM([]byte(PASSWORD), TPM_PATH)
	if err != nil {
		return "", err
	}

	member := ecdaa.NewMember(tpm)
	// todo: set handle

	signature, err := member.Sign(
		[]byte("test"),
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
