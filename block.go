package main

import (
	"time"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"crypto/sha256"
)

type Block struct {
	Version uint64
	//头hash
	PrevHash   []byte
	Merkle     []byte
	TimeStamp  uint64
	Difficulty uint64
	Nonce      uint64
	//自身hash
	SelfHash    []byte
	Transaction []*Transaction
}

//创建一个区块
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	block := Block{
		Version:   00,
		PrevHash:  prevHash,
		Merkle:    []byte{},
		TimeStamp: uint64(time.Now().Unix()),
		//难度和随机数先为空
		Difficulty:  1,
		Nonce:       1,
		SelfHash:    []byte{},
		Transaction: txs,
	}
	//生成自身hash
	//block.SetHash()
	pow := NewProofOfWork(&block)
	hash, nonce := pow.CrashHash()
	block.SelfHash = hash
	block.Nonce = nonce
	return &block
}

//将属性中的uint64转换成[]byte，方便拼接生成hash
func Uint64ToByte(num uint64) []byte {
	buffer := bytes.Buffer{}
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

//将区块信息序列化
func (bc *Block) Serialize() []byte {
	buffer := bytes.Buffer{}
	//创建编码器
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&bc)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

//将区块信息反序列化
func Deserialize(blockData []byte) *Block {
	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(blockData))
	var block Block
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}

//生成Merkel Tree Root哈希值，使用所有交易的哈希值生成一个平衡二叉树，
// 此处，为了简化代码，我们目前直接将区块中交易的id进行拼接后进行哈希操作即可
func (block *Block) MerkleRootHash() []byte {
	//定义一个储存所有交易id的变量
	var TxIds [][]byte
   //遍历区块，获得所有id
   for _,tx:=range block.Transaction{
   	TxIds=append(TxIds,tx.TxId)
   }
   //将元素拼接
   TxIdStr:=bytes.Join(TxIds,[]byte{})
   //生成hash
	TxIdsHash:=sha256.Sum256(TxIdStr)
	return TxIdsHash[:]
}
