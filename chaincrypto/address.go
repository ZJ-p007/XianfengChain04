package chaincrypto

/**
 *新生成一个比特币地址
 */
//func NewAddress() (string, error) {
/*	curve := elliptic.P256()
	//1、生成私钥
	pri, err := NewPriKey(curve)
	if err != nil {
		return "", err
	}
	//2、获取公钥
	pub := GetPubByPriv(curve, pri)

	//3、对公钥进行sha256
	pubHash := utils.Hash256(pub)

	//4、ripemd160计算
	//ripe := ripemd160.New()
	ripemdPub := utils.HashRipemd160(pubHash)

	//5、添加版本号
	versionPub := append([]byte{0x00}, ripemdPub...)

	//6、两次哈希(双哈希)
	firstHash := utils.Hash256(versionPub)
	sencondHash := utils.Hash256(firstHash)

	//7、截取前四个字节作为地址校验位
	check := sencondHash[:4]

	//8、拼接到versionPub后面
	originAddress := append(versionPub, check...)

	//9、base58编码
	return base58.Encode(originAddress), nil*/
//}

/**
 *该函数用于检查地址合法，如果符合地址规范，返回true，
 *如果不符合，返回false
 */
/*func CheckAddress(addr string) bool {
	//1、使用base58对传入的地址进行解码
	reAddBytes := base58.Decode(addr)// versionPubHash  -> check

	if len(reAddBytes) < 4{
		return false
	}

	//2、取出校验位
	reCheck := reAddBytes[len(reAddBytes) - 4:]

	//3、截取得到versionPubHah
	reVersionPubHash := reAddBytes[:len(reAddBytes) - 4]

	//4、versionPubHah双哈希
	reFirstHash := utils.Hash256(reVersionPubHash)
	reSencondHash := utils.Hash256(reFirstHash)

	//5、对双哈希以后的内容进行前四个字节的截取
	check := reSencondHash[:4]

	return bytes.Compare(reCheck,check) == 0

}*/
//1J1suFWuuY9Xzic3NNPXnRtA7TwyrtNnZR
