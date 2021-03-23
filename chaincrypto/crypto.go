package chaincrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

/**
 *使用密码学随机生成私钥：椭圆曲线数字签名算法ECDSA
 *ECDSA: elliptic  curve digital signature algorithm
 *ECC：elliptical curve crypto
 */

func NewPriKey(curve elliptic.Curve)(*ecdsa.PrivateKey,error){
	//curve := elliptic.P256()
	/*priv,err :=ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil{
		return nil,err
	}
	return priv,nil*/
	return ecdsa.GenerateKey(curve,rand.Reader)
}

/**
 *根据私钥获得公钥
 */
func GetPubByPriv(curve elliptic.Curve,pri *ecdsa.PrivateKey)[]byte{
	return elliptic.Marshal(curve,pri.X,pri.Y)
}