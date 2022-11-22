package core

import (
	"fmt"
	"testing"
)

const BENCH_DB_PATH = "./%v_%v_bench.db"
const BENCH_ORIGIN = "www.example.com"

const BENCH_UNIT = 1000
const BENCH_MAX = 10000

type targetFunc = func(b *testing.B, db *DB, logSize int)

func prepareDB(conf DB, entityType string, logSize int, value string) {
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
		return
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
		err := Insert(db, value)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v/%v\n", i, logSize)
	}

	db.DB.Close()
}

func prepare(b *testing.B, conf DB, path, value string) error {
	db, err := SetupDB(conf, path)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	query := fmt.Sprintf("DELETE FROM %v WHERE %v = ?", conf.Table, conf.Column)
	_, err = db.DB.Exec(query, value)

	if err != nil {
		return err
	}

	return nil
}

func benchmarkSign(b *testing.B) {
	SIGNER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "signer_log", 0)
	VERIFIER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "verifier_log", 0)
	VERIFIER_RL_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "verifier_revocation", 0)

	period := Now()

	basename := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)
	b.Run("sign", func(b *testing.B) {
		b.ResetTimer()
		b.StopTimer()

		for i := 0; i < b.N; i++ {
			prepare(b, SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, basename)

			b.StartTimer()

			signature, err := Sign(BENCH_ORIGIN, period)

			if err != nil {
				b.Fatalf("%v: ", err)
			}

			b.StopTimer()

			fmt.Printf("signature: %v\n", signature)
			fmt.Printf("signature size: %v\n", len(signature))
		}
	})
}

func benchmarkVerify(b *testing.B) {
	SIGNER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "signer_log", 0)
	VERIFIER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "verifier_log", 0)
	VERIFIER_RL_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "verifier_revocation", 0)

	period := Now()
	basename := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)

	prepare(b, SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, basename)

	signature, err := Sign(BENCH_ORIGIN, period)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	K, err := GetK(signature)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	b.Run("verify", func(b *testing.B) {
		b.ResetTimer()
		b.StopTimer()

		for i := 0; i < b.N; i++ {
			prepare(b, VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH, K)

			b.StartTimer()

			err := Verify(signature, BENCH_ORIGIN, period)

			if err != nil {
				b.Fatalf("%v: ", err)
			}

			b.StopTimer()
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

	// K := "AxXNV9CJnzDdcJJ+Pm6N8rlLY2zRHYI0g78FqTt1iUYC"
	// h := "www.example.com_1668902640"
	// rogueSK := "/9uRzMxx+phPfrU8qvmxuO7HpfEEF2Ol4Uw84n8VNC8="

	benchmarkSign(b)
	benchmarkVerify(b)

	// benchmarkAll(b, SIGNER_LOG_DB_CONF, "signer_log", h, benchSearchSignerLog)
	// benchmarkAll(b, VERIFIER_LOG_DB_CONF, "verifier_log", K, benchSearchVerifierLog)
	// benchmarkAll(b, VERIFIER_RL_DB_CONF, "verifier_revocation", rogueSK, benchVerifierVerifyRevocation)
}
