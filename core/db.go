package core

import (
	"database/sql"
	"fmt"
	"os/exec"
)

const CLIENT_DB_PATH = "./client.db"
const SERVER_DB_PATH = "./server.db"
const TEST_DB_PATH = ":memory:"

func SetupDB(path string) (*sql.DB, error) {
	if path != TEST_DB_PATH {
		exec.Command("touch", path).Run()
	}

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS "BASENAMES" ("BASENAME" VARCHAR(1024));`,
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func HasExist(db *sql.DB, basename string) (bool, error) {
	i := 0

	row := db.QueryRow(
		`SELECT 1 FROM BASENAMES WHERE BASENAME=?`,
		basename,
	)

	err := row.Scan(&i)

	result := err == nil

	if err == sql.ErrNoRows {
		err = nil
	}

	return result, err

}

func Insert(db *sql.DB, basename string) error {
	_, err := db.Exec(
		`INSERT INTO BASENAMES (BASENAME) VALUES (?)`,
		basename,
	)

	return err
}

func InsertIfItDoesNotExist(db *sql.DB, basename string) error {
	hasExist, err := HasExist(db, basename)

	if err != nil {
		return err
	}

	if hasExist {
		return fmt.Errorf("basename %s already exists", basename)
	}

	return Insert(db, basename)
}
