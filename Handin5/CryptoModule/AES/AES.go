package AES

import (
	"CryptoModule/RSA"
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

func Generate(filename string, password string) *RSA.PublicKey {
	pub, pri := RSA.KeyGen(2000)
	priBytes, _ := json.Marshal(pri)
	fmt.Println("before", priBytes)
	EncryptToFileWithHash(filename, priBytes, password)
	return pub
}

func Sign(filename string, password string, msg []byte) (Signature []byte) {
	priKey_bytes := DecryptFromFileWithHash(filename, password)
	if string(priKey_bytes) == string([]byte("Wrong password!")) {
		return []byte("Wrong password!")
	}

	priKey := new(RSA.PrivateKey)
	fmt.Println("after", priKey_bytes)
	json.Unmarshal(priKey_bytes, priKey)
	msg_bytes := big.NewInt(0).SetBytes(msg)
	fmt.Println(priKey)
	return RSA.Sign(msg_bytes, priKey).Bytes()

}

func DecryptFromFileWithHash(filename string, password string) []byte {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return []byte("DecryptFromFile error")
	}
	defer fp.Close()

	KF := new(RSA.KeyFile)

	fileinfo, _ := fp.Stat()
	filesize := fileinfo.Size()
	ciphertext := make([]byte, filesize)
	len, _ := fp.Read(ciphertext)

	json.Unmarshal(ciphertext, KF)
	fmt.Println(RSA.Hash(password))
	fmt.Println(KF.H)
	if RSA.Hash(password).Cmp(KF.H) != 0 {
		return []byte("Wrong password!")
	}

	key, _ := hex.DecodeString(password)
	block, _ := aes.NewCipher(key)

	// fmt.Println(ciphertext)
	iv := KF.PriKey[:aes.BlockSize]
	plaintext2 := make([]byte, len-aes.BlockSize)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext2, KF.PriKey[aes.BlockSize:])
	return plaintext2
}

func EncryptToFileWithHash(filename string, pri []byte, password string) {
	fp, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()
	// convert password to bytes and use it as AES key
	key, _ := hex.DecodeString(password)
	block, _ := aes.NewCipher(key)

	aes_ciphertext := make([]byte, aes.BlockSize+len(pri))
	iv := aes_ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(aes_ciphertext[aes.BlockSize:], pri)
	// Use password to encrypt private key
	KF := new(RSA.KeyFile)
	KF.PriKey = aes_ciphertext
	// hash password
	KF.H = RSA.Hash(password)

	fileContent, _ := json.Marshal(KF)
	// save encrypted privateKey and hash of password to file
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, fileContent)
	fp.Write(buf.Bytes())

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
