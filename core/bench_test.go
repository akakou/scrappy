package core

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

const BENCH_DB_PATH = "./%v_%v_bench.db"
const BENCH_ORIGIN = "www.example.com"

const BENCH_UNIT = 1000
const BENCH_MAX = 10000

const K = "AxXNV9CJnzDdcJJ+Pm6N8rlLY2zRHYI0g78FqTt1iUYC"
const h = "www.example.com_1668902640"
const rogueSK = "/9uRzMxx+phPfrU8qvmxuO7HpfEEF2Ol4Uw84n8VNC8="

func prepareDB(conf DB, entityType string, logSize int, value string) string {
	path := fmt.Sprintf(BENCH_DB_PATH, entityType, logSize)
	db, err := SetupDB(conf, path)

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
		err = Insert(db, value)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v/%v\n", i, logSize)
	}

	db.DB.Close()

	return path
}

func sweepDB(conf DB, path, value string) error {
	db, err := SetupDB(conf, path)

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

func benchmarkSign(b *testing.B, logSize int) {
	SIGNER_LOG_DB_PATH = prepareDB(SIGNER_LOG_DB_CONF, "signer_log", logSize, h)

	var signature string

	db, err := SetupDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	defer db.DB.Close()

	period := Now()
	hashed_origin := sha256.New().Sum([]byte(BENCH_ORIGIN))
	basename := fmt.Sprintf("%v_%v", hashed_origin, period)

	b.Run("sign_log", func(b *testing.B) {
		// if !IsValidPeriod(period) {
		// 	b.Fatalf("invalid period %d, but now %d", period, Now())
		// }

		for i := 0; i < b.N; i++ {
			hasExist, err := HasExist(db, basename)

			if err != nil {
				b.Fatalf("has exist: %v", err)
			}

			if hasExist {
				b.Fatalf(HAS_EXIST_ERROR, basename)
			}

			err = Insert(db, basename)

			if err != nil {
				b.Fatalf("has exist: %v", err)
			}

			b.StopTimer()
			sweepDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, basename)
			b.StartTimer()
		}
	})

	b.Run("sign_crypto", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			signature, err = CryptoSign(basename)

			if err != nil {
				b.Fatalf("%v: ", err)
			}
		}
	})

	fmt.Printf("basename: %v\n", basename)

	fmt.Printf("signature: %v\n", signature)
	fmt.Printf("encoded signature size: %v\n", len(signature))
}

func benchmarkVerify(b *testing.B, logSize, rlSize int) {
	SIGNER_LOG_DB_PATH = prepareDB(SIGNER_LOG_DB_CONF, "signer_log", 0, h)
	VERIFIER_LOG_DB_PATH = prepareDB(VERIFIER_LOG_DB_CONF, "verifir_log", logSize, K)
	VERIFIER_RL_DB_PATH = prepareDB(VERIFIER_RL_DB_CONF, "verifier_revocation", rlSize, rogueSK)

	logDB, err := SetupDB(VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	rlDB, err := SetupDB(VERIFIER_RL_DB_CONF, VERIFIER_RL_DB_PATH)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	rl, err := SelectAllRL(rlDB)
	if err != nil {
		b.Fatalf("%v: ", err)
	}

	period := Now()
	basename := getBasename(BENCH_ORIGIN, period)

	signature, err := Sign(BENCH_ORIGIN, period)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	sweepDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, basename)

	K, err := GetK(signature)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	b.Run("verify_log", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hasExist, err := HasExist(logDB, K)

			if err != nil {
				b.Errorf("has exist: %v", err)
			}

			if hasExist {
				b.Errorf(HAS_EXIST_ERROR, signature)
			}

			err = Insert(logDB, K)

			if err != nil {
				b.Errorf("insert: %v", err)
			}

			b.StopTimer()
			sweepDB(VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH, K)
			b.StartTimer()
		}
	})

	b.Run("verify_crypto", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err = CryptoVerify(signature, basename, rl)

			if err != nil {
				b.Errorf("verify: %v", err)
			}
		}
	})
}

// sudo go test -benchmem -run=^$ -bench ^BenchmarkAll$ example.com/m/v2 -benchtime 20x
func BenchmarkAll(b *testing.B) {
	fmt.Println("Ready...")

	err := Setup()
	if err != nil {
		b.Fatalf("%v: ", err)
	}

	benchmarkSign(b, 1000)
	// benchmarkVerify(b, 100000, 0)
	// benchmarkVerify(b, 0, 50)

	// benchmarkAll(b, SIGNER_LOG_DB_CONF, "signer_log", h, benchSearchSignerLog)
	// benchmarkAll(b, VERIFIER_LOG_DB_CONF, "verifier_log", K, benchSearchVerifierLog)
	// benchmarkAll(b, VERIFIER_RL_DB_CONF, "verifier_revocation", rogueSK, benchVerifierVerifyRevocation)
}
