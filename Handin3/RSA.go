package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
)

type PublicKey struct {
	N *big.Int // modulus
	E *big.Int // public exponent
}

type PrivateKey struct {
	N *big.Int
	D *big.Int // private exponent
}

func KeyGen(k int) (*PublicKey, *PrivateKey) {
	for {
		// set E to 3
		E := big.NewInt(3)
		random := rand.Reader

		p, _ := rand.Prime(random, k/2)
		q, _ := rand.Prime(random, k/2)

		p_1 := big.NewInt(0).Sub(p, big.NewInt(1))
		q_1 := big.NewInt(0).Sub(q, big.NewInt(1))

		N := big.NewInt(0).Mul(p, q)
		Fi := big.NewInt(0).Mul(p_1, q_1)

		// fmt.Println(p, p.BitLen())
		// fmt.Println(q, q.BitLen())
		// fmt.Println(N, Fi)

		d := big.NewInt(0)

		d.ModInverse(E, Fi)
		// fmt.Println(d)
		if d.Cmp(big.NewInt(0)) != 0 {
			pri := new(PrivateKey)
			pri.N = N
			pri.D = d

			pub := new(PublicKey)
			pub.E = E
			pub.N = N
			return pub, pri
		}
	}
}

func Encrypt(plaintext *big.Int, pk *PublicKey) *big.Int {
	res := big.NewInt(0).Exp(plaintext, pk.E, pk.N)
	return res
}

func Decrypt(ciphertext *big.Int, pk *PrivateKey) *big.Int {
	res := big.NewInt(0).Exp(ciphertext, pk.D, pk.N)
	return res
}

func EncryptToFile(filename string, plaintext []byte, key_str string) {
	fp, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()

	key, _ := hex.DecodeString(key_str)
	block, _ := aes.NewCipher(key)

	aes_ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := aes_ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(aes_ciphertext[aes.BlockSize:], plaintext)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, aes_ciphertext)
	fp.Write(buf.Bytes())

}

func DecryptFromFile(filename string, key_str string) []byte {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return []byte("DecryptFromFile error")
	}
	defer fp.Close()

	fileinfo, _ := fp.Stat()
	filesize := fileinfo.Size()
	ciphertext := make([]byte, filesize)
	len, _ := fp.Read(ciphertext)

	key, _ := hex.DecodeString(key_str)
	block, _ := aes.NewCipher(key)

	// fmt.Println(ciphertext)
	iv := ciphertext[:aes.BlockSize]
	plaintext2 := make([]byte, len-aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])
	return plaintext2
}

func main() {
	pub, pri := KeyGen(128)
	fmt.Println("RSA PublicKey:", pub)
	fmt.Println("RSA PrivateKey:", pri)
	random := rand.Reader
	plain, _ := rand.Int(random, big.NewInt(2048))
	fmt.Println("plaintext before RSA encryption:", plain)
	ciphertext := Encrypt(plain, pub)
	fmt.Println("RSA ciphertext :", ciphertext)

	aes_plaintext, _ := json.Marshal(pri)
	key := "6368616e676520746869732070617373"
	fmt.Println("AES_CTR key:", key)
	EncryptToFile("bin", aes_plaintext, key)

	decrypted_text := DecryptFromFile("bin", key)

	decrypted_pri := new(PrivateKey)
	json.Unmarshal(decrypted_text, decrypted_pri)
	fmt.Println("AES_decrypted private key :", decrypted_pri)

	plaintext := Decrypt(ciphertext, decrypted_pri)
	fmt.Println("plaintext after RSA decryption:", plaintext)

	// -----------------------------------------

}
