package core

import "testing"

func TestAll(t *testing.T) {
	origin := "aaa"
	period := 1000

	err := Setup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signatuere, err := Sign(origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	_, err = Sign(origin, period)

	if err.Error() != "basename aaa with period 1000 already exists" {
		t.Fatalf("%v: ", "failed to check basename deplication")
	}

	err = Verify(signatuere, origin, period)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = Verify(signatuere, origin, period)

	if err.Error() != "basename aaa with period 1000 already exists" {
		t.Fatalf("%v: ", "failed to check basename deplication")
	}
}
