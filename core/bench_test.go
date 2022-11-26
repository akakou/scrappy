package core

import (
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
		err := InsertHash(db, value)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%v/%v\n", i, logSize)
	}

	db.DB.Close()

	return path
}

func sweepDB(b *testing.B, conf DB, path, value string) error {
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

func benchmarkSign(b *testing.B, logSize int) {
	SIGNER_LOG_DB_PATH = prepareDB(SIGNER_LOG_DB_CONF, "signer_log", logSize, h)

	period := Now()

	var signature string
	var err error

	basename := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)
	b.Run("sign", func(b *testing.B) {
		b.ResetTimer()
		b.StopTimer()

		for i := 0; i < b.N; i++ {
			b.StartTimer()

			signature, err = Sign(BENCH_ORIGIN, period)

			if err != nil {
				b.Fatalf("%v: ", err)
			}

			b.StopTimer()

			sweepDB(b, SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, basename)
		}

		fmt.Printf("signature: %v\n", signature)
		fmt.Printf("encoded signature size: %v\n", len(signature))

	})
}

func benchmarkVerify(b *testing.B, logSize, rlSize int) {
	SIGNER_LOG_DB_PATH = prepareDB(SIGNER_LOG_DB_CONF, "signer_log", 0, h)
	VERIFIER_LOG_DB_PATH = prepareDB(VERIFIER_LOG_DB_CONF, "verifir_log", logSize, K)
	VERIFIER_RL_DB_PATH = prepareDB(VERIFIER_RL_DB_CONF, "verifier_revocation", rlSize, rogueSK)

	period := Now()
	basename := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)

	sweepDB(b, SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, basename)

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
			b.StartTimer()

			err := Verify(signature, BENCH_ORIGIN, period)

			if err != nil {
				b.Fatalf("%v: ", err)
			}

			b.StopTimer()
			sweepDB(b, VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH, K)
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
	benchmarkVerify(b, 10000, 0)
	benchmarkVerify(b, 0, 50)

	// benchmarkAll(b, SIGNER_LOG_DB_CONF, "signer_log", h, benchSearchSignerLog)
	// benchmarkAll(b, VERIFIER_LOG_DB_CONF, "verifier_log", K, benchSearchVerifierLog)
	// benchmarkAll(b, VERIFIER_RL_DB_CONF, "verifier_revocation", rogueSK, benchVerifierVerifyRevocation)
}
