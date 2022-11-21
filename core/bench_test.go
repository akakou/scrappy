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

// func benchSign(b *testing.B, db *DB) {
// 	period := Now()

// 	h := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)

// 	b.StartTimer()
// 	_, err := Sign(BENCH_ORIGIN, period)
// 	b.StopTimer()

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	query := fmt.Sprintf("DELETE FROM %v WHERE %v = ?", SIGNER_DB_CONF.Table, SIGNER_DB_CONF.Column)
// 	_, err = db.DB.Exec(query, h)

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}
// }

// func benchVerify(b *testing.B, db *DB) {
// 	period := Now()

// 	signature, err := Sign(BENCH_ORIGIN, period)

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	b.StartTimer()
// 	// err = Verify(signature, BENCH_ORIGIN, period)

// 	hasExist, err := HasExist(db, signature)

// 	if err != nil {
// 		b.Fatalf("has exist: %v", err)
// 	}

// 	if hasExist {
// 		b.Fatalf(HAS_EXIST_ERROR, signature)
// 	}

// 	b.StopTimer()

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	K, err := GetK(signature)

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	query := fmt.Sprintf("DELETE FROM %v WHERE %v = ?", VERIFIER_DB_CONF.Table, VERIFIER_DB_CONF.Column)
// 	_, err = db.DB.Exec(query, K)

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}
// }

func benchSearchSignerLog(b *testing.B, db *DB, logSize int) {
	SIGNER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "signer_log", logSize)

	b.StartTimer()

	period := Now()

	basename := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)

	b.StartTimer()
	hasExist, err := HasExist(db, basename)

	if err != nil {
		b.Fatalf("has exist: %v", err)
	}

	if hasExist {
		b.Fatalf(HAS_EXIST_ERROR, basename)
	}

	b.StopTimer()
}

func benchSearchVerifierLog(b *testing.B, db *DB, logSize int) {
	SIGNER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "signer_log", 1)
	VERIFIER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "verifier_log", logSize)

	period := Now()

	basename := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)
	prepare(b, SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH, basename)

	signature, err := Sign(BENCH_ORIGIN, period)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	b.StartTimer()
	// err = Verify(signature, BENCH_ORIGIN, period)

	K, err := GetK(signature)
	if err != nil {
		b.Fatalf("%v: ", err)
	}

	hasExist, err := HasExist(db, K)

	if err != nil {
		b.Fatalf("has exist: %v", err)
	}

	if hasExist {
		b.Fatalf(HAS_EXIST_ERROR, signature)
	}

	b.StopTimer()
}

func benchVerifierVerifyRevocation(b *testing.B, db *DB, logSize int) {
	SIGNER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "signer_log", 1)
	VERIFIER_LOG_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "verifier_log", 1)
	VERIFIER_RL_DB_PATH = fmt.Sprintf(BENCH_DB_PATH, "verifier_revocation", logSize)

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

	prepare(b, VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH, K)

	b.StartTimer()
	err = Verify(signature, BENCH_ORIGIN, period)

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	b.StopTimer()
}

func benchmarkAllCore(b *testing.B, logSize int, conf DB, path string, target targetFunc) {
	db, err := SetupDB(conf, path)
	defer db.DB.Close()

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	b.ResetTimer()
	b.StopTimer()

	for i := 0; i < b.N; i++ {
		target(b, db, logSize)
	}
}

func benchmarkAll(b *testing.B, conf DB, entityType string, value string, target targetFunc) {
	for i := 0; i < BENCH_MAX; i += BENCH_UNIT {
		fmt.Printf("set up db with %d records\n", i)
		prepareDB(conf, entityType, i, value)
	}

	fmt.Println("Ready!!")
	fmt.Println("Start!!")

	for i := 0; i < BENCH_MAX; i += BENCH_UNIT {
		name := fmt.Sprintf("log_size(%v):%d", entityType, i)

		b.Run(name, func(b *testing.B) {
			path := fmt.Sprintf(BENCH_DB_PATH, entityType, i)
			benchmarkAllCore(b, i, conf, path, target)
		})
	}
}

// sudo go test -benchmem -run=^$ -bench ^BenchmarkSign$ example.com/m/v2 -benchtime 20x
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

			_, err := Sign(BENCH_ORIGIN, period)

			if err != nil {
				b.Fatalf("%v: ", err)
			}

			b.StopTimer()
		}
	})
}

// sudo go test -benchmem -run=^$ -bench ^BenchmarkVerify$ example.com/m/v2 -benchtime 20x
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
