package main

import (
	"XianfengChain04/chain"

	"fmt"
)

func main() {
	//fmt.Println("BlockChain")
	/*block0 := chain.Block{
		Height:    0,
		Version:   0x00,
		PrevHash:  [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		//Hash:      [32]byte{},
		TimeStamp: time.Now().Unix(),
		Nonce:     0,
		Data:      nil,
	}*/
	/*block0 := chain.CreateGenesis([]byte("Blockchian"))
	block1 := chain.NewBlock(block0.Height,block0.Hash,[]byte("Hello word"))*/
	blockChain := chain.CreatChainWithGensis([]byte("HelloWord"))
	blockChain.CreateNewBlock([]byte("Hello"))

	fmt.Println("区块链的个数:",len(blockChain.Blocks))

	fmt.Println("区块0:",blockChain.Blocks[0])
	fmt.Println("区块1:",blockChain.Blocks[1])

	firstBlock := blockChain.Blocks[0]
	firstBytes,err := firstBlock.Serialize()
	if err != nil{
		panic(err.Error())
	}
	//反序列化，验证逆过程
	deFirstBlock,err :=chain.Deserialize(firstBytes)
	if err != nil{
		panic(err.Error())
	}
	fmt.Println(string(deFirstBlock.Data))
}
