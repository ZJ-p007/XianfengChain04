package chain

import (
	"XianfengChain04/consensus"
	"XianfengChain04/utils"
	"bytes"
	"crypto/sha256"
	"time"
)

const VERSION = 0x00
/**
 *区块结构定义
 */
type Block struct {
	Height int64//高度
	Version int64
	PrevHash [32]byte
	Hash [32]byte
	//默克尔根
	TimeStamp int64
	//Difficulty int64
	Nonce int64
	//区块体
	Data []byte
}

/**
 *计算区块的哈希值并进行赋值
 */
func (block *Block)CalculateBlockHash(){
	heightByte,_ := utils.Int2Byte(block.Height)
	versionByte,_:= utils.Int2Byte(block.Version)
	timeByte,_ := utils.Int2Byte(block.TimeStamp)
	nonceByte,_:= utils.Int2Byte(block.Nonce)

	blockBytes:= bytes.Join([][]byte{heightByte,versionByte,block.PrevHash[:],timeByte,nonceByte,block.Data},[]byte{})
	//为区块的哈希字段赋值
	block.Hash = sha256.Sum256(blockBytes)
	//fmt.Println("区块的哈希值是:",block.Hash)
}

/**
 *生成创世区块的函数
 */
func CreateGenesis(data []byte) Block {
	genesis := Block{
		Height:    0,
		Version:   VERSION,
		PrevHash:  [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		Hash:      [32]byte{},
		TimeStamp: time.Now().Unix(),
		//Nonce:     0,
		Data:      data,
	}
	//todo 计算设置哈希 寻找并设置nonce
	//计算并设置哈希
	genesis.CalculateBlockHash()

	proof := consensus.NewPow(genesis)
	genesis.Nonce = proof.FindNonce()
	return genesis
}

/**
 *生成新区块的功能函数
 */
func NewBlock(height int64,prev [32]byte,data []byte) Block {
	newBlock := Block{
		Height:    height+1,
		Version:   VERSION,
		PrevHash:  prev,
		Hash:      [32]byte{},
		TimeStamp: time.Now().Unix(),
		//Nonce:     0,
		Data:      data,
	}
	//todo 设置哈希 寻找并设置nonce
	//设置区块哈希
	newBlock.CalculateBlockHash()
	proof := consensus.NewPow(newBlock)
	newBlock.Nonce = proof.FindNonce()
	return newBlock
}

