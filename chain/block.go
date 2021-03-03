package chain

import "time"

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


	return newBlock
}

