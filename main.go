package main

import (
	"XianfengChain04/chain"
	"fmt"
	"github.com/bolt"
)

const BLOCKS = "Genesis.db"

func main() {
	//打开数据库文件
	db, err := bolt.Open(BLOCKS, 0600, nil)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close() //xxx.db.lock

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

	blockChain := chain.CreateChain(db)
	//创世区块
	err = blockChain.CreatGenesis([]byte("hello word"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//新增一个区块
	err = blockChain.CreateNewBlock([]byte("hello word"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//
	lastBlock := blockChain.GetLastBlock()
	/*if err !=nil{
		fmt.Println(err.Error())
		return
	}*/
	fmt.Println("最新区块：",lastBlock)

	blocks,err := blockChain.GetAllBlocks()
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	for index,block := range blocks{
		fmt.Printf("第%d个区块:",index)
		fmt.Println(block)
	}

	/*
		blockChain := chain.CreatChainWithGensis([]byte("HelloWord"))
		blockChain.DB = db
		blockChain.CreateNewBlock([]byte("Hello"))
	*/

	/*	fmt.Println("区块链的个数:",len(blockChain.Blocks))
		fmt.Println("区块0:",blockChain.Blocks[0])
		fmt.Println("区块1:",blockChain.Blocks[1])
	*/

	/*firstBlock := blockChain.Blocks[0]
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
	*/
}
