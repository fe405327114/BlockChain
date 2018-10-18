package main

import (
	"os"
	"fmt"
	"strconv"
)

//利用命令行实现区块的添加和区块链的打印
type Cli struct {
	bc *BlockChain
}

const Usage = `
addBlock --data DATA (add a block to the chain)
getBalance  --address ADDRESS(get the balance)
printChain  (print the blockchain)
send from to amount miner data (send some coins to another)
`

func (cli *Cli) Run() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}
	cmd := args[1]
	switch cmd {
	case "printChain":
		//打印区块链
		fmt.Println("打印区块链")
		cli.PrintChain()
	case "getBalance":
		if len(args) > 3 && args[2] == "--address" {
			if args[3]==""{
				fmt.Println("address should not be empty!")
			}
			//获取余额
			fmt.Println("获取余额")
			cli.GetBalance(args[3])
		}
	case "send":
		if len(args)!=7{
			fmt.Println(Usage)
			os.Exit(1)
		}
		from:=args[2]
		to:=args[3]
		amount,_:=strconv.ParseFloat(os.Args[4], 64)
		miner:=args[5]
		data:=args[6]
		cli.Send(from,to,amount,miner,data)
	default:
		fmt.Println(Usage)
	}

}
