package core

import (
	"database/sql"
	"fmt"
	"os/exec"
)

var SIGNER_DB_PATH = "./signer.db"
var VERIFIER_DB_PATH = "./verifier.db"
var TEST_DB_PATH = "./test.db"

const HAS_EXIST_ERROR = "k %v already exists"

type DB struct {
	Table  string
	Column string
	DB     *sql.DB
}

func SetupSignerDB(path string) (*DB, error) {
	signerDB := DB{
		Table:  "SIGNER_LOG",
		Column: "BASENAME",
	}

	return SetupDB(&signerDB, path)
}

func SetupVerifierDB(path string) (*DB, error) {
	verifirDB := DB{
		Table:  "VERIFIER_LOG",
		Column: "K",
	}

	return SetupDB(&verifirDB, path)
}

func SetupDB(db *DB, path string) (*DB, error) {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v VARCHAR(1024))", db.Table, db.Column)

	if path != TEST_DB_PATH {
		exec.Command("touch", path).Run()
	}

	_db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	_, err = _db.Exec(query)

	if err != nil {
		return nil, err
	}

	db.DB = _db

	return db, nil
}

func HasExist(db *DB, value string) (bool, error) {
	i := 0

	query := fmt.Sprintf("SELECT 1 FROM %v WHERE %v=?", db.Table, db.Column)

	row := db.DB.QueryRow(
		query,
		value,
	)

	err := row.Scan(&i)

	result := err == nil

	if err == sql.ErrNoRows {
		err = nil
	}

	return result, err

}

func Insert(db *DB, value string) error {
	query := fmt.Sprintf("INSERT INTO %v (%v) VALUES (?)", db.Table, db.Column)

	_, err := db.DB.Exec(
		query,
		value,
	)

	return err
}

// func InsertIfItHasNotExist(db *DB, value string) error {
// 	hasExist, err := HasExist(db, value)

// 	if err != nil {
// 		return fmt.Errorf("has exist: %v", err)
// 	}

// 	if hasExist {
// 		return fmt.Errorf(HAS_EXIST_ERROR, value)
// 	}

// 	return Insert(db, value)
// }
