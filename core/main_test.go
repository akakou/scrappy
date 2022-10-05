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

	signatuere, err := Sign(origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	_, err = Sign(origin, period)

	if err.Error() != fmt.Sprintf("basename %v with period %d already exists", origin, period) {
		t.Fatalf("%v: ", "failed to check basename deplication")
	}

	err = Verify(signatuere, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = Verify(signatuere, origin, period)

	if err.Error() != fmt.Sprintf("basename %v with period %d already exists", origin, period) {
		t.Fatalf("%v: ", "failed to check basename deplication")
	}
}
