package paillier

import (
  "crypto/rand"
  "math/big"
  "log"
  "fmt"
)


type Key struct {
  p           *big.Int      // p
  q           *big.Int      // q
  N           *big.Int      // N, should be pq
  M           *big.Int      // N ^ 2
  Gamma       *big.Int      // 1 + N
  totient     *big.Int      // phi(N)
  totientInv  *big.Int      // phi(N) ^ -1 (mod N)
}


type PublicKey struct {
  N     *big.Int            // N
  M     *big.Int            // N ^ 2
  Gamma *big.Int            // 1 + N
}


func NewKey(p *big.Int, q *big.Int) *Key {

  one := big.NewInt(1)

  N := new(big.Int).Mul(p, q)                         // N = pq
  M := new(big.Int).Exp(N, big.NewInt(2), nil)        // N ^ 2
  Gamma := new(big.Int).Add(N, one)                   // 1 + N
  pMinusOne := new(big.Int).Sub(p, one)               // p - 1
  qMinusOne := new(big.Int).Sub(q, one)               // q - 1
  totient := new(big.Int).Mul(pMinusOne, qMinusOne)   // phi(N) = (p - 1)(q - 1)
  totientInv := new(big.Int).ModInverse(totient, N)   // phi(N) ^ -1 (mod N)

  return &Key {
    p:          p,
    q:          q,
    N:          N,
    M:          M,
    Gamma:      Gamma,
    totient:    totient,
    totientInv: totientInv,
  }
}


func (key *Key) Public() *PublicKey {
  return &PublicKey {
    N:      key.N,
    M:      key.M,
    Gamma:  key.Gamma,
  }
}


func (public *PublicKey) Encrypt(message *big.Int) *big.Int {
  r := randInt(public.N)
  rToN := new(big.Int).Exp(r, public.N, public.M) // r ^ N (mod N ^ 2)
  GammaToM := new(big.Int).Exp(public.Gamma, message, public.M) // (1 + N) ^ m (mod N ^ 2)
  cipher := new(big.Int).Mul(GammaToM, rToN) // (1 + N) ^ m * r ^ N
  cipher = cipher.Mod(cipher, public.M) // (1 + N) ^ m * r ^ N (mod N ^ 2)
  return cipher
}


func (public *PublicKey) EncryptWithProof(message *big.Int) (*big.Int, *ZKProof) {
  // TODO: Implement
  r := randInt(public.N)
  rToN := new(big.Int).Exp(r, public.N, public.M) // r ^ N (mod N ^ 2)
  fmt.Println(rToN)


  cipher := big.NewInt(0)

  z := big.NewInt(0)
  u1 := big.NewInt(0)
  u2 := big.NewInt(0)
  u3 := big.NewInt(0)
  e := big.NewInt(0)
  s1 := big.NewInt(0)
  s2 := big.NewInt(0)
  s3 := big.NewInt(0)

  proof := &ZKProof{
    z:  z,
    u1: u1,
    u2: u2,
    u3: u3,
    e:  e,
    s1: s1,
    s2: s2,
    s3: s3,
  }

  return cipher, proof
}


func (key *Key) Decrypt(cipher *big.Int) *big.Int {

  c_hat := new(big.Int).Exp(cipher, key.totient, key.M) // c^ = c ^ phi(N) (mod N ^ 2)
  tmp := new(big.Int).Sub(c_hat, big.NewInt(1)) // c^ - 1
  m_hat := new(big.Int).Div(tmp, key.N) // (c^ - 1) / N

  result := new(big.Int).Mul(m_hat, key.totientInv) // ((c^ -1) / N) * (phi(N) ^ -1 (mod N))
  result = result.Mod(result, key.N)  // ((c^ -1) / N) * (phi(N) ^ -1 (mod N))  (mod N)
  return result
}


func randInt(max *big.Int) *big.Int {
  r, err := rand.Int(rand.Reader, max)

  if err != nil {
    log.Fatal(err)
  }

  return r
}


// Generate prime numbers p, q = (p - 1)/2 with bitlength(p) >= bitLen
func GenerateSafePrimes(bitLength int) (*big.Int, *big.Int) {
  p := new(big.Int)

  count := 0
  for {
    fmt.Println(count)
    count ++
    q, err := rand.Prime(rand.Reader, bitLength - 1)

    if err != nil {
      log.Fatal(err)
    }

    // p = 2 * q + 1 = (q << 1) ^ 1
    p.Lsh(q, 1)
    p.SetBit(p, 0, 1)

    // Miller-Rabin primality test for p failing with probability at most equal
    // to (1/4) ^ 25
    if p.ProbablyPrime(25) {
      return p, q
    }
  }
}
