package scrappy

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/crypto"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var encodeBase64 = base64.StdEncoding.EncodeToString
var decodeBase64 = base64.StdEncoding.DecodeString

func getBasename(origin string, period int) string {
	hashed_origin := hashAndEncodeBase64([]byte(origin))
	basename := fmt.Sprintf("%v_%v", hashed_origin, period)

	return basename
}

func hashAndEncodeBase64(data []byte) string {
	hash := sha256.Sum256(data)
	encoded := encodeBase64(hash[:])
	return encoded
}

func GetK(signatureBuf []byte) (string, error) {
	var signature ecdaa.Signature
	err := signature.Decode(signatureBuf)

	if err != nil {
		return "nil", err
	}

	var KBuf [33]byte
	signature.Proof.K.ToBytes(KBuf[:], true)

	result := encodeBase64(KBuf[:])
	return result, nil
}

var SetupIssuerAndSave = crypto.SetupIssuerAndSave

func ginServ() *gin.Engine {
	secret := []byte("secret")
	r := gin.Default()
	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*.html")

	return r
}
