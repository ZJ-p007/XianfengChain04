package transaction

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
