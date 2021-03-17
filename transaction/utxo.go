package transaction

/**
 *定义结构体utxo表示未花费的交易输出
 */
type UTXO struct {
	TxId     [32]byte //该笔收入来自哪个交易
	Vout     int      //该笔交易来自交易的哪个输出
	TxOutput          //该笔输入的面额和收入者

}
