package scrappy

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os/exec"

	"github.com/akakou/ecdaa"
)

var SIGNER_LOG_DB_PATH = "./signer_log.db"
var VERIFIER_LOG_DB_PATH = "./verifier_log.db"
var VERIFIER_RL_DB_PATH = "./verifier_rl.db"

var TEST_DB_PATH = "./test.db"

const HAS_EXIST_ERROR = "%v already exists in Logs"

type DB struct {
	Table  string
	Column string
	DB     *sql.DB
	Size   int
}

var SIGNER_LOG_DB_CONF = DB{
	Table:  "SIGNER_LOG",
	Column: "BASENAME",
	Size:   55 + 1,
}

var VERIFIER_LOG_DB_CONF = DB{
	Table:  "VERIFIER_LOG",
	Column: "K",
	Size:   44 + 1,
}

var VERIFIER_RL_DB_CONF = DB{
	Table:  "VERIFIER_RL",
	Column: "ROGUE_SK",
	Size:   44 + 1,
}

func SetupDB(db DB, path string) (*DB, error) {
	query1 := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %v"+
			"(id INTEGER PRIMARY KEY AUTOINCREMENT, "+
			"%v VARCHAR(%v))",
		db.Table, db.Column, db.Size)

	query2 := fmt.Sprintf(
		"CREATE UNIQUE INDEX IF NOT EXISTS id ON %v(id)",
		db.Table)

	if path != TEST_DB_PATH {
		exec.Command("touch", path).Run()
	}

	_db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, err
	}

	_, err = _db.Exec(query1)

	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	_, err = _db.Exec(query2)

	if err != nil {
		return nil, fmt.Errorf("failed to set up index: %v", err)
	}

	db.DB = _db

	return &db, nil
}

func HasExist(db *DB, value string) (bool, error) {
	var i int

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

func SelectAll(db *DB) ([]string, error) {
	result := []string{}

	query := fmt.Sprintf("SELECT %v FROM %v", db.Column, db.Table)
	rows, err := db.DB.Query(
		query,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var str string

		err := rows.Scan(&str)
		if err != nil {
			return nil, err
		}

		result = append(result, str)
	}

	return result, nil
}

func SelectAllRL(db *DB) (ecdaa.RevocationList, error) {
	base64decoded := [][]byte{}

	rl, err := SelectAll(db)

	if err != nil {
		return nil, err
	}

	for _, r := range rl {
		decode, err := base64.StdEncoding.DecodeString(r)
		if err != nil {
			return nil, err
		}

		base64decoded = append(base64decoded, decode)
	}

	result := ecdaa.DecodeRevocationList(base64decoded)
	return result, nil
}
