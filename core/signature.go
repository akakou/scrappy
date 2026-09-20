package scrappy

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/ecdaa_helper"
)

// RateLimitK is k in the Scrappy protocol. Signers and verifiers must use
// the same value.
var RateLimitK = 2

func pickUnusedIndex(origin string, period int, db *DB) (int, error) {
	if RateLimitK <= 0 {
		return 0, fmt.Errorf("invalid rate limit k: %d", RateLimitK)
	}

	basename := HashBasename(origin)
	usedIndexes, err := SelectUsedIndexes(db, period, basename)
	if err != nil {
		return 0, fmt.Errorf("select used indexes: %v", err)
	}

	used := make(map[int]struct{}, len(usedIndexes))
	for _, i := range usedIndexes {
		used[i] = struct{}{}
	}

	unused := make([]int, 0, RateLimitK)
	for i := 0; i < RateLimitK; i++ {
		if _, ok := used[i]; !ok {
			unused = append(unused, i)
		}
	}

	if len(unused) == 0 {
		return 0, fmt.Errorf("rate limit reached for %s in period %d", origin, period)
	}

	choice, err := rand.Int(rand.Reader, big.NewInt(int64(len(unused))))
	if err != nil {
		return 0, err
	}

	return unused[int(choice.Int64())], nil
}

func encodeProof(i int, signature string) string {
	return strconv.Itoa(i) + ":" + signature
}

func decodeProof(proof string) (int, string, error) {
	iString, signature, ok := strings.Cut(proof, ":")
	if !ok || signature == "" {
		return 0, "", fmt.Errorf("invalid proof encoding")
	}

	i, err := strconv.Atoi(iString)
	if err != nil {
		return 0, "", fmt.Errorf("invalid proof index: %v", err)
	}

	return i, signature, nil
}

func Sign(origin string, period int, signer ecdaa.Signer) (string, error) {
	db, err := SetupDB(SIGNER_LOG_DB_CONF, SIGNER_LOG_DB_PATH)

	if err != nil {
		return "", err
	}

	defer db.DB.Close()

	if !IsValidPeriod(period) {
		return "", fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	i, err := pickUnusedIndex(origin, period, db)
	if err != nil {
		return "", err
	}

	basename := GetBasenameWithIndex(origin, period, i)
	signature, err := ecdaa_helper.SignWithEncoding(basename, signer)

	if err != nil {
		return "", err
	}

	err = InsertSignerLog(db, i, period, HashBasename(origin))
	if err != nil {
		return "", err
	}

	return encodeProof(i, signature), nil
}

func Verify(proof, origin string, period int, ipk *ecdaa.IPK, rl *ecdaa.RevocationList) error {
	if RateLimitK <= 0 {
		return fmt.Errorf("invalid rate limit k: %d", RateLimitK)
	}

	i, signatureString, err := decodeProof(proof)
	if err != nil {
		return err
	}
	if i < 0 || i >= RateLimitK {
		return fmt.Errorf("invalid proof index %d for rate limit %d", i, RateLimitK)
	}

	signature, err := ecdaa_helper.DecodeSignature(signatureString)

	if err != nil {
		return err
	}

	basename := GetBasenameWithIndex(origin, period, i)

	K := GetKBytes(signature)

	logDB, err := SetupDB(VERIFIER_LOG_DB_CONF, VERIFIER_LOG_DB_PATH)

	if err != nil {
		return err
	}

	defer logDB.DB.Close()

	rlDB, err := SetupDB(VERIFIER_RL_DB_CONF, VERIFIER_RL_DB_PATH)

	if err != nil {
		return err
	}

	defer rlDB.DB.Close()

	if !IsValidPeriod(period) {
		return fmt.Errorf("invalid period %d, but now %d", period, Now())
	}

	hasExist, err := HasExist(logDB, K)

	if err != nil {
		return fmt.Errorf("has exist: %v", err)
	}

	if hasExist {
		return fmt.Errorf(HAS_EXIST_ERROR, signature)
	}

	err = ecdaa.Verify([]byte{}, []byte(basename), signature, ipk, *rl)

	if err != nil {
		return err
	}

	err = Insert(logDB, K)

	return err
}
