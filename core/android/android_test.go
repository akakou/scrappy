package android_scrappy

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy"
	"github.com/akakou/scrappy/crypto"
	scrappy_gin "github.com/akakou/scrappy/gin"
	_ "github.com/mattn/go-sqlite3"
)

const WAIT_TIME = 2

func checkError(err error, t *testing.T, i int) {
	if err != nil {
		t.Fatalf("%d: %v", i, err)
	}
}

func TestAndroidAndGin(t *testing.T) {
	now := scrappy.Now()
	rng := ecdaa.InitRandom()
	issuer, err := crypto.SetupIssuer(rng)

	checkError(err, t, 0)

	err = crypto.WriteConfig(issuer, crypto.ISSUER_CONFIG_PATH)

	checkError(err, t, 1)

	verifier := crypto.VerifierConfig{
		IPK: issuer.IPK,
	}

	err = crypto.WriteConfig(&verifier, crypto.VERIFIER_CONFIG_PATH)

	if err != nil {
		t.Fatal(err)
	}

	go scrappy_gin.RunIssuer("../templates/*.html")

	time.Sleep(WAIT_TIME * time.Second)

	var respJoin AndroidResponse

	base64IPK := base64.StdEncoding.EncodeToString(issuer.IPK)

	respJson := AndroidJoin("http://127.0.0.1:8080", base64IPK)
	err = json.Unmarshal([]byte(respJson), &respJoin)
	checkError(err, t, 2)
	if respJoin.Status != "ok" {
		t.Fatalf("status should be ok: %v", respJoin.Error)
	}

	signerConfig := respJoin.Data.(map[string]interface{})

	var respSign AndroidResponse
	respJson = AndroidSign("http://localhost:8080", now, signerConfig["SK"].(string), signerConfig["Cred"].(string), signerConfig["IPK"].(string))
	err = json.Unmarshal([]byte(respJson), &respSign)
	if respSign.Status != "ok" {
		t.Fatalf("status should be ok: %v", respSign.Error)
	}
	checkError(err, t, 3)

	err = scrappy.Verify(respSign.Data.(string), "http://localhost:8080", now)

	checkError(err, t, 4)

}
