package scrappy

import (
	"io"
	"miracl/core"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/akakou/scrappy/crypto"
)

func JoinForSigner(HOST string, IPKBuf []byte, rng *core.RAND) (*crypto.SignerConfig, error) {
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

	req, sk, err := crypto.GenJoinReq(seed, rng)

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

	err = crypto.VerifyCred(cred, IPKBuf)

	if err != nil {
		return nil, err
	}

	config := crypto.SignerConfig{
		IPK:  IPKBuf,
		Cred: cred,
		SK:   sk,
	}

	return &config, err

}
