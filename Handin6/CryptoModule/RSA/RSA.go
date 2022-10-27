package RSA

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"math/big"
)

type PublicKey struct {
	N *big.Int // modulus
	E *big.Int // public exponent
}

type KeyFile struct {
	PriKey []byte
	H      *big.Int
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

func Sign(plaintext *big.Int, pk *PrivateKey) *big.Int {
	return Decrypt(plaintext, pk)
}

func SignStr(plaintext []byte, pk string, pub string) string {
	plainBI := big.NewInt(0).SetBytes(plaintext)
	pri := new(PrivateKey)
	json.Unmarshal([]byte(pk), pri)

	pubL := new(PublicKey)
	json.Unmarshal([]byte(pub), pubL)

	// fmt.Println("SIGNATURE Bytes", plaintext)
	// fmt.Println("plainBI", plainBI)
	// SigBI := Decrypt(plainBI, pri)
	// fmt.Println("sigBI:", SigBI)
	// deBI := Encrypt(SigBI, pubL)
	// fmt.Println("deBI:", deBI)
	return Decrypt(plainBI, pri).String()
}

func Verify(message string, signature *big.Int, pk *PublicKey) bool {
	check_hash := Encrypt(signature, pk)
	real_hash := Hash(message)
	// fmt.Println("check_hash", check_hash)
	// fmt.Println("real_hash", real_hash)
	return check_hash.Cmp(real_hash) == 0
}

// func VerifyTrans(sig string, pk string) []byte {
// 	sigBI, _ := new(big.Int).SetString(sig, 10)
// 	pub := new(PublicKey)
// 	json.Unmarshal([]byte(pk), pub)

// 	Body := Encrypt(sigBI, pub).Bytes()
// 	fmt.Println("VERIFY BI", Encrypt(sigBI, pub))
// 	fmt.Println("VERIFY Btyes", Body)
// 	// return check_hash.Cmp(real_hash) == 0
// 	return Body
// }

// 添加 hash string
// kwz 改
func VerifyTrans(hash string, sig string, pk string) bool {
	sigBI, _ := new(big.Int).SetString(sig, 10)
	pub := new(PublicKey)
	json.Unmarshal([]byte(pk), pub)
	return VerifyMessage(hash, sigBI, pub)
}

func VerifyMessage(message string, signature *big.Int, pk *PublicKey) bool {
	sigMessage := string(Encrypt(signature, pk).Bytes())
	return sigMessage == message
}

func Hash(message string) *big.Int {
	bytes := HashRaw(message)
	return big.NewInt(0).SetBytes(bytes)
}

func HashRaw(message string) []byte {
	h := sha256.New()
	h.Write([]byte(message))
	return h.Sum(nil)
}
