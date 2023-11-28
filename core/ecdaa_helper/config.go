package ecdaa_helper

import (
	"encoding/json"
	"io/ioutil"

	"github.com/akakou/ecdaa"
)

type IssuerConfig struct {
	ISK []byte
	IPK []byte
}

func IssuerFromConfigFile() (*ecdaa.Issuer, error) {
	var config IssuerConfig

	ipk := ecdaa.IPK{}
	isk := ecdaa.ISK{}

	err := ReadConfig(&config, ISSUER_CONFIG_PATH)
	if err != nil {
		return nil, err
	}

	err = ipk.Decode(config.IPK)
	if err != nil {
		return nil, err
	}

	err = isk.Decode(config.ISK)

	issuer := ecdaa.NewIssuer(isk, ipk)

	return &issuer, err
}

type SignerConfigTPM struct {
	IPK        []byte
	Cred       []byte
	HandleName []byte
	HandleNum  uint32
}

type SignerConfig struct {
	IPK  []byte
	Cred []byte
	SK   []byte
}

type VerifierConfig struct {
	IPK []byte
	RL  [][]byte
}

func ReadConfig[T any](config T, path string) error {
	configBuf, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	err = json.Unmarshal(configBuf, config)
	return err
}

func WriteConfig[T any](config T, path string) error {
	buf, err := json.Marshal(config)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf, 0644)

	return err
}
