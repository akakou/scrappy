package core

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	origin := "aaa"
	period := Now()

	err := Setup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signature, err := Sign(origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	K, err := GetK(signature)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	_, err = Sign(origin, period)

	if err.Error() != fmt.Sprintf(HAS_EXIST_ERROR, K) {
		t.Fatalf("%v: ", "failed to check basename deplication")
	}

	err = Verify(signature, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = Verify(signature, origin, period)

	if err.Error() != fmt.Sprintf(HAS_EXIST_ERROR, K) {
		t.Fatalf("%v: ", "failed to check basename deplication")
	}
}
