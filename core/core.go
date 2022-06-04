package core

import (
  "crypto/ecdsa"
  "crypto/elliptic"
  "crypto/sha256"
  "crypto/rand"
  "math/big"
  "hash"
  "fmt"
  "log"
  "io"
)


func Hash(text string) []byte {

  var h hash.Hash
  h = sha256.New()
  io.WriteString(h, text)

  return h.Sum(nil)
}


func Setup() elliptic.Curve {
  return elliptic.P256()
}


func KeyGen(curve elliptic.Curve) (*ecdsa.PrivateKey, ecdsa.PublicKey) {

  key, err := ecdsa.GenerateKey(curve, rand.Reader)

  if err != nil {
    log.Fatal(err)
  }

  return key, key.PublicKey
}


func SignMessage(text string, key *ecdsa.PrivateKey) (*big.Int, *big.Int) {

  r, s, err := ecdsa.Sign(rand.Reader, key, Hash(text))

  if err != nil {
    log.Fatal(err)
  }

  return r, s
}


func VerifySignature(text string, public *ecdsa.PublicKey, r, s *big.Int) bool {

  return ecdsa.Verify(public, Hash(text), r, s)
}


func DemoFlow() {

  curve := Setup()

  key, public := KeyGen(curve)

  var message string
  message = "to-be-signed"

  // low level version
  r, s := SignMessage(message, key)
  vrf := VerifySignature(message, &public, r, s)

  // ASN.1 version
  sig2, err := ecdsa.SignASN1(rand.Reader, key, Hash(message))
  if err != nil {
    log.Fatal(err)
  }

  var vrf2 bool
  vrf2 = ecdsa.VerifyASN1(&public, Hash(message), sig2)

  // Displays
  fmt.Println(vrf)
  fmt.Println(vrf2)
}
