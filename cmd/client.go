package cmd

import (
	"flag"
	"fmt"
	"les-miserables-chain/chain"
	"les-miserables-chain/database"
	"les-miserables-chain/utils"
	"log"
	"os"
)

type CLI struct {
	Chain *chain.Chain
}

//校验参数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

//客户端运行
func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("请先完成本机节点ID的相关配置!\n")
		os.Exit(1)
	}
	fmt.Printf("当前运行节点：%s\n", nodeID)
	database.GenerateDatabase(nodeID) //生成节点数据库

	CmdPrintChain := flag.NewFlagSet("printchain", flag.ExitOnError)     //打印区块链
	CmdDelete := flag.NewFlagSet("delete", flag.ExitOnError)             //删除区块链
	CmdInit := flag.NewFlagSet("init", flag.ExitOnError)                 //初始化区块链
	CmdGetBalance := flag.NewFlagSet("balance", flag.ExitOnError)        //获取账户余额
	CmdSendToken := flag.NewFlagSet("send", flag.ExitOnError)            //转账
	CmdCreateWallet := flag.NewFlagSet("createwallet", flag.ExitOnError) //创建钱包
	CmdAddressLists := flag.NewFlagSet("addresslists", flag.ExitOnError) //获取所有钱包地址
	CmdStartNode := flag.NewFlagSet("startnode", flag.ExitOnError)       //启动节点

	cbAddr := CmdInit.String("address", "", "创世区块奖励人")
	balanceAddr := CmdGetBalance.String("addr", "", "获取指定地址的余额")
	sendFrom := CmdSendToken.String("from", "", "转账源地址")
	sendTo := CmdSendToken.String("to", "", "转账目的地址")
	sendAmount := CmdSendToken.String("amount", "", "转账金额")
	sendMine := CmdSendToken.Bool("mine", false, "是否启用本地节点验证")
	MinerAddress := CmdStartNode.String("miner", "", "挖矿奖励的地址")

	switch os.Args[1] {
	case "printchain":
		err := CmdPrintChain.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "init":
		err := CmdInit.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "delete":
		err := CmdDelete.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getbalance":
		err := CmdGetBalance.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := CmdSendToken.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := CmdCreateWallet.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "addresslists":
		err := CmdAddressLists.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := CmdStartNode.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if CmdInit.Parsed() {
		if *cbAddr == "" {
			cli.printUsage()
			os.Exit(1)
		}
		if !chain.CheckAddress([]byte(*cbAddr)) {
			fmt.Println("地址格式错误!")
			cli.printUsage()
			os.Exit(1)
		}
		cli.initialize(*cbAddr)
	}
	if CmdPrintChain.Parsed() {
		cli.printChain()
	}
	if CmdDelete.Parsed() {
		cli.deleteChain()
	}
	if CmdGetBalance.Parsed() {
		//fmt.Println(*balanceAddr)
		if *balanceAddr == "" {
			cli.printUsage()
			os.Exit(1)
		}
		cli.getBalance(*balanceAddr)
	}
	if CmdSendToken.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount == "" {
			cli.printUsage()
			os.Exit(1)
		}
		from := utils.JsonToArray(*sendFrom)
		to := utils.JsonToArray(*sendTo)
		for index, fromAddress := range from {
			if chain.CheckAddress([]byte(fromAddress)) == false || chain.CheckAddress([]byte(to[index])) == false {
				fmt.Println("地址格式错误!")
				cli.printUsage()
				os.Exit(1)
			}
		}
		amount := utils.JsonToArray(*sendAmount)
		cli.sendToken(from, to, amount, *sendMine)
	}
	if CmdCreateWallet.Parsed() {
		cli.createWallet()
	}
	if CmdAddressLists.Parsed() {
		cli.addresslists()
	}
	if CmdStartNode.Parsed() {
		if *MinerAddress == "" {
			cli.printUsage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *MinerAddress)
	}
}
