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

const blockChainDb = "blockChainDb.db"
const blockBucket = "blockBucket"

//创建一个区块链
//创世块信息
const genesisInfo = "genesis block"

func NewBlockChain(address string) *BlockChain {
	//打开数据库
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		panic(err)
	}
	var lastHash []byte
	//查看bucket是否存在，不存在就创建，存在就直接返回对象
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		//不存在就创建
		if bucket == nil {
			bucket, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic(err, "创建bucket失败")
			}
			//创建挖矿交易
			transCoinBase := NewCoinBase(address, genesisInfo)
			//生成创世块
			genesisBlock := NewBlock([]*Transaction{transCoinBase}, []byte{})
			//将区块信息序列化
			blockData := genesisBlock.Serialize()
			bucket.Put(genesisBlock.SelfHash, blockData)
			bucket.Put([]byte("lastHashKey"), genesisBlock.SelfHash)
			//将lastHash更新
			lastHash = genesisBlock.SelfHash
		} else {
			lastHash = bucket.Get([]byte("lastHashKey"))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
	//将创世区块放入区块链中返回
	//genesisBlock:=GenesisBlock()
	//return &BlockChain{
	//	blocks:[]*Block{genesisBlock},
	//}

}

//将区块添加到链中
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	db := bc.db
	lastHash := bc.lastBlockHash
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("该bucket不存在")
		}
		//生成新区块
		block := NewBlock(txs, lastHash)
		//将区块写入数据库
		blockData := block.Serialize()
		bucket.Put(block.SelfHash, blockData)
		bucket.Put([]byte("lastHashKey"), block.SelfHash)
		//更新区块链中的lashHash
		bc.lastBlockHash = block.SelfHash
		return nil
	})
}

//某个地址所有的utxo(未花费的输出)所在的交易的集合
func (bc *BlockChain) AllTransOfUtxo(address string) []Transaction {
	//存储要返回的交易
	var transactions []Transaction
	//已经花费的utxo
	spentUtxo := make(map[string][]int)
	//创建一个迭代器
	iterator := NewIterator(bc)
	//遍历区块
	for {
		block := iterator.Next()
		//遍历区块中的交易
		output:
		for _, tx := range block.Transaction {
			//遍历所有的output，收款信息
			for outputIndex, output := range tx.Outputs {
				//过滤掉该交易中已经花费过的output
				if spentUtxo[string(tx.TxId)]!=nil{
					for _,payIndex:=range spentUtxo[string(tx.TxId)]{
                       if outputIndex==payIndex{
                        continue output
					   }
					}
				}
				//continue是跳过当次循环中剩下的语句，执行下一次循环
				if output.PayBy(address) {
					transactions = append(transactions, *tx)
				}
			}

			//遍历所有的input
			if !tx.IsCoinBase() {
				for _, input := range tx.Inputs {
					//检查当前的input是否是由该地址产生
					if input.UnlockedBy(address) {
						//如果是，则将input中的output标记添加进已花费的utxo中
						spentUtxo[string(input.Txid)] = append(spentUtxo[string(input.Txid)], input.PayIndex)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return transactions
}

//某个地址所有的utxo(未花费的所有输出)的集合
func (bc *BlockChain) AllUtxoOf(address string) []TxOutput {
	//获得所有utxo交易的集合，进行遍历
	txs := bc.AllTransOfUtxo(address)
	//创建需要返回的utxo切片
	utxos := []TxOutput{}
	for _, tx := range txs {
		//遍历获得output
		for _, output := range tx.Outputs {
			//获取该地址的utxos(该地址可以解锁的即为属于该地址的output)
			if output.PayBy(address){
				utxos = append(utxos, *output)
			}
		}
	}
	return utxos
}

//为某个交易查找用到的utxo集合
func (bc *BlockChain) AllUtxoOfTrans(address string, amount float64) (map[string][]int, float64) {
	//创建一个map，存储utxos
	//key为交易id，值为utxo切片
	var validUtxos =make(map[string][]int)
	var transations []Transaction
	var total float64
	//找到某个地址所有的utxo(未花费的输出)所在的交易的集合
	transations=bc.AllTransOfUtxo(address)
    //遍历交易集合，找到足够多的value就结束遍历
    label:
    for _,tx:=range transations{
    	for outputIndex,output:=range tx.Outputs{ //output的索引即为对应utxo
    		if total>=amount{
    	     break label
			}else{
				total+=output.ReValue
				//将用到的utxo索引添加到map中
				validUtxos[string((tx.TxId))]=append(validUtxos[string(tx.TxId)],outputIndex)
			}
		}
	}
	return validUtxos,total
}
