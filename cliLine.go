package main

import "fmt"

//func(cli *Cli) AddBlock(data string){
//	cli.bc.AddBlock(data)
//}
//打印区块链
func (cli *Cli) PrintChain() {
	//创建一个迭代器
	iterator := NewIterator(cli.bc)
	for {
		block := iterator.Next()
		fmt.Printf("===========================\n\n")
		fmt.Printf("区块的头hash：%x\n", block.PrevHash)
		fmt.Printf("区块自身的hash：%x\n", block.SelfHash)
		fmt.Printf("区块的Nonce：%d\n", block.Nonce)
		//fmt.Printf("区块的信息为：%s\n", block.Data)
		fmt.Printf("区块的数据为:%s", block.Transaction[0].Inputs[0].ScriptSig)
		if len(block.PrevHash) == 0 {
			fmt.Println("已经打印所有区块")
			break
		}
	}
}
//查询余额
func (cli *Cli) GetBalance(address string) {
	//获取该地址的utxos
	utxos := cli.bc.AllUtxoOf(address)
	//遍历，汇总
	var total float64
	for _, utxo := range utxos {
		total += utxo.ReValue
	}
	fmt.Printf("The balance of %s is : %f\n", address, total)
}
//发起交易
func (cli *Cli) Send(from, to string, amount float64,miner,data string) {
	fmt.Println("开始转账")
	//cli.GetBalance(from)
	//添加挖矿交易
	coinBase:=NewCoinBase(miner,data)
	//添加普通交易
	transaction := NewTransaction(from, to, amount, cli.bc)
    //添加到区块链
    cli.bc.AddBlock([]*Transaction{coinBase,transaction})
	//fmt.Println("from:",from)
	//fmt.Println("to:",to)
	//fmt.Println("amount:",amount)
	//fmt.Println("miner:",miner)
	//fmt.Println("data:",data)
	//cli.GetBalance(from)
	//cli.GetBalance(to)
	//fmt.Println()
    fmt.Println("send successfully!")
}
//创建钱包
func (cli *Cli)CreateWallet(){
	//创建ws对象
	ws:=NewWallets()
	//调用创建account的方法
	//该方法会生成一个wallet对象，同时创建新的地址，
	address:=ws.NewAccount()
	//打印地址
	fmt.Printf("New address:%s\n",address)
}
//钱包列表
func (cli *Cli)List(){
	//获取ws对象
	ws:=NewWallets()
	addressLidt:=ws.AllAddress()
	for _, address := range addressLidt {
		fmt.Printf("地址：%s\n", address)
	}
}