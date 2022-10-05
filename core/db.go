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
		`CREATE TABLE IF NOT EXISTS "BASENAMES" ("ORIGIN" VARCHAR(1024), "PERIOD" INTEGER);`,
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func HasExist(db *sql.DB, origin string, period int) (bool, error) {
	i := 0

	row := db.QueryRow(
		`SELECT 1 FROM BASENAMES WHERE ORIGIN=? AND PERIOD=?`,
		origin,
		period,
	)

	err := row.Scan(&i)

	result := err == nil

	if err == sql.ErrNoRows {
		err = nil
	}

	return result, err

}

func Insert(db *sql.DB, origin string, period int) error {
	_, err := db.Exec(
		`INSERT INTO BASENAMES (ORIGIN, PERIOD) VALUES (?, ?)`,
		origin,
		period,
	)

	return err
}

func InsertIfItDoesNotExist(db *sql.DB, origin string, period int) error {
	hasExist, err := HasExist(db, origin, period)

	if err != nil {
		return err
	}

	if hasExist {
		return fmt.Errorf("basename %s with period %d already exists", origin, period)
	}

	return Insert(db, origin, period)
}
