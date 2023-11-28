package scrappy_test

import (
	"fmt"
	"miracl/core"
	"testing"

	"github.com/akakou/ecdaa"
	"github.com/akakou/mcl_utils"
	"github.com/akakou/scrappy"
	"github.com/akakou/scrappy/ecdaa_helper"
	_ "github.com/mattn/go-sqlite3"
)

const BENCH_DB_PATH = "./%v_%v_bench.db"
const BENCH_ORIGIN = "www.example.com"

var unixtime = 1560000000

var signer *ecdaa.SWSigner
var rng *core.RAND

func sign() (string, string) {
	var err error

	if signer == nil || rng == nil {
		rng = mcl_utils.InitRandom()

		_, signer, err = ecdaa.ExampleInitialize(rng)

		if err != nil {
			panic(err)
		}
	}

	basename := scrappy.GetBasename("www.example.com", unixtime)

	signature, err := ecdaa_helper.SignWithEncoding(basename, signer)

	if err != nil {
		panic(err)
	}

	decoded, err := ecdaa_helper.DecodeSignature(signature)

	if err != nil {
		panic(err)
	}

	unixtime++

	return signature, scrappy.GetKBytes(decoded)
}

// func onlySignature() string {
// 	signature, _ := sign()
// 	return signature
// }

func onlyKBytes() string {
	_, kBytes := sign()
	return kBytes
}

func randomBasename() string {
	now := scrappy.Now()
	if signer == nil || rng == nil {
		rng = mcl_utils.InitRandom()
	}

	origin := mcl_utils.RandomBytes(rng, 32)
	return string(scrappy.GetBasename(string(origin), now))
}

func prepareDB(conf scrappy.DB, entityType string, logSize int, data func() string) string {
	path := fmt.Sprintf(BENCH_DB_PATH, entityType, logSize)
	db, err := scrappy.SetupDB(conf, path)

	if err != nil {
		panic(err)
	}

	countQuery := fmt.Sprintf(`SELECT COUNT (*) FROM %v`, conf.Table)
	deleteQuery := fmt.Sprintf("DELETE FROM %v", conf.Table)

	if err != nil {
		panic(err)
	}

	defer db.DB.Close()

	if err != nil {
		panic(err)
	}

	res, err := db.DB.Query(countQuery)

	if err != nil {
		panic(err)
	}

	var count int

	res.Next()
	res.Scan(&count)
	res.Next()

	if count == logSize {
		return path
	}

	if count > logSize {
		fmt.Printf("%v != %v\n", count, logSize)

		fmt.Println("clear!")

		_, err = db.DB.Exec(deleteQuery)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("insert!")

	for i := count; i < logSize; i++ {
		err = scrappy.Insert(db, data())

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v/%v\n", i, logSize)
	}

	db.DB.Close()

	return path
}

func sweepDB(conf scrappy.DB, path, value string) error {
	db, err := scrappy.SetupDB(conf, path)

	if err != nil {
		panic(err)
	}

	query := fmt.Sprintf("DELETE FROM %v WHERE %v = ?", conf.Table, conf.Column)
	_, err = db.DB.Exec(query, value)

	if err != nil {
		return err
	}

	return nil
}

// sudo go test ./bench -bench ^BenchmarkAll$ -benchtime 20x
func BenchmarkAll(b *testing.B) {
	fmt.Println("Ready...")

	// err := Setup()
	// if err != nil {
	// 	b.Fatalf("%v: ", err)
	// }

	// benchmarkSign(b, 0)

	// for i := 1; i <= 100000; i *= 10 {
	// benchmarkSign(b, i)
	// }

	// benchmarkSign(b, 100000)
	//benchmarkVerify(b, 100000)

	// benchmarkAll(b, SIGNER_LOG_DB_CONF, "signer_log", h, benchSearchSignerLog)
	// benchmarkAll(b, VERIFIER_LOG_DB_CONF, "verifier_log", K, benchSearchVerifierLog)
	// benchmarkAll(b, VERIFIER_RL_DB_CONF, "verifier_revocation", rogueSK, benchVerifierVerifyRevocation)
}
