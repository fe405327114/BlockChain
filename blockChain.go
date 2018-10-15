package main

import (
	"BlockChain/bolt"
	"log"
)

//定义一个区块链结构
type BlockChain struct {
	//blocks []*Block
	db *bolt.DB
	//保存最后一个hash值
	lastBlockHash []byte
}

const blockChainDb  = "blockChainDb.db"
const blockBucket  ="blockBucket"
//创建一个区块链
func NewBlockChain()*BlockChain{
	//打开数据库
	db,err:=bolt.Open(blockChainDb,0600,nil)
	if err!=nil{
		panic(err)
	}
	var lastHash []byte
	//查看bucket是否存在，不存在就创建，存在就直接返回对象
	db.Update(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(blockBucket))
		//不存在就创建
		if bucket==nil{
			bucket,err:=tx.CreateBucket([]byte(blockBucket))
			if err!=nil{
				log.Panic(err,"创建bucket失败")
			}
			//将创世块添加进去
			genesisBlock:=GenesisBlock()
			//将区块信息序列化
			blockData:=genesisBlock.Serialize()
			bucket.Put(genesisBlock.SelfHash,blockData)
			bucket.Put([]byte("lastHashKey"),genesisBlock.SelfHash)
			//将lastHash更新
			lastHash=genesisBlock.SelfHash
		}else{
			lastHash=bucket.Get([]byte("lastHashKey"))
		}
		return nil
	})
	return &BlockChain{db,lastHash}
	//将创世区块放入区块链中返回
	//genesisBlock:=GenesisBlock()
	//return &BlockChain{
	//	blocks:[]*Block{genesisBlock},
	//}

}
//将区块添加到链中
func (bc *BlockChain)AddBlock(data string)  {
	db:=bc.db
	lastHash:=bc.lastBlockHash
	db.Update(func(tx *bolt.Tx) error {
		bucket:=tx.Bucket([]byte(blockBucket))
		if bucket==nil{
			log.Panic("该bucket不存在")
		}
		//生成新区块
		block:=NewBlock(data,lastHash)
		//将区块写入数据库
		blockData:=block.Serialize()
		bucket.Put(block.SelfHash,blockData)
		bucket.Put([]byte("lastHashKey"),block.SelfHash)
		//更新区块链中的lashHash
		bc.lastBlockHash=block.SelfHash
		return nil
	})
}