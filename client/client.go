package client

import (
	"XianfengChain04/chain"
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
	case CREATEBLOCK: //前提是创建区块已存在
		//fmt.Println("调用创建新区块功能")
		cmd.CreateBlock()
		//blockchain := cmd.Chain
	case GETLASTBLOCK:
		//fmt.Println("调用获取最新区块功能")
		cmd.GetLastBlock()
	case GETALLBLOCKS:
		//fmt.Println("调用获取所有区块功能")
		cmd.GetAllBlocks()
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

func (cmd *CmdClient) GenerateGenesis() {
	generategenesis := flag.NewFlagSet(GENRATEGENSIS, flag.ExitOnError)
	var genesis string
	generategenesis.StringVar(&genesis, "genesis", "", "创世区块的数据")
	generategenesis.Parse(os.Args[2:])
	fmt.Println("用户输入的自定义创世区块数据:", genesis)
	blockchain := cmd.Chain
	//1、先判断该blockchain是否已存在创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(blockchain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 {
		fmt.Println("创世区块已存在，不能重复生成")
		return
	}
	//2、调用方法实现创世区块的操作
	blockchain.CreatGenesis([]byte(genesis))
	fmt.Println("创世区块已生成，并保存到文件中。")
}

func (cmd *CmdClient) CreateBlock() {
	createblock := flag.NewFlagSet(CREATEBLOCK, flag.ExitOnError)
	//	var create string
	data := createblock.String("data", "", "创建的新区块的自定义内容")

	if len(os.Args[2:]) > 2 {
		fmt.Println("CREATEBLOCK只支持data一个参数，请重试")
		return
	}
	//args := os.Args[2:]
	createblock.Parse(os.Args[2:])
	//1、先判断是否已生成创世区块，如果没创世区块，提示用户先生成
	hashBig := new(big.Int)
	hashBig.SetBytes(cmd.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 { //没有创世区块
		fmt.Println("That not a genesis block in blockchain, please use go run main.go generategenesis comand create a genesis block first")
		return
	}
	//2、生成一个新区块，存到文件中
	err := cmd.Chain.CreateNewBlock([]byte(*data))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("新区块创建成功,并成功保存到文件中")
}

func (cmd *CmdClient) GetLastBlock() {
	lastBlock := cmd.Chain.LastBlock
	//1、判断是否唯恐
	hashBig := new(big.Int)
	hashBig.SetBytes(lastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 { //当前没有最新区块
		fmt.Println("当前暂五最新区块")
		return
	}
	fmt.Println("获取到最新区块")
	fmt.Printf("最新区块的高度:%d\n", lastBlock.Height)
	fmt.Printf("最新区块的哈希:%x\n", lastBlock.Hash)
	for _,tx := range lastBlock.Transactions{
		fmt.Printf("区块交易:%d,交易：%v\n", lastBlock.Transactions,tx)
	}

}

func (cmd *CmdClient) GetAllBlocks() {
	blocks, err := cmd.Chain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("获取到所有区块数据")
	for _, block := range blocks {
		fmt.Printf("区块高度:%d,区块哈希:%x,区块数据:%s\n", block.Height, block.Hash, block.Transactions)
	}
}

func (cmd *CmdClient) Default() {
	fmt.Println("go run main.go: Unknown subcommand.")
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
	fmt.Println("    createblock    this command used to create a new block, that can specified an argument ")
	fmt.Println("    getlastblock    get the lasted block data")
	fmt.Println("    getallblocks    return a blocks data to user.")
	fmt.Println("    help    use the command can print usage infomation")
	fmt.Println()
	fmt.Println("Use bee help [command] for more information about a command")
}
