package core

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestDB(t *testing.T) {
	value := "value"

	db := &DB{
		Table:  "TEST_TABLE",
		Column: "TEST_COLUMN",
	}

	db, err := SetupDB(db, ":memory:")

	if err != nil {
		panic(err)
	}

	hasExist, err := HasExist(db, value)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if hasExist {
		t.Fatalf("%v", value)
	}

	err = Insert(db, value)

	if err != nil {
		t.Fatalf("%v", err)
	}

	hasExist, err = HasExist(db, value)

	if err != nil {
		t.Fatalf("%v", err)
	}

if !hasExist {
		t.Fatalf("basename %s does not exists", value)
	}
}
