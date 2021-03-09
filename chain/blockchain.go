package chain

import (
	"errors"
	"github.com/bolt"
)

const BLOCKS = "blocks"
const LASTHASH = "lasthash"

/**
 *定义区块链结构体，该结构体用于管理区块
 */
type BlockChain struct {
	//Blocks []Block
	DB *bolt.DB
	LastBlock Block//最新最后的区块
}

func CreateChain(db *bolt.DB) BlockChain {
	return BlockChain{DB:db}
}

/**
 *创建一个区块链对象，包含一个创世区块
 */
func (chain *BlockChain) CreatGenesis(data []byte) error {
	/*genesis := CreateGenesis(data)
	genSerBytes,err :=gensis.Serialize()*/
	var err error
	//gensis持久化到db中去
	engine := chain.DB
	engine.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil { //没有桶
			bucket, err = tx.CreateBucket([]byte(BLOCKS))
			if err != nil {
				return err
				//panic("操作区块存储文件失败，请重试")
			}
		}
		//先查看
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) == 0 {
			genesis := CreateGenesis(data)
			genSerBytes, _ := genesis.Serialize()
			//bucket已经存在
			bucket.Put(genesis.Hash[:], genSerBytes) //把创世区块保存到boltdb中
			//使用一个标志用来记录最新区块的hash，以标明当前文件中存储到了最新的哪个区块
			bucket.Put([]byte(LASTHASH), genesis.Hash[:])
			//把genesis赋值给chain的lastBlock
			chain.LastBlock = genesis
		} else {
			//lasthash有值，长度不为0，什么都不干

			//从文件中读取出最新的区块，并赋值给内存中的chain中的LastBlock
			lastHash := bucket.Get([]byte(LASTHASH))
			lastBlockBytes := bucket.Get(lastHash)
			//把反序列化的最后最新区块赋值给chain.LastBlock
			chain.LastBlock,err = Deserialize(lastBlockBytes)
		}
		return nil
	})
	/*blocks := make([]Block,0)
	blocks = append(blocks,gensis)*/
	//return BlockChain{blocks}
	return err
}

/**
 * 生成一个新区快
 */
func (chain *BlockChain) CreateNewBlock(data []byte) error {
	/**
	 *目的：生成一个新区块，并存到bolt.DB中（持久化）
	 *手段(步骤)：
	     a、从文件中查到当前存储的最新区块数据
	     b、反序列化得到区块
		 c、根据获取的最新区块生成一个新区块
		 d、将最新区块序列化，得到序列化数据
	     e、将序列化数据存储到文件，同时更新最新区块的标记lasthash，更新为最新区块的hash
	*/

	//1、从文件中查到当前存储的最新区块数据
	lastBlock := chain.LastBlock
	//var lastBlock Block
	/*
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil { //桶为空
			err = errors.New("区块数据库操作失败，请重试！")
			return err
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		lastBlockBytes := bucket.Get(lastHash)
		//2、反序列化得到区块
		lastBlock, err = Deserialize(lastBlockBytes)
		if err != nil {
			return err
		}
		return nil
	})
	*/
	//lastBlock := chain.LastBlock

	//3、
	var err error
	newBlock := NewBlock(lastBlock.Height, lastBlock.Hash, data)
	//4、将最新区块序列化，得到序列化数据
	newBlockSerBytes, err := newBlock.Serialize()
	if err != nil {
		return err
	}
	//5、将序列化数据存储到文件，同时更新最新区块的标记lasthash，更新为最新区块的hash
	db := chain.DB
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块数据库操作失败，请重试！")
		}
		//将新生成的区块保存到文件中去
		bucket.Put(newBlock.Hash[:], newBlockSerBytes)
		//同时更新最新区块的标记lasthash
		bucket.Put([]byte(LASTHASH), newBlock.Hash[:])
		//更新内存中的blockchain的lastBlock
		chain.LastBlock = newBlock
		return nil
	})

	/*blocks := chain.Blocks//获取到当前所有的区块
	lastBlock := blocks[len(blocks)-1]//最后最新的区块
	newBlock :=NewBlock(lastBlock.Height,lastBlock.Hash,lastBlock.Data)
	chain.Blocks = append(chain.Blocks,newBlock)*/
	return err
}

/**
 *获取最新区块数据
 *获取所有区块数据
 */

//获取最新区块数据
func (chain *BlockChain) GetLastBlock() Block{
	return chain.LastBlock
}

/*func (chain *BlockChain) GetLastBlock() (Block, error) {
	db := chain.DB
	var err error
	var lastBlock Block
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据库操作失败，请重试！")
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, err = Deserialize(lastBlockBytes)
		if err != nil {
			return err
		}
		return nil
	})
	return lastBlock, err
}*/


//获取所有区块数据
func (chain *BlockChain) GetAllBlocks() ([]Block,error) {
	/**
	 *目的：获取所有的区块
	 *具体步骤：
	    a、找到最后一个区块
	    b、根据最后一个区块依次往前找
	    c、每次找到区块放入到一个[]Block容器中
	    d、找到最开始的创世区块时，结束，不再找了
	 */

	//1、找到最后一个区块
	db := chain.DB
	var err error
	//var lastBlock Block
	blocks := make([]Block,0)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil{
			 err = errors.New("区块数据库操作失败，请重试！")
			return err
		}
		var currentHash []byte   //当前遍历到hash
		currentHash = bucket.Get([]byte(LASTHASH))
		//lastHash := bucket.Get([]byte(LASTHASH))
		//2、根据最后一个区块依次往前找
		for{
		currentBlockBytes := bucket.Get(currentHash)
		currentBlock,err := Deserialize(currentBlockBytes)
		if err != nil{
			//lastBlockBytes := bucket.Get(lastHash)
			//return err
			break
		}
		//3、每次找到区块放入到一个[]Block容器中
		blocks = append(blocks,currentBlock)
		//4、找到最开始的创世区块时，结束，不再找了
		if currentBlock.Height == 0{
			break
		}
		currentHash = currentBlock.PrevHash[:]
		}
		return nil
	})
	return blocks,err
}