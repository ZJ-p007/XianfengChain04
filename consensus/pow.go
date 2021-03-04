package consensus

import (
	"XianfengChain04/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

/**目的：拿到区块的属性数据(属性值)
  *1、通过结构体引用，引用block结构体，然后访问其属性
  *2、 接口
 */

const DIFFICULTY  = 15//难度值系数

type PoW struct {
	Block BlockInterface
	Target *big.Int
}

/**
 *
 */
func (pow PoW) FindNonce() ([32]byte,int64) {
	//fmt.Println("这里是共识机制POW")
	var nonce int64
	nonce = 0

	//无限循环
	hashBig := new(big.Int)
	for {
	/*	heightByte, _ := utils.Int2Byte(pow.Block.Height)
		versionByte, _ := utils.Int2Byte(pow.Block.Version)
		timeByte, _ := utils.Int2Byte(pow.Block.TimeStamp)
		nonceByte, _ := utils.Int2Byte(nonce)
		blockBytes := bytes.Join([][]byte{heightByte, versionByte, pow.Block.PrevHash[:], timeByte, nonceByte, pow.Block.Data}, []byte{})
		//1、计算区块的哈希
		hash := sha256.Sum256(blockBytes)*/
		hash := CalculateHash(pow.Block,nonce)
		/**2、拿到系统的nonce
		  *①难度值可调整 difficulty
		  *②难难度和目标值的关系:难度越大，目标值越小
		     实现思路:big.Int -> 初始一 ->初始一 ->左移
		*/
		target := pow.Target
		//3、比较大小
		hashBig = hashBig.SetBytes(hash[:])
		//result := bytes.Compare(hash[:], target.Bytes())
		result := hashBig.Cmp(target)
		if result == -1 {
			return hash,nonce
		}
		nonce++ //否则nonce自增
	}
	//return 0
}

/**
 *根据区块已有的信息和当前nonce的赋值，计算区块的hash
 */
func CalculateHash(block BlockInterface,nonce int64) [32]byte {
	heightByte, _ := utils.Int2Byte(block.GetHeight())
	versionByte, _ := utils.Int2Byte(block.GetVersion())
	timeByte, _ := utils.Int2Byte(block.GetTimeStamp())
	nonceByte, _ := utils.Int2Byte(nonce)
	prev := block.GetPreHash()
	blockByte := bytes.Join([][]byte{heightByte, versionByte, prev[:], timeByte, nonceByte, block.GetData()}, []byte{})
	//1、计算区块的哈希
	hash := sha256.Sum256(blockByte)
	return hash
}
