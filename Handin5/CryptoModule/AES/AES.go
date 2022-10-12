package AES

import(
	"os"
	"io"
	"fmt"
	"encoding/json"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"bytes"
	"encoding/binary"
	"crypto/rand"
	"CryptoModule/RSA"
)

func Generate(filename string, password string) *RSA.PublicKey{
	pub, pri := RSA.KeyGen(200)
	priBytes, _ := json.Marshal(pri)
	EncryptToFileWithHash(filename, priBytes, password)
	return pub
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

	fileContent, _ := json.Marshal(pri)
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