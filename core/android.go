package scrappy

import (
	"fmt"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/crypto"
)

var AndroidSign = androidSign
var AndroidJoin = androidJoin

func androidJoin(host string, ipk []byte) ([]byte, error) {
	rng := ecdaa.InitRandom()

	config, err := JoinForSigner(host, ipk, rng)

	if err != nil {
		fmt.Printf("err")
		return nil, err
	}

	configBuf, err := ecdaa.Encode(config)
	return configBuf, err
}

func androidSign(origin string, period int, configBuf []byte) (string, error) {
	var config crypto.SignerConfig
	err := ecdaa.Decode(&config, configBuf)

	if err != nil {
		return "", err
	}

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	basename := getBasename(origin, period)

	signature, err := crypto.Sign([]byte(basename), &config)

	if err != nil {
		return "", err
	}

	return encodeBase64(signature), err
}
