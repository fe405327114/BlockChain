package main

func main() {

	blockChain := NewBlockChain("myaddress")
	cli:=Cli{blockChain}
	cli.Run()
	/*
	blockChain.AddBlock("这是第二个区块")
	blockChain.AddBlock("这是第三个区块")

	*/
}
