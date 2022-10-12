package RSA

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"
)

type PublicKey struct {
	N *big.Int // modulus
	E *big.Int // public exponent
}

type KeyFile struct {
	PriKey []byte
	H *big.Int
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

func Verify(message string, signature *big.Int, pk *PublicKey) bool {
	check_hash := Encrypt(signature, pk)
	real_hash := Hash(message)
	// fmt.Println("check_hash", check_hash)
	// fmt.Println("real_hash", real_hash)
	return check_hash.Cmp(real_hash) == 0
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