package main

import (
	"fmt"
)

func main(){
	//block:=NewBlock("生成一个新的区块")
	//fmt.Printf("区块的头hash：%x\n",block.PrevHash)
	//fmt.Printf("区块自身的hash：%x\n",block.SelfHash)
	//fmt.Printf("区块的信息为：%s\n",block.Data)
	blockChain:=NewBlockChain()
	blockChain.AddBlock("这是第二个区块")
	blockChain.AddBlock("这是第三个区块")
	for i,block:=range blockChain.blocks{
		fmt.Printf("区块的高度为：%d\n",i)
		fmt.Printf("区块的头hash：%x\n",block.PrevHash)
		fmt.Printf("区块自身的hash：%x\n",block.SelfHash)
		fmt.Printf("区块的Nonce：%d\n",block.Nonce)
		fmt.Printf("区块的信息为：%s\n",block.Data)
	}
}