package consensus

import (
	"XianfengChain04/chain"
	"fmt"
)

type PoW struct {
	Block chain.Block
}

func (pow PoW) FindNonce() int64 {
	fmt.Println("这里是共识机制POW")
	return 0
}