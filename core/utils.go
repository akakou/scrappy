package scrappy

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/ecdaa_helper"
)

var encodeBase64 = base64.StdEncoding.EncodeToString
var decodeBase64 = base64.StdEncoding.DecodeString

func GetBasename(origin string, period int) string {
	hashed_origin := hashAndEncodeBase64([]byte(origin))
	basename := fmt.Sprintf("%v_%v", hashed_origin, period)

	return basename
}

func hashAndEncodeBase64(data []byte) string {
	hash := sha256.Sum256(data)
	encoded := encodeBase64(hash[:])
	return encoded
}

func GetKBytes(signature *ecdaa.Signature) string {
	var KBuf [33]byte
	signature.Proof.K.ToBytes(KBuf[:], true)

	result := encodeBase64(KBuf[:])
	return result
}

var SetupIssuerAndSave = ecdaa_helper.SetupIssuerAndSave

func CheckBasenameExists(basename string, db *DB) error {
	hasExist, err := HasExist(db, basename)

	if err != nil {
		return fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return fmt.Errorf(HAS_EXIST_ERROR, basename)
	}

	return nil
}
