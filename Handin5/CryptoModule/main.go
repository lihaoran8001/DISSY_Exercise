package main

import (
	"CryptoModule/AES"
	"CryptoModule/randString"

	// "CryptoModule/RSA"
	// "CryptoModule/randString"
	// "crypto/rand"
	// "encoding/json"
	"fmt"
	// "math/big"
	// "time"
)

func main() {
	// --------- Exercise 9.11 --------------//
	password := randString.RandStringNums(32)
	pub := AES.Generate("911bin", password)
	fmt.Println(pub.E)

	signature := AES.Sign("911bin", password, []byte("123"))
	fmt.Println(signature)
}
