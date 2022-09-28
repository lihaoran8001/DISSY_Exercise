package AES

import(
	"os"
	"io"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"bytes"
	"encoding/binary"
	"crypto/rand"
)

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