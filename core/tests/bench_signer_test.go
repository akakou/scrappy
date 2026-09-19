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
	basename := scrappy.GetBasename(BENCH_ORIGIN, period)

	b.Run("sign_log", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hasExist, err := scrappy.HasExist(db, basename)

			if err != nil {
				b.Fatalf("has exist: %v", err)
			}

			if hasExist {
				b.Fatalf(scrappy.HAS_EXIST_ERROR, basename)
			}

			err = scrappy.Insert(db, basename)

			if err != nil {
				b.Fatalf("has exist: %v", err)
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
