package android_scrappy

import (
	"encoding/json"
	"fmt"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy"
	"github.com/akakou/scrappy/crypto"
	"github.com/pkg/errors"
)

type AndroidResponse struct {
	Status string `json:"status"`
	Buffer []byte `json:"buffer"`
}

func AndroidJoin(host string, ipk []byte) string {
	configBuf, err := androidJoin(host, ipk)

	var resp AndroidResponse

	if err != nil {
		err = errors.Wrap(err, "android-error")
		resp = AndroidResponse{
			Status: "error",
			Buffer: []byte(err.Error()),
		}
	} else {
		resp = AndroidResponse{
			Status: "ok",
			Buffer: configBuf,
		}
	}

	r, _ := json.Marshal(&resp)

	return string(r)

}

func AndroidSign(origin string, period int, configBuf []byte) string {
	signature, err := androidSign(origin, period, configBuf)

	var resp AndroidResponse

	if err != nil {
		resp = AndroidResponse{
			Status: "error",
			Buffer: []byte(err.Error()),
		}
	} else {
		resp = AndroidResponse{
			Status: "ok",
			Buffer: []byte(signature),
		}
	}
	r, _ := json.Marshal(&resp)

	return string(r)
}

func androidJoin(host string, ipk []byte) ([]byte, error) {
	rng := ecdaa.InitRandom()

	config, err := scrappy.JoinForSigner(host, ipk, rng)

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

	if !scrappy.IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, scrappy.Now())
	}

	basename := scrappy.GetBasename(origin, period)

	signature, err := crypto.Sign([]byte(basename), &config)

	return signature, err
}
