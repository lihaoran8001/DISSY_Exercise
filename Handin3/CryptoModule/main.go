package main

import (
	"time"
	"crypto/rand"
	"math/big"
	"encoding/json"
	"fmt"
	"CryptoModule/RSA"
	"CryptoModule/AES"
	"CryptoModule/randString"
)

func main() {
	pub, pri := RSA.KeyGen(2000)
	// fmt.Println("RSA PublicKey:", pub)
	// fmt.Println("RSA PrivateKey:", pri)

	// create a random integer as plaintext
	plain, _ := rand.Int(rand.Reader, big.NewInt(204800))
	fmt.Println("plaintext before RSA encryption:", plain)

	// use RSA public key to encrypt the plaintext
	ciphertext := RSA.Encrypt(plain, pub)
	fmt.Println("RSA ciphertext :", ciphertext)

	// use AES to encrypt the RSA private key
	aes_plaintext, _ := json.Marshal(pri)
	key := "6368616e676520746869732070617373"
	// fmt.Println("AES_CTR key:", key)
	AES.EncryptToFile("bin", aes_plaintext, key)

	decrypted_text := AES.DecryptFromFile("bin", key)

	decrypted_pri := new(RSA.PrivateKey)
	json.Unmarshal(decrypted_text, decrypted_pri)
	// fmt.Println("AES_decrypted private key :", decrypted_pri)

	// used Decrypted RSA private key to get plaintext
	plaintext := RSA.Decrypt(ciphertext, decrypted_pri)
	fmt.Println("plaintext after RSA decryption:", plaintext)


	// generate and verify the signature
	message := "hello world!"
	// message2 := "goodbye world!"
	hashValue := RSA.Hash(message)
	signature := RSA.Sign(hashValue, pri)
	result := RSA.Verify(message, signature, pub)
	// result := RSA.Verify(message2, signature, pub)
	if (result){
		fmt.Println("Verify signature success!")
	}else{
		fmt.Println("Verify signature fail!")
	}


	// measure the time of hashing
	var startTime time.Time
	var duration time.Duration
	for i := 0; i < 10; i++{
		big_string := randString.RandStringRunes(10240)
		startTime = time.Now()
		RSA.HashRaw(big_string)
		duration += time.Since(startTime)
	}
	fmt.Println("Time used for hasing 10kb data:", duration/10)


	// measure the time of RSA signature
	var startTime2 time.Time
	var duration2 time.Duration
	for i := 0; i < 10; i++{
		startTime2 = time.Now()
		RSA.Sign(hashValue, pri)
		duration2 += time.Since(startTime2)
	}
	fmt.Println("Time used for sign with 2000bits keylen:", duration2/10)

	//
}
