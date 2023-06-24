package scrappy

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestDB(t *testing.T) {
	value := "valuea"

	conf := DB{
		Table:  "TEST_TABLE",
		Column: "TEST_COLUMN",
	}

	db, err := SetupDB(conf, "test.db")

	if err != nil {
		panic(err)
	}

	hasExist, err := HasHashExist(db, []byte(value))

	if err != nil {
		t.Fatalf("%v", err)
	}

	if hasExist {
		t.Fatalf("%v", value)
	}

	err = InsertHash(db, value)
	if err != nil {
		t.Fatalf("%v", err)
	}

	hasExist, err = HasHashExist(db, []byte(value))

	if err != nil {
		t.Fatalf("%v", err)
	}

	if !hasExist {
		t.Fatalf("basename %s does not exists", value)
	}

	strs, err := SelectAll(db)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(strs) != 1 {
		t.Fatalf("%v", err)
	}
}
