package core

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestDB(t *testing.T) {
	basename := "basename"

	db, err := SetupDB(TEST_DB_PATH)

	if err != nil {
		panic(err)
	}

	err = InsertIfItHasNotExist(db, basename)

	if err != nil {
		t.Fatalf("%v", err)
	}

	hasExist, err := HasExist(db, basename)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if !hasExist {
		t.Fatalf("basename %s does not exists", basename)
	}
}
