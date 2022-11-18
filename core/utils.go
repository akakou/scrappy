package core

import (
	b64 "encoding/base64"
	"miracl/core/FP256BN"
)

func curveToBase64(curve *FP256BN.ECP) string {
	var buf [int(FP256BN.MODBYTES) + 1]byte
	curve.ToBytes(buf[:], true)

	encoded := b64.StdEncoding.EncodeToString(buf[:])
	return encoded
}
