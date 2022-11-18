package core

import (
	"fmt"
	"testing"
)

const BENCH_DB = "./%v_bench.db"
const BENCH_ORIGIN = "www.example.com"

const BENCH_UNIT = 10

const BENCH_MAX = 1000000

const TMP_K = "0000000000000000000000000000000000000000000000000000000000000000"

func prepareDB(logSize int) {
	db_name := fmt.Sprintf(BENCH_DB, logSize)
	SIGNER_DB_PATH = db_name

	db, err := SetupVerifierDB(db_name)
	defer db.DB.Close()

	if err != nil {
		panic(err)
	}

	res, err := db.DB.Query(
		`SELECT COUNT (*) FROM BASENAMES`,
	)

	if err != nil {
		panic(err)
	}

	var count int

	res.Scan(&count)

	if count == logSize {
		return
	}

	fmt.Print("Delete DB\n")

	_, err = db.DB.Exec("DROP TABLE basenames")
	if err != nil {
		panic(err)
	}

	db.DB.Close()

	db, err = SetupVerifierDB(db_name)

	if err != nil {
		panic(err)
	}

	for i := 0; i < logSize; i++ {
		err := Insert(db, TMP_K)

		if err != nil {
			panic(err)
		}
	}
}

func benchmarkOfSignerLog(b *testing.B, logSize int) {
	db, err := SetupVerifierDB(SIGNER_DB_PATH)
	defer db.DB.Close()

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	for i := 0; i < logSize; i++ {
		err := Insert(db, TMP_K)

		if err != nil {
			b.Fatalf("%v: ", err)
		}
	}

	for i := 0; i < b.N; i++ {
		b.ResetTimer()
		period := Now()

		b.StartTimer()
		signature, err := Sign(BENCH_ORIGIN, period)
		b.StopTimer()

		if err != nil {
			b.Fatalf("%v: ", err)
		}

		K, err := GetK(signature)

		if err != nil {
			b.Fatalf("%v: ", err)
		}

		_, err = db.DB.Exec("DELETE FROM basenames WHERE K = ?", K)

		if err != nil {
			b.Fatalf("%v: ", err)
		}
	}

}

func BenchmarkOfSignerLog(b *testing.B) {
	fmt.Println("Ready...")

	err := Setup()
	if err != nil {
		b.Fatalf("%v: ", err)
	}

	for i := 1; i < BENCH_MAX; i *= BENCH_UNIT {
		fmt.Printf("set up db with %d records\n", i)
		prepareDB(i)
	}

	fmt.Println("Ready!!")
	fmt.Println("Start!!")

	for i := 1; i < BENCH_MAX; i *= BENCH_UNIT {
		name := fmt.Sprintf("log_size:%d", i)

		b.Run(name, func(b *testing.B) {
			benchmarkOfSignerLog(b, i)

		})
	}
}
