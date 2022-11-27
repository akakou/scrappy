package core

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"miracl/core/FP256BN"
)

func curveToBase64(curve *FP256BN.ECP) string {
	var buf [int(FP256BN.MODBYTES) + 1]byte
	curve.ToBytes(buf[:], true)

	encoded := b64.StdEncoding.EncodeToString(buf[:])
	return encoded
}

func hashAndEncodeBase64(data []byte) string {
	hash := sha256.Sum256(data)
	encoded := b64.StdEncoding.EncodeToString(hash[:])
	return encoded
}

func getBasename(origin string, period int) string {
	hashed_origin := hashAndEncodeBase64([]byte(origin))
	basename := fmt.Sprintf("%v_%v", hashed_origin, period)

	return basename
}
