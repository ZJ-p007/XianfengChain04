package wallet

import (
	"XianfengChain04/utils"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	base58 "github.com"
	"github.com/bolt"
)

const KEYSTORE = "keystore"
const ADDANDPRIR = "addrs_keypairs"

/**
 *定义wallet结构体，用于管理地址和对应得秘钥对信息
 */
type Wallet struct {
	Address map[string]*KeyPair
	Engine  *bolt.DB
}

func (wallet *Wallet) NewAddress() (string, error) {
	keyPair, err := NewKeyPair()
	if err != nil {
		return "", err
	}
	//3、对公钥进行sha256
	pubHash := utils.Hash256(keyPair.Pub)

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
	address, err := base58.Encode(originAddress), nil
	if err != nil {
		return "", err
	}
	//把新生成得地址和对应得秘钥存入到wallet中得map中管理起来
	wallet.Address[address] = keyPair //仅仅是在内存中

	//把更新了地址信息和对应秘钥对的map结构中的数据持久化存在文件中
	wallet.SaveAddAndKeyPairs2DB()
	return address, nil

}

/**
 *地址校验
 */
func (wallet *Wallet) CheckAddress(addr string) bool {
	//1、使用base58对传入的地址进行解码
	reAddBytes := base58.Decode(addr) // versionPubHash  -> check

	if len(reAddBytes) < 4 {
		return false
	}

	//2、取出校验位
	reCheck := reAddBytes[len(reAddBytes)-4:]

	//3、截取得到versionPubHah
	reVersionPubHash := reAddBytes[:len(reAddBytes)-4]

	//4、versionPubHah双哈希
	reFirstHash := utils.Hash256(reVersionPubHash)
	reSencondHash := utils.Hash256(reFirstHash)

	//5、对双哈希以后的内容进行前四个字节的截取
	check := reSencondHash[:4]

	return bytes.Compare(reCheck, check) == 0

}

/**
 *该方法用于将内存中的map数据的地址和秘钥对保存到持久化文件中
 */
func (wallet *Wallet) SaveAddAndKeyPairs2DB() {
	var err error
	wallet.Engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(KEYSTORE))
			if err != nil {
				return err
			}
		}
		//桶keystores已存在,可以向桶中存放map的数据
		//addBytes,_ := utils.Encode(wallet.Address)
		gob.Register(elliptic.P256())
		buff := new(bytes.Buffer)
		encoder := gob.NewEncoder(buff)
		err := encoder.Encode(wallet.Address)
		if err != nil {
			return err
		}
		bucket.Put([]byte(ADDANDPRIR), buff.Bytes())
		return nil
	})
}

/**
 *从文件中读取已经存在的地址和对应的秘钥对信息
 */
func LoadAddrAndKeyPairsFromDB(engine *bolt.DB) (*Wallet, error) {
	address := make(map[string]*KeyPair)
	var err error
	engine.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(KEYSTORE))
		if bucket == nil {
			return nil
		}
		//如果有keystore存在，从keystore桶中读取
		addsAndKeyPairsBytes := bucket.Get([]byte(ADDANDPRIR))
		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(addsAndKeyPairsBytes))
		err = decoder.Decode(&address)
		return err
	})
	if err != nil {
		return nil, err
	}
	walet := &Wallet{
		Address: address,
		Engine:  engine,
	}
	return walet, err
}