package chain

/**
 *定义区块链结构体，该结构体用于管理区块
 */
type BlockChain struct {
	Blocks []Block
}

/**
  *创建一个区块链对象，包含一个创世区块
 */
func CreatChainWithGensis (data []byte) BlockChain {
	gensis := CreateGenesis(data)
	blocks := make([]Block,0)
	blocks = append(blocks,gensis)
	return BlockChain{blocks}
}

/**
 * 生成一个新区快
 */
func (chain *BlockChain)CreateNewBlock(data []byte) {
	blocks := chain.Blocks//获取到当前所有的区块
	lastBlock := blocks[len(blocks)-1]//最后最新的区块
	newBlock :=NewBlock(lastBlock.Height,lastBlock.Hash,lastBlock.Data)
	chain.Blocks = append(chain.Blocks,newBlock)
}
