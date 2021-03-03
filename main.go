package main

import (
	"XianfengChain04/chain"
	"fmt"
)

func main() {
	fmt.Println("0000000")
	/*block0 := chain.Block{
		Height:    0,
		Version:   0x00,
		PrevHash:  [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		//Hash:      [32]byte{},
		TimeStamp: time.Now().Unix(),
		Nonce:     0,
		Data:      nil,
	}*/
	block0 := chain.CreateGenesis([]byte("HelloWord"))
	block1 := chain.NewBlock(block0.Height,block0.Hash,[]byte("Hello word"))

	fmt.Println(block0)
	fmt.Println(block1)
}
