package chain

import (
	"XianfengChain04/transaction"
	"errors"
	"fmt"

	"github.com/bolt"
	"math/big"
)

const BLOCKS = "blocks"
const LASTHASH = "lasthash"

/**
 *定义区块链结构体，该结构体用于管理区块
 */
type BlockChain struct {
	//Blocks []Block
	DB                *bolt.DB
	LastBlock         Block    //最新最后的区块
	IteratorBlockHash [32]byte //表示当前迭代到了哪个区块，改变了用于记录跌倒到的区块hash
}

func CreateChain(db *bolt.DB) BlockChain {
	var lastBlock Block
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(BLOCKS))
		}
		lastHash := bucket.Get([]byte(LASTHASH))
		if len(lastHash) <= 0 {
			return nil
		}
		lastBlockBytes := bucket.Get(lastHash)
		lastBlock, _ = Deserialize(lastBlockBytes)
		return nil
	})
	return BlockChain{
		DB:                db,
		LastBlock:         lastBlock,
		IteratorBlockHash: lastBlock.Hash,
	}
}

/**
*创建coinbaase交易的方法
 */
func (chain *BlockChain) CreateCoinBase(addr string) error {
	//1、创建一笔coinbase交易
	coinbase, err := transaction.CreateCoinBase(addr)
	if err != nil {
		//fmt.Println("创建coinbase交易遇到错误")
		return err
	}
	//2、把coinbase交易到区块中
	err = chain.CreatGenesis([]transaction.Transaction{*coinbase})
	return err
}

/**
 *创建一个区块链对象，包含一个创世区块
 */
func (chain *BlockChain) CreatGenesis(txs []transaction.Transaction) error {
	/*genesis := CreateGenesis(data)
	genSerBytes,err :=gensis.Serialize()*/
	hashBig := new(big.Int)
	hashBig.SetBytes(chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 { //最新区块hash有值，则说明创世区块已存在
		return nil
	}

	var err error
	//genesis持久化到db中去
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
			genesis := CreateGenesis(txs)
			genSerBytes, _ := genesis.Serialize()
			//bucket已经存在
			bucket.Put(genesis.Hash[:], genSerBytes) //把创世区块保存到bolt.db中
			//使用一个标志用来记录最新区块的hash，以标明当前文件中存储到了最新的哪个区块
			bucket.Put([]byte(LASTHASH), genesis.Hash[:])
			//把genesis赋值给chain的lastBlock
			chain.LastBlock = genesis
			chain.IteratorBlockHash = genesis.Hash
		} /* else {
			//lasthash有值，长度不为0，什么都不干
			//从文件中读取出最新的区块，并赋值给内存中的chain中的LastBlock
			lastHash := bucket.Get([]byte(LASTHASH))
			lastBlockBytes := bucket.Get(lastHash)
			//把反序列化的最后最新区块赋值给chain.LastBlock
			chain.LastBlock, err = Deserialize(lastBlockBytes)
			chain.IteratorBlockHash = chain.LastBlock.Hash
		}*/
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
func (chain *BlockChain) CreateNewBlock(txs []transaction.Transaction) error {
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
		client.View(func(tx *bolt.Tx) error {
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
	newBlock := NewBlock(lastBlock.Height, lastBlock.Hash, txs)
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
		chain.IteratorBlockHash = newBlock.Hash
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
func (chain *BlockChain) GetLastBlock() Block {
	return chain.LastBlock
}

/*func (chain *BlockChain) GetLastBlock() (Block, error) {
	client := chain.DB
	var err error
	var lastBlock Block
	client.View(func(tx *bolt.Tx) error {
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
func (chain *BlockChain) GetAllBlocks() ([]Block, error) {
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
	blocks := make([]Block, 0)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			err = errors.New("区块数据库操作失败，请重试！")
			return err
		}
		var currentHash []byte //当前遍历到hash
		currentHash = bucket.Get([]byte(LASTHASH))
		//lastHash := bucket.Get([]byte(LASTHASH))
		//2、根据最后一个区块依次往前找
		for {
			currentBlockBytes := bucket.Get(currentHash)
			currentBlock, err := Deserialize(currentBlockBytes)
			if err != nil {
				//lastBlockBytes := bucket.Get(lastHash)
				//return err
				break
			}
			//3、每次找到区块放入到一个[]Block容器中
			blocks = append(blocks, currentBlock)
			//4、找到最开始的创世区块时，结束，不再找了
			if currentBlock.Height == 0 {
				break
			}
			currentHash = currentBlock.PrevHash[:]
		}
		return nil
	})
	return blocks, err
}

/**
 *该方法用于实现迭代器Iterator的HHasNext方法用于判断是否还有数据
 *如果有数据，返回true,否则返回false
 */
func (chain *BlockChain) HasNext() bool {
	/**
	 *区块0 -> 区块2 —> 区块3
	 *最新区块3
	 *步骤：当前区块在哪-> preHash -> client
	 */
	//lastBlock := chain.LastBlock
	db := chain.DB
	var hasNext bool
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败，请重试！")
		}
		//preBlockBytes := bucket.Get(lastBlock.PrevHash[:]) //获取前一个区块数据
		BlockBytes := bucket.Get(chain.IteratorBlockHash[:]) //获取前一个区块数据
		//如果获取不到前一个区块的数据，说明前面没有区块了
		/*	if len(preBlockBytes) == 0{
				hasNext = false
			}else {
				hasNext = true
			}*/
		hasNext = len(BlockBytes) != 0
		return nil
	})
	return hasNext
}

/**
 *该方法用于实现迭代器Iterator的Next方法，用于取出一个数据
 *此处，因为区块链数据集合，因此返回的数据类型是Block
 */
func (chain *BlockChain) Next() Block {
	//1、当前在哪个区块
	//lastBlock := chain.LastBlock
	//2、找当前区块的前一个区块
	//3、找到区块返回
	db := chain.DB
	var iteratorBlock Block
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BLOCKS))
		if bucket == nil {
			return errors.New("区块数据文件操作失败。")
		}
		//前一个区块的数据
		BlockBytes := bucket.Get(chain.IteratorBlockHash[:])
		iteratorBlock, _ = Deserialize(BlockBytes)
		//迭代到当前区块后，更新游标的区块内容
		chain.IteratorBlockHash = iteratorBlock.PrevHash
		return nil
	})
	return iteratorBlock
}

/**
 *该方法用于查询出指定地址的UTXO集合并返回
 */
func (chain *BlockChain) SearchUTXOS(addr string) []transaction.UTXO {
	//花费记录的容器
	spend := make([]transaction.TxInput, 0)
	//收入记录的容器
	inCome := make([]transaction.UTXO, 0)

	//迭代遍历每一个区块
	for chain.HasNext() {
		block := chain.Next() //拿到区块
		//遍历区块中得交易
		for _, tx := range block.Transactions {
			//a、遍历每个交易得交易输入
			for _, input := range tx.Inputs {
				//找到了花费记录
				if string(input.SciptSig) != addr {
					continue
				}
				spend = append(spend, input)
			}
			//b、遍历每个交易得交易输出
			for index, output := range tx.Outputs {
				if string(output.ScriptPub) != addr {
					continue
				}
				utxo := transaction.UTXO{
					TxId:     tx.TxHash,
					Vout:     index,
					TxOutput: output,
				}
				inCome = append(inCome, utxo)
			}
		}

	}
	//遍历收入集合和花费集合，把已花的剔除，找出未花费的记录
	utxos := make([]transaction.UTXO, 0)
	var isComeSpent bool
	for _, come := range inCome {
		isComeSpent = false
		for _, spen := range spend {
			if come.TxId == spen.TxId && come.Vout == spen.Vout {
				//该笔收入已被花费
				isComeSpent = true
				break
			}
		}
		if !isComeSpent { //当前遍历到come未被消费
			utxos = append(utxos, come)
		}
	}
	return utxos
}

/**
 *定义区块链的发送交易的功能
 */
func (chain *BlockChain) SendTransaction(from string, to string, amount float64) error {
	//1、先把form的可花费的utxos找出来
	/*utxos := chain.SearchUTXOS(from)

	//1、先判断当前所有的utxo的可花费的总额是否比转账数额要大
	var totalBlance float64
	var utxoNum int
	for index, utxo := range utxos {
		totalBlance += utxo.Value
		if totalBlance >= amount{
			utxoNum = index
		}*/
	utxos, totalBalance := chain.GetTUXOsWithBalance(from)
	if totalBalance < amount {
		return errors.New("余额不足！！！")
	}
	totalBalance = 0
	var utxoNum int
	for index, utxo := range utxos {
		totalBalance += utxo.Value
		if totalBalance > amount {
			utxoNum = index
			break
		}
	}

	//2、可花费的钱总额比要花费的钱数额大，才创建交易
	newTx, err := transaction.CreateNewTransaction(utxos[0:utxoNum+1], from, to, amount)
	if err != nil {
		return err
	}
	err = chain.CreateNewBlock([]transaction.Transaction{*newTx})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

/**
 *用于实现地址余额查询
 */
func (chain *BlockChain) GetBalane(addr string) float64 {
	_, totalBalance := chain.GetTUXOsWithBalance(addr)
	return totalBalance
}

/**
 *该方法用于实现地址余额统计和地址可以花费的utxo
 */

func (chain *BlockChain) GetTUXOsWithBalance(addr string) ([]transaction.UTXO, float64) {
	utxos := chain.SearchUTXOS(addr)
	var totalBalance float64
	for _, utxo := range utxos {
		totalBalance += utxo.Value
	}
	return utxos, totalBalance
}
