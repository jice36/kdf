package main

import (
	"fmt"
	"kdf/internal/checksum"
	"kdf/internal/kdf"
	"kdf/internal/logger"
	"math/rand"
	"time"
)

func main() {
	checksum.InitCheckSum("main")
	go checksum.PeriodicCheckSum("main")

	str, t := timeTest("test gen key 100000 ")
	k := kdf.New()
	log := logger.NewLogger()

	S := make([]byte, 16)
	T := make([]byte, 16)

	for i := 0; i < 1000; i++ {
		randSlice(S)
		randSlice(T)
		res, err := k.KDF(S, T)
		fmt.Println(res)
		if err != nil {
			log.Fatalf("kdf: %v", err)
		}
	}

	defer func() {
		testing(str, t)
	}()

}

func timeTest(test string) (string, time.Time) {
	return test, time.Now()
}

func testing(test string, start time.Time) {
	fmt.Println(test, time.Since(start).Seconds())
}

func randSlice(s []byte) {
	i := 0
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 255
	for i < len(s) {
		s[i] = byte(rand.Intn(max-min+1) + min)
		i++
	}
}
