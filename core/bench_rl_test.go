package core

// type targetFunc = func(b *testing.B, db *DB)

// func prepareDB(conf DB, entityType string, logSize int, value string) {
// 	path := fmt.Sprintf(BENCH_DB_PATH, entityType, logSize)
// 	db, err := SetupDB(conf, path)

// 	if err != nil {
// 		panic(err)
// 	}

// 	countQuery := fmt.Sprintf(`SELECT COUNT (*) FROM %v`, conf.Table)
// 	deleteQuery := fmt.Sprintf("DELETE FROM %v", conf.Table)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer db.DB.Close()

// 	if err != nil {
// 		panic(err)
// 	}

// 	res, err := db.DB.Query(countQuery)

// 	if err != nil {
// 		panic(err)
// 	}

// 	var count int

// 	res.Next()
// 	res.Scan(&count)
// 	res.Next()

// 	if count == logSize {
// 		return
// 	}

// 	if count > logSize {
// 		fmt.Printf("%v != %v\n", count, logSize)

// 		fmt.Println("clear!")

// 		_, err = db.DB.Exec(deleteQuery)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	fmt.Println("insert!")

// 	for i := count; i < logSize; i++ {
// 		err := Insert(db, value)

// 		if err != nil {
// 			panic(err)
// 		}

// 		fmt.Printf("%v/%v\n", i, logSize)
// 	}

// 	db.DB.Close()
// }

// // func benchSign(b *testing.B, db *DB) {
// // 	period := Now()

// // 	h := fmt.Sprintf("%v_%v", BENCH_ORIGIN, period)

// // 	b.StartTimer()
// // 	_, err := Sign(BENCH_ORIGIN, period)
// // 	b.StopTimer()

// // 	if err != nil {
// // 		b.Fatalf("%v: ", err)
// // 	}

// // 	query := fmt.Sprintf("DELETE FROM %v WHERE %v = ?", SIGNER_DB_CONF.Table, SIGNER_DB_CONF.Column)
// // 	_, err = db.DB.Exec(query, h)

// // 	if err != nil {
// // 		b.Fatalf("%v: ", err)
// // 	}
// // }

// // func benchVerify(b *testing.B, db *DB) {
// // 	period := Now()

// // 	signature, err := Sign(BENCH_ORIGIN, period)

// // 	if err != nil {
// // 		b.Fatalf("%v: ", err)
// // 	}

// // 	b.StartTimer()
// // 	// err = Verify(signature, BENCH_ORIGIN, period)

// // 	hasExist, err := HasExist(db, signature)

// // 	if err != nil {
// // 		b.Fatalf("has exist: %v", err)
// // 	}

// // 	if hasExist {
// // 		b.Fatalf(HAS_EXIST_ERROR, signature)
// // 	}

// // 	b.StopTimer()

// // 	if err != nil {
// // 		b.Fatalf("%v: ", err)
// // 	}

// // 	K, err := GetK(signature)

// // 	if err != nil {
// // 		b.Fatalf("%v: ", err)
// // 	}

// // 	query := fmt.Sprintf("DELETE FROM %v WHERE %v = ?", VERIFIER_DB_CONF.Table, VERIFIER_DB_CONF.Column)
// // 	_, err = db.DB.Exec(query, K)

// // 	if err != nil {
// // 		b.Fatalf("%v: ", err)
// // 	}
// // }

// func benchSearchVerifierRL(b *testing.B, db *DB) {
// 	period := Now()

// 	signature, err := Sign(BENCH_ORIGIN, period)

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	b.StartTimer()

// 	err = Verify(signature, BENCH_ORIGIN, period)

// 	hasExist, err := HasExist(db, signature)

// 	if err != nil {
// 		b.Fatalf("has exist: %v", err)
// 	}

// 	if hasExist {
// 		b.Fatalf(HAS_EXIST_ERROR, signature)
// 	}

// 	b.StopTimer()
// }

// func benchmarkOfLogCore(b *testing.B, logSize int, conf DB, path string, target targetFunc) {
// 	db, err := SetupDB(conf, path)
// 	defer db.DB.Close()

// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	b.ResetTimer()
// 	b.StopTimer()

// 	for i := 0; i < b.N; i++ {
// 		target(b, db)
// 	}
// }

// func benchmarkOfLog(b *testing.B, conf DB, entityType string, value string, target targetFunc) {
// 	for i := 1; i < BENCH_MAX; i *= BENCH_UNIT {
// 		fmt.Printf("set up db with %d records\n", i)
// 		prepareDB(conf, entityType, i, value)
// 	}

// 	fmt.Println("Ready!!")
// 	fmt.Println("Start!!")

// 	for i := 1; i < BENCH_MAX; i *= BENCH_UNIT {
// 		name := fmt.Sprintf("log_size(%v):%d", entityType, i)

// 		b.Run(name, func(b *testing.B) {
// 			path := fmt.Sprintf(BENCH_DB_PATH, entityType, i)
// 			benchmarkOfLogCore(b, i, conf, path, target)
// 		})
// 	}
// }

// // sudo go test -benchmem -run=^$ -bench ^BenchmarkOfLog$ example.com/m/v2
// func BenchmarkOfLog(b *testing.B) {
// 	fmt.Println("Ready...")

// 	err := Setup()
// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	h := fmt.Sprintf("%v_%v", BENCH_ORIGIN, Now())

// 	signature, err := Sign(BENCH_ORIGIN, Now())
// 	if err != nil {
// 		b.Fatalf("%v: ", err)
// 	}

// 	K, _ := GetK(signature)

// 	benchmarkOfLog(b, SIGNER_DB_CONF, "signer", h, benchSearchSignerLog)
// 	benchmarkOfLog(b, VERIFIER_DB_CONF, "verifier", K, benchSearchVerifierLog)
// }
