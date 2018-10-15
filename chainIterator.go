package main

import "BlockChain/bolt"

//定义迭代器，用于循环遍历区块链，返回各个块数据
type Iterator struct {
	db *bolt.DB
	//当前指向的块hash
	currentHash []byte
}
//生成一个迭代器
func NewIterator (bc *BlockChain)*Iterator{
	return &Iterator{bc.db,bc.lastBlockHash}
}
//迭代器运作
func (it *Iterator)Next()*Block{
	var block *Block
	//输出第最后一个区块信息
    it.db.View(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(blockBucket))
		if bucket==nil{
			panic("bucket不存在")
		}
		blockData:=bucket.Get(it.currentHash)
		//对数据进行解码
		block=Deserialize(blockData)
		//将hash指向前一个区块
		it.currentHash=block.PrevHash
		return nil
	})
	return block
}