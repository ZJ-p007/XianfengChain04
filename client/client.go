package client

import (
	"XianfengChain04/chain"
	"XianfengChain04/utils"
	"flag"
	"fmt"
	"math/big"
	"os"
)

/**
 *该结构体定义了用于实现命令行参数解析的结构体
 */
type CmdClient struct {
	Chain chain.BlockChain
}

func (cmd *CmdClient) Run() {
	args := os.Args
	//1、处理用户没有输入任何命令和参数，打印输出说明书
	if len(args) == 1 {
		cmd.Help()
		return
	}

	//2、解析用户输入的第一个参数，作为功能命令进行解析
	switch os.Args[1] {
	case GENRATEGENSIS:
		//fmt.Println("调用创建创世区块功能")
		cmd.GenerateGenesis()
	case SENDTRANSACTION: //前提是创建区块已存在
		//fmt.Println("调用创建新区块功能")
		cmd.SendTransaction()
		//blockchain := cmd.Chain
	case GETLASTBLOCK:
		//fmt.Println("调用获取最新区块功能")
		cmd.GetLastBlock()
	case GETALLBLOCKS:
		//fmt.Println("调用获取所有区块功能")
		cmd.GetAllBlocks()
	case GETBALANCE:
		cmd.GetBalance()
	case GETNEWADDRESS:
		cmd.GetNewAddress()
	case LISTADDRESS:
		cmd.ListAddress()
	case DUMPPRIVKEY:
		cmd.DumpPrivKey()
	case HELP:
		//fmt.Println("调用帮助说明")
		cmd.Help()
	default:
		//fmt.Println("不支持该命令")
		cmd.Default()
	}

	/*
		createBlock := flag.NewFlagSet("createblock", flag.ExitOnError)
		data := createBlock.String("data", "默认值", "新区块的内容")
		createBlock.Parse(os.Args[2:])
		cmd.Chain.CreateNewBlock([]byte(*data))
	*/
}

/**
 *用户发起交易
 */

func (cmd *CmdClient) GenerateGenesis() {
	generategenesis := flag.NewFlagSet(GENRATEGENSIS, flag.ExitOnError)
	//解析参数
	var addr string
	generategenesis.StringVar(&addr, "address", "", "用户指定的矿工的地址")
	generategenesis.Parse(os.Args[2:])
	//fmt.Println("用户输入的自定义创世区块数据:", addr)
	blockchain := cmd.Chain
	//1、先判断该blockchain是否已存在创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(blockchain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 {
		fmt.Println("创世区块已存在，不能重复生成")
		return
	}

	//2、调用方法实现创世区块的操作
	err := blockchain.CreateCoinBase(addr)
	if err != nil {
		fmt.Println("抱歉创建coinbase交易遇到错误", err.Error())
	}
	fmt.Println("恭喜!生成一笔COINBASE交易，奖励已到账")
}

func (cmd *CmdClient) SendTransaction() {
	createBlock := flag.NewFlagSet(SENDTRANSACTION, flag.ExitOnError)
	from := createBlock.String("from", "", "交易发起人地址")
	to := createBlock.String("to", "", "交易接收者地址")
	amount := createBlock.String("amount", "", "转账的数量")

	if len(os.Args[2:]) > 6 {
		fmt.Println("sendTransaction命令只支持三个参数和参数值，请重试")
		return
	}
	createBlock.Parse(os.Args[2:])

	//from，to，amount三个参数是字符串类型，同时需要满足符合JSON格式
	fromSlice, err := utils.JSONArray2String(*from)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}
	toSlice, err := utils.JSONArray2String(*to)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}
	amountSlice, err := utils.JSONArray2Float(*amount)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}

	//先看看参数个数是否一致
	fromLen := len(fromSlice)
	toLen := len(toSlice)
	amountLen := len(amountSlice)
	if fromLen != toLen || fromLen != amountLen || toLen != amountLen {
		fmt.Println("参数个数不一致，请检查参数后重试")
		return
	}

	//1、先判断是否已生成创世区块，如果没有创世区块，提示用户先生成
	hashBig := new(big.Int)
	hashBig.SetBytes(cmd.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 { //没有创世区块
		fmt.Println("That not a gensis block in blockchain，please use go run main.go generategensis command to create a gensis block first.")
		return
	}

	err = cmd.Chain.SendTransaction(fromSlice, toSlice, amountSlice)
	if err != nil {
		fmt.Println("抱歉，发送交易出现错误：", err.Error())
		return
	}
	fmt.Println("交易发送成功")
}

func (cmd *CmdClient) GetLastBlock() {
	lastBlock := cmd.Chain.GetLastBlock()
	//1、判断是否为空
	hashBig := new(big.Int)
	hashBig.SetBytes(lastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 { //没有最新区块
		fmt.Println("抱歉，当前暂无最新区块.")
		return
	}
	fmt.Println("恭喜，获取到最新区块数据")
	fmt.Printf("最新区块高度:%d\n", lastBlock.Height)
	fmt.Printf("最新区块哈希:%x\n", lastBlock.Hash)

	for index, tx := range lastBlock.Transactions {
		fmt.Printf("区块交易%d,交易:%v\n", index, tx)
	}
}

func (cmd *CmdClient) GetAllBlocks() {
	blocks, err := cmd.Chain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("恭喜，查询到所有区块数据")
	for _, block := range blocks {
		fmt.Printf("区块高度:%d,区块哈希:%x\n", block.Height, block.Hash)
		fmt.Print("区块中的交易信息：\n")
		for index, tx := range block.Transactions {
			fmt.Printf("   第%d笔交易,交易hash:%x\n", index, tx.TxHash)
			for inputIndex, input := range tx.Inputs {
				fmt.Printf("       第%d笔交易输入,%s花了%x的%d的钱\n", inputIndex, input.SciptSig, input.TxId, input.Vout)
			}
			for outputIndex, output := range tx.Outputs {
				fmt.Printf("       第%d笔交易输出,%s实现收入%f\n", outputIndex, output.ScriptPub, output.Value)
			}
		}
		fmt.Println()
	}
}

/**
 *获取地址余额
 */
func (cmd *CmdClient) GetBalance() {
	getbalance := flag.NewFlagSet(GETBALANCE, flag.ExitOnError)
	var addr string
	getbalance.StringVar(&addr, "address", "", "用户地址")
	getbalance.Parse(os.Args[2:])

	blockChain := cmd.Chain

	//1、先判断是否有创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(blockChain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 {
		fmt.Println("抱歉，该链不存在，无法查询")
		return
	}

	balance,err := blockChain.GetBalane(addr)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("地址%s的余额：%f\n", addr, balance)
}

func (cmd *CmdClient) Default() {
	fmt.Println("go run main.go：Unknown subcommand.")
	fmt.Println("Run 'go run main.go help' for usage.")
}

/**
 *定义新方法：用于生成新地址
 */
func (cmd *CmdClient) GetNewAddress() {
	getNewAddress := flag.NewFlagSet(GETNEWADDRESS,flag.ExitOnError)
	getNewAddress.Parse(os.Args[2:])

	if len(os.Args[2:]) > 0{
		fmt.Println("抱歉生成新地址功能无法解析参数")
		return
	}
	address,err := cmd.Chain.GetNewAddress()
	if err !=nil{
		fmt.Println("生成地址遇到错误",err)
		return
	}
	fmt.Println("生成新地址:",address)
}

func (cmd *CmdClient) ListAddress() {
	listAddress := flag.NewFlagSet(LISTADDRESS,flag.ExitOnError)
	listAddress.Parse(os.Args[2:])
	if len(os.Args[2:]) > 0{
		fmt.Println("无法解析参数，请检查后重试！！！")
		return
	}
	addList ,err := cmd.Chain.GetAddressList()
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println("获取地址列表成功，地址信息如下：")
	for index, add := range addList {
		fmt.Printf("[%d]:%s\n", index+1, add)
	}

}

/**
 *该方法用于导出指定地址的私钥
 */
func (cmd *CmdClient) DumpPrivKey() {
	dumpPrivkey := flag.NewFlagSet(DUMPPRIVKEY,flag.ExitOnError)
	address := dumpPrivkey.String("address","","要导出的私钥的地址")
	dumpPrivkey.Parse(os.Args[2:])
	if len(os.Args[2:]) > 2{
		fmt.Println("无法解析参数，请检查后重试！！！")
		return
	}
	pri,err:= cmd.Chain.DumpPrivkey(*address)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("私钥是:%x\n",pri.D.Bytes())
}



/**
 *该方法用于打印输出项目的使用和说明信息，相当于项目的帮助文档和说明书
 */
func (cmd *CmdClient) Help() {
	fmt.Println("------Welcome to XianfengChain04 Project------")
	fmt.Println("XianfengChain04 is a custom blockchain project, the projects plan to a very simple public chain")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("go run main.go command [arguments]")
	fmt.Println("AVAILABLE COMMANDS")
	fmt.Println("    generategenesis    use the command can create a genesis block and save to the boltdb file. use genesis argument to set")
	fmt.Println("    sendtransaction    this command used to create a new transaction, that can specified an argument ")
	fmt.Println("    getlastblock    get the lasted block data")
	fmt.Println("    getbalance    this is a command that can the balance of specified address")
	fmt.Println("    getallblocks    return a blocks data to user.")
	fmt.Println("        getnewaddress       this command used to create a new address by bitcoin algorithm")
	fmt.Println("    help    use the command can print usage infomation")
	fmt.Println()
	fmt.Println("Use bee help [command] for more information about a command")
}
