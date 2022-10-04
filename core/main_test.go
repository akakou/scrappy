package core

import "testing"

func TestAll(t *testing.T) {
	basename := "aaa"

	err := Setup()

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	signatuere, err := Sign(basename)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

	err = Verify(signatuere, basename)

	if err != nil {
		t.Fatalf("%v: ", err)
	}

}
