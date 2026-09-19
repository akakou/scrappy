//go:build tpm
// +build tpm

package scrappy

import (
	"testing"
)

func TestTPM(t *testing.T) {
	SetupAllTPM()
	signature, err := SignTPM("example.com", 1)
	if err != nil {
		t.Fatalf("%v: ", err)
	}

	isNotValid := Verify(signature, "example.com", 1)
	if isNotValid != nil {
		t.Fatalf("%v: ", err)
	}
}
