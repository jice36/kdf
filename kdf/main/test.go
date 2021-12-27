package main

import (
	"fmt"
	supporting "github.com/jice36/cipher_support"
	"kdf"
	"log"
	"math/rand"
	"os"
	"time"
)

var logger *log.Logger
var f *os.File

func main(){
	str, t := timeTest("test gen key 100000 " )
	k := kdf.New()
	logger, f = supporting.CreateLogger(logger)
	S := make([]byte, 16)
	T := make([]byte, 16)
	for i := 0; i < 100000; i++ {
		randSlice(S)
		randSlice(T)
		res, err := k.KDF(S, T)
		fmt.Println(res)
		if err != nil{
			logger.Println(err)
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
	fmt.Println( test, time.Since(start).Seconds())
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