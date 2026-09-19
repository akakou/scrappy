package scrappy

import (
	"io"
	"miracl/core"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/akakou/scrappy/ecdaa_helper"
)

func JoinForSigner(HOST string, IPKBuf []byte, rng *core.RAND) (*ecdaa_helper.SignerConfig, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.Get(HOST + "/gen_seed")

	if err != nil {
		return nil, err
	}

	seed, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	req, sk, err := ecdaa_helper.GenJoinReq(seed, rng)

	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("join_req", string(req))

	resp, err = client.PostForm(HOST+"/make_cred", values)

	if err != nil {
		return nil, err
	}

	cred, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = ecdaa_helper.VerifyCred(cred, IPKBuf)

	if err != nil {
		return nil, err
	}

	config := ecdaa_helper.SignerConfig{
		IPK:  IPKBuf,
		Cred: cred,
		SK:   sk,
	}

	return &config, err

}
