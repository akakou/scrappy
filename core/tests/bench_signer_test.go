package scrappy_test

import (
	"testing"

	"github.com/akakou/scrappy"
)

func benchmarkSignLog(b *testing.B, logSize int) {
	scrappy.SIGNER_LOG_DB_PATH = prepareDB(scrappy.SIGNER_LOG_DB_CONF, "signer_log", logSize, randomBasename)

	db, err := scrappy.SetupDB(scrappy.SIGNER_LOG_DB_CONF, scrappy.SIGNER_LOG_DB_PATH)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	defer db.DB.Close()

	period := scrappy.Now()
	basename := scrappy.HashBasename(BENCH_ORIGIN)

	b.Run("sign_log", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			usedIndexes, err := scrappy.SelectUsedIndexes(db, period, basename)
			hasExist := len(usedIndexes) > 0

			if err != nil {
				b.Fatalf("has exist: %v", err)
			}

			if hasExist {
				b.Fatalf(scrappy.HAS_EXIST_ERROR, basename)
			}

			err = scrappy.InsertSignerLog(db, 0, period, basename)

			if err != nil {
				b.Fatalf("insert: %v", err)
			}

			b.StopTimer()
			sweepDB(scrappy.SIGNER_LOG_DB_CONF, scrappy.SIGNER_LOG_DB_PATH, basename)
			b.StartTimer()
		}
	})
}

func BenchmarkSignLog(b *testing.B) {
	benchmarkSignLog(b, 0)
	benchmarkSignLog(b, 1000)
}
