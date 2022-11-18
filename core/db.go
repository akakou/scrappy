package core

import (
	"database/sql"
	"fmt"
	"os/exec"
)

var SIGNER_DB_PATH = "./signer.db"
var VERIFIER_DB_PATH = "./verifier.db"
var TEST_DB_PATH = ":memory:"

const HAS_EXIST_ERROR = "k %v already exists"

func SetupDB(path string) (*sql.DB, error) {
	if path != TEST_DB_PATH {
		exec.Command("touch", path).Run()
	}

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS "BASENAMES" ("K" VARCHAR(1024));`,
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func HasExist(db *sql.DB, K string) (bool, error) {
	i := 0

	row := db.QueryRow(
		`SELECT 1 FROM BASENAMES WHERE K=?`,
		K,
	)

	err := row.Scan(&i)

	result := err == nil

	if err == sql.ErrNoRows {
		err = nil
	}

	return result, err

}

func Insert(db *sql.DB, K string) error {
	_, err := db.Exec(
		`INSERT INTO BASENAMES (K) VALUES (?)`,
		K,
	)

	return err
}

func InsertIfItHasNotExist(db *sql.DB, K string) error {
	hasExist, err := HasExist(db, K)

	if err != nil {
		return err
	}

	if hasExist {
		return fmt.Errorf(HAS_EXIST_ERROR, K)
	}

	return Insert(db, K)
}
