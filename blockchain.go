package main

//定义一个区块链结构
type BlockChain struct {
	blocks []*Block
}
//创建一个区块链
func NewBlockChain()*BlockChain{
	//将创世区块放入区块链中返回
	genesisBlock:=GenesisBlock()
	return &BlockChain{
		blocks:[]*Block{genesisBlock},
	}
}
//将区块添加到链中
func (bc *BlockChain)AddBlock(data string)  {
	//获取区块的头哈希
	lastBlock:=bc.blocks[len(bc.blocks)-1]
	prevHash:=lastBlock.SelfHash
	//创建一个新的区块
	block:=NewBlock(data,prevHash)
	//将区块添加进去
	bc.blocks=append(bc.blocks,block)
}