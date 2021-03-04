package consensus

import (
	"XianfengChain04/chain"
	"math/big"
)

type Consensus interface {
	FindNonce() int64
}

func NewPow(block chain.Block) Consensus {
	initTarget := big.NewInt(1)
	initTarget.Lsh(initTarget, 255 - DIFFICULTY)
	return PoW{block,initTarget}

}

func NewPos(block chain.Block) PoS {
	return PoS{Block:block}
}
