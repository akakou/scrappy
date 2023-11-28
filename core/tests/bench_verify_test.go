package scrappy_test

import (
	"fmt"
	"testing"

	"github.com/akakou/ecdaa"
	"github.com/akakou/mcl_utils"
	"github.com/akakou/scrappy"
)

func benchmarkVerifyLog(b *testing.B, logSize int) {
	rng := mcl_utils.InitRandom()

	scrappy.SIGNER_LOG_DB_PATH = prepareDB(scrappy.SIGNER_LOG_DB_CONF, "signer_log", 0, func() string { return "" })
	scrappy.VERIFIER_LOG_DB_PATH = prepareDB(scrappy.VERIFIER_LOG_DB_CONF, "verifir_log", logSize, onlyKBytes)

	logDB, err := scrappy.SetupDB(scrappy.VERIFIER_LOG_DB_CONF, scrappy.VERIFIER_LOG_DB_PATH)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	period := scrappy.Now()
	basename := scrappy.GetBasename(BENCH_ORIGIN, period)

	_, signer, err := ecdaa.ExampleInitialize(rng)

	if err != nil {
		b.Fatalf("%v: ", err)
	}
	signature, err := signer.Sign([]byte{}, []byte(basename), rng)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	sweepDB(scrappy.SIGNER_LOG_DB_CONF, scrappy.SIGNER_LOG_DB_PATH, basename)

	if err != nil {
		b.Fatalf("%v: ", err)
	}

	k2 := onlyKBytes()

	name := fmt.Sprintf("verify_log (%v)", logSize)
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			hasExist, err := scrappy.HasExist(logDB, k2)

			if err != nil {
				b.Errorf("has exist: %v", err)
			}

			if hasExist {
				b.Errorf(scrappy.HAS_EXIST_ERROR, signature)
			}

			err = scrappy.Insert(logDB, k2)

			if err != nil {
				b.Errorf("insert: %v", err)
			}

			b.StopTimer()
			sweepDB(scrappy.VERIFIER_LOG_DB_CONF, scrappy.VERIFIER_LOG_DB_PATH, k2)
			k2 = onlyKBytes()
			b.StartTimer()
		}
	})
}

func BenchmarkVerifyLog(b *testing.B) {
	benchmarkVerifyLog(b, 0)
	benchmarkVerifyLog(b, 100000)
}
