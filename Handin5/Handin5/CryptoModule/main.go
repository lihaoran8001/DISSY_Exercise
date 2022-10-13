package main

import (
	"CryptoModule/AES"
	"CryptoModule/RSA"
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
	message := "DISSY2022"
	// password = "anything else"
	signature := AES.Sign("911bin", password, []byte(message))
	fmt.Println(signature)
	isVerified := RSA.VerifyMessage(message, signature, pub)
	if isVerified{
		fmt.Println("verify success")
	}else{
		fmt.Println("verify fail")
	}
}
