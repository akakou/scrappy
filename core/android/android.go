package android_scrappy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	amclutils "github.com/akakou/fp256bn-amcl-utils"
	"github.com/akakou/scrappy"
	"github.com/akakou/scrappy/ecdaa_helper"
	"github.com/pkg/errors"
)

type AndroidResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Data   any    `json:"data"`
}

func buildAndroidMessage(data any, err error) string {
	var resp AndroidResponse

	if err != nil {
		err = errors.Wrap(err, "android-error")
		resp = AndroidResponse{
			Status: "error",
			Error:  string(err.Error()),
		}
	} else {
		resp = AndroidResponse{
			Status: "ok",
			Data:   data,
		}
	}

	r, _ := json.Marshal(&resp)

	return string(r)
}

func Hello() {
}

func RateLimitK() int {
	return scrappy.RateLimitK
}

func HashBasename(origin string) string {
	return scrappy.HashBasename(origin)
}

func AndroidJoin(host, base64Ipk string) string {
	config, err := androidJoin(host, base64Ipk)
	resp := buildAndroidMessage(config, err)

	return resp
}

func androidJoin(host, base64Ipk string) (*ecdaa_helper.SignerConfig, error) {
	rng := amclutils.InitRandom()

	ipk, err := base64.StdEncoding.DecodeString(base64Ipk)
	if err != nil {
		return nil, err
	}

	config, err := scrappy.JoinForSigner(host, ipk, rng)
	return config, err
}

func AndroidSign(origin string, period, i int, sk, cred, ipk string) string {
	signature, err := androidSign(origin, period, i, sk, cred, ipk)
	resp := buildAndroidMessage(signature, err)

	return resp
}

func androidSign(origin string, period, i int, sk, cred, ipk string) (string, error) {
	SK, err := base64.StdEncoding.DecodeString(sk)
	if err != nil {
		return "", err
	}

	Cred, err := base64.StdEncoding.DecodeString(cred)
	if err != nil {
		return "", err
	}

	IPK, err := base64.StdEncoding.DecodeString(ipk)
	if err != nil {
		return "", err
	}

	var config = ecdaa_helper.SignerConfig{
		IPK: IPK, Cred: Cred, SK: SK,
	}

	if !scrappy.IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, scrappy.Now())
	}
	if i < 0 || i >= scrappy.RateLimitK {
		return "", fmt.Errorf("invalid proof index %d for rate limit %d", i, scrappy.RateLimitK)
	}

	basename := scrappy.GetBasenameWithIndex(origin, period, i)

	signer, err := ecdaa_helper.PrepareSWSigner(&config)
	if err != nil {
		return "", err
	}

	signature, err := ecdaa_helper.SignWithEncoding(basename, signer)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d:%s", i, signature), nil
}
