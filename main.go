package main

import (
  "core"
  "fmt"
  "math/big"
)


func demoFlow() {

  curve := core.Setup()
  // fmt.Println(core.CryptoParams(curve))

  key, public := core.KeyGen(curve)

  message := "to-be-signed"
  var verified bool

  // low level version
  r, s := core.Sign(message, key)
  verified = core.VerifySignature(message, &public, r, s)
  fmt.Println(verified)

  // ASN.1 version
  signature := core.SignASN1(message, key)
  verified = core.VerifySignatureASN1(message, &public, signature)
  fmt.Println(verified)

  // Paillier enc/dec
  p := big.NewInt(11)
  q := big.NewInt(17)

  paillierKey := core.NewPaillierKey(p, q)
  paillierPub := paillierKey.Public()

  fmt.Println(paillierKey)
  fmt.Println(paillierPub)

  N := paillierPub.N
  fmt.Println("N:", N)

  m := big.NewInt(175)
  fmt.Println("m:", m)

  c := core.Encrypt(paillierPub, m)
  fmt.Println("c:", c)

  // d := core.Decrypt(p, q, c)
  d := core.Decrypt(paillierKey, c)
  fmt.Println("d:", d)
}


func main() {

  demoFlow()
}
