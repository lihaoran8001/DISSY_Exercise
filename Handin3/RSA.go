package main

import(
	"math/big"
	"crypto/rand"
	"fmt"
)

type PublicKey struct {
	N *big.Int // modulus
	E *big.Int // public exponent
}

type PrivateKey struct{
	N *big.Int
	D *big.Int // private exponent
}

func KeyGen(k int)(*PublicKey, *PrivateKey){
	for{
		// set E to 3
		E := big.NewInt(3)
		random := rand.Reader

		p, _ := rand.Prime(random, k / 2)
		q, _ := rand.Prime(random, k / 2)

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
		if (d.Cmp(big.NewInt(0)) != 0){
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

func Encrypt(plaintext *big.Int, pk *PublicKey)(*big.Int){
	res := big.NewInt(0).Exp(plaintext, pk.E, pk.N)
	return res
}

func Decrypt(ciphertext *big.Int, pk *PrivateKey)(*big.Int){
	res := big.NewInt(0).Exp(ciphertext, pk.D, pk.N)
	return res
}

func main(){
	pub, pri := KeyGen(128)
	fmt.Println("PublicKey:", pub)
	fmt.Println("PrivateKey:", pri)
	random := rand.Reader
	plain, _ := rand.Int(random, big.NewInt(2048))
	fmt.Println("plaintext before encryption:", plain)
	ciphertext := Encrypt(plain, pub)
	fmt.Println("ciphertext:", ciphertext)
	plaintext := Decrypt(ciphertext, pri)
	fmt.Println("plaintext after decryption:", plaintext)
}
