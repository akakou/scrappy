package core

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestDB(t *testing.T) {
	origin := "aaa"
	period := 1000

	db, err := SetupDB(TEST_DB_PATH)

	if err != nil {
		panic(err)
	}

	err = InsertIfItDoesNotExist(db, origin, period)

	if err != nil {
		t.Fatalf("%v", err)
	}

	hasExist, err := HasExist(db, origin, period)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if !hasExist {
		t.Fatalf("basename %s with period %d does not exists", origin, period)
	}
}
