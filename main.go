package main

import "fmt"

func main() {

	blockChain := NewBlockChain()
	blockChain.AddBlock("这是第二个区块")
	blockChain.AddBlock("这是第三个区块")
	//创建一个迭代器
	iterator := NewIterator(blockChain)
	for {
		block := iterator.Next()
		fmt.Printf("===========================\n\n")
		fmt.Printf("区块的头hash：%x\n", block.PrevHash)
		fmt.Printf("区块自身的hash：%x\n", block.SelfHash)
		fmt.Printf("区块的Nonce：%d\n", block.Nonce)
		fmt.Printf("区块的信息为：%s\n", block.Data)
		if len(block.PrevHash)==0{
			fmt.Println("已经打印所有区块")
			break
		}
	}
	//for i,block:=range blockChain.blocks{
	//	fmt.Printf("区块的高度为：%d\n",i)
	//	fmt.Printf("区块的头hash：%x\n",block.PrevHash)
	//	fmt.Printf("区块自身的hash：%x\n",block.SelfHash)
	//	fmt.Printf("区块的Nonce：%d\n",block.Nonce)
	//	fmt.Printf("区块的信息为：%s\n",block.Data)
	//}
}
