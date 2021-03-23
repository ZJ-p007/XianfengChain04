package chaincrypto

import (
	"XianfengChain04/utils"
	"crypto/elliptic"
	base58 "github.com"
)

/**
 *新生成一个比特币地址
 */
func NewAddress() (string, error) {
	curve :=elliptic.P256()
	//1、生成私钥
	pri,err := NewPriKey(curve)
	if err != nil{
		return "",err
	}
	//2、获取公钥
	pub := GetPubByPriv(curve,pri)

	//3、对公钥进行sha256
	pubHash := utils.Hash256(pub)

    //4、ripemd160计算
    //ripe := ripemd160.New()
    ripemdPub := utils.HashRipemd160(pubHash)

    //5、添加版本号
    versionPub := append([]byte{0x00},ripemdPub...)

    //6、两次哈希(双哈希)
    firstHash := utils.Hash256(versionPub)
    sencondHash := utils.Hash256(firstHash)

    //7、截取前四个字节作为地址校验位
    check := sencondHash[:4]

    //8、拼接到versionPub后面
    originAddress := append(versionPub,check...)

    //9、base58编码
	return base58.Encode(originAddress),nil
}
