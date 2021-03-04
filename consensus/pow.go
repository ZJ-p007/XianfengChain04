package consensus

import (
	"XianfengChain04/chain"
	"XianfengChain04/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

const DIFFICULTY  = 10//难度值系数

type PoW struct {
	Block chain.Block
	Target *big.Int
}

/**
 *
 */
func (pow PoW) FindNonce() int64 {
	//fmt.Println("这里是共识机制POW")
	var nonce int64
	nonce = 0

	//无限循环
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
		result := bytes.Compare(hash[:], target.Bytes())
		if result == -1 {
			return nonce
		}
		nonce++ //否则nonce自增
	}
	return 0
}

/**
 *根据区块已有的信息和当前nonce的赋值，计算区块的hash
 */
func CalculateHash(block chain.Block,nonce int64) [32]byte {
	heightByte, _ := utils.Int2Byte(block.Height)
	versionByte, _ := utils.Int2Byte(block.Version)
	timeByte, _ := utils.Int2Byte(block.TimeStamp)
	nonceByte, _ := utils.Int2Byte(nonce)
	blockByte := bytes.Join([][]byte{heightByte, versionByte, block.PrevHash[:], timeByte, nonceByte, block.Data}, []byte{})
	//1、计算区块的哈希
	hash := sha256.Sum256(blockByte)
	return hash
}