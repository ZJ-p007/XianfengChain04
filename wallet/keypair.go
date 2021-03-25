package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

/**
 *地址所对应得秘钥对(私钥+公钥) 封装在一个自定义得结构体中
 */
type KeyPair struct {
	Pri *ecdsa.PrivateKey
	Pub []byte
}

/**
 *生成一对钥对
 */
func NewKeyPair() (*KeyPair,error) {
	curve := elliptic.P256()
	pri,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil{
		return nil,err
	}
	pub :=elliptic.Marshal(curve,pri.X,pri.Y)

	keyPair := KeyPair{
		Pri: pri,
		Pub: pub,
	}
	return &keyPair,nil
}