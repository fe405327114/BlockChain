package main

import (
	"time"
	"crypto/sha256"
	"bytes"
	"encoding/binary"
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
	SelfHash []byte
	Data     []byte
}

//创建一个区块
func NewBlock(data string,prevHash []byte) *Block {
	block := Block{
		Version:   00,
		PrevHash:prevHash,
		Merkle:    []byte{},
		TimeStamp: uint64(time.Now().Unix()),
		//难度和随机数先为空
		Difficulty: 1,
		Nonce:      1,
		SelfHash:   []byte{},
		Data:       []byte(data),
	}
	//生成自身hash
	//block.SetHash()
	pow:=NewProofOfWork(&block)
	hash,nonce:=pow.CrashHash()
	block.SelfHash=hash
	block.Nonce=nonce
	return &block
}

//将属性中的uint64转换成[]byte，方便拼接生成hash
func  Uint64ToByte(num uint64)[]byte {
	buffer := bytes.Buffer{}
	err := binary.Write(&buffer, binary.BigEndian, num)
     if err!=nil{
     	panic(err)
	 }
	 return buffer.Bytes()
}
func (b *Block) SetHash() {
	blockInfo := bytes.Join([][]byte{
	Uint64ToByte(b.Version),
	b.Merkle,
	Uint64ToByte(b.TimeStamp),
	Uint64ToByte(b.Difficulty),
	Uint64ToByte(b.Nonce),
	b.Data,
	}, []byte{})
	hash:= sha256.Sum256(blockInfo)
	b.SelfHash=hash[:]
}
//创建一个创世区块
func GenesisBlock() *Block{
	return NewBlock("这是一个创世区块",[]byte{})
}
