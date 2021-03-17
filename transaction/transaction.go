package transaction

import (
	"XianfengChain04/utils"
	"crypto/sha256"
)

/**
 *定义交易的结构体
 */

type Transaction struct {
	//交易哈希
	TxHash [32]byte

	//交易输入  数据结构：{“txid”,“vount”,"scriptSig"}
	Inputs []TxInput

	//交易输出 数据结构{“value","scriptPub"}
	Outputs []TxOutput
}

/**
 *该函数用于定义一个coinbase交易，并返回该交易结构体
 */
func CreateCoinBase(addr string) (*Transaction, error) {
	//生成一个交易
	output0 := TxOutput{
		Value:     50,
		ScriptPub: []byte(addr),
	}

	/**
	  *构建创世区块
	  *解析矿工的地址
	  *利用定义的交易相关结构体，构建COINBASE transaction
	     特点：只有输出、没有输入
	     output:{
	       value:
	     }
	  *将coinbase 交易传入CreateGenesis方法参数中
	  *CreateGenesis将生成的区块保存到bolt.db文件中
	*/
	coinbase := Transaction{
		//TxHash:  [32]byte{},
		//Inputs:  nil,
		Outputs: []TxOutput{output0},
	}
	coinbaseBytes, err := utils.Encode(coinbase)
	if err != nil {
		//fmt.Println("构建coinbase出现错误，请重试！")
		return nil, err
	}
	coinbase.TxHash = sha256.Sum256(coinbaseBytes)
	return &coinbase, nil
}

/**
 *该函数用于构建一笔普通交易，返回构建好的交易实例
 */

func CreateNewTransaction(from string, to string, amount float64) (*Transaction, error) {
	//1、构建inputs
	inputs := make([]TxInput, 0) //用于存放交易输入的容器

	//改变量用于记录转账发起者一共付了多少钱
	var inputAmount float64

	//从最后一个区块开始，遍历每一个区块
	//在遍历到每一个区块中，遍历区块中所有的交易
	//在遍历到的交易中，遍历所有的交易输入和交易输出

	//a、先充已有的区块账本数据中找到跟from有关的交易输出（所有付给from的钱 ）
	//b、从已有的区块账本数据中找到所有跟from有关的交易输入(所有from花的钱 )
	//c、从所有跟from有关的建议输出记录中提出已经花费的，剩下的就是可花费的

	//2、构建outputs
	outputs := make([]TxOutput, 0) //用于存放交易输出的容器

	//构建转账接收者的交易输出
	output0 := TxOutput{
		Value:     amount,
		ScriptPub: []byte(to),
	}

	//把第一个交易输出放入到专门存交易输出的容器中
	outputs = append(outputs, output0)

	//判断是否需要找零，如果需要，则需要构建一个找零输出
	if inputAmount-amount > 0 {
		output1 := TxOutput{
			Value:     inputAmount - amount,
			ScriptPub: []byte(from),
		}
		outputs = append(outputs,output1)
	}

	//3、构建transaction
	newTransaction := Transaction{
		//TxHash:  ,
		Inputs:  inputs,
		Outputs: outputs,
	}
	//4、计算transaction的哈希，并赋值
	transactionBytes, err := utils.Encode(newTransaction)
	newTransaction.TxHash = sha256.Sum256(transactionBytes)
	if err != nil {
		return nil, err
	}
	//5、将构建的transaction的实例进行返回

	return &newTransaction, err
}
