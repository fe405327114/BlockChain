package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

//定义一个工作量证明结构
type ProofOfWork struct {
	//区块信息
	block *Block
	//目标值
	target *big.Int
}
//定义函数，返回对象
func NewProofOfWork(block *Block)*ProofOfWork{
	var pow=ProofOfWork{
		block:block,
	}
	targetString:="0001000000000000000000000000000000000000000000000000000000000000"
	var tarTem big.Int
	//将string转为bigint
	tarTem.SetString(targetString,16)
	pow.target=&tarTem
	return &pow
}
//计算hash,返回hash和nonce
func (pow *ProofOfWork)CrashHash()([]byte,uint64){
	//1,将区块信息拼接算出hash值，转为big.int
	var Nonce uint64=0
	var hash=[]byte{}
	for  {
		hash=pow.JoinHash(Nonce)
		var bigIntHash=big.Int{}
		bigIntHash.SetBytes(hash)
		//2，与目标值进行比较
		if bigIntHash.Cmp(pow.target)==-1{
			fmt.Printf("found hash : %x, %d\n", hash, Nonce)
			break
		}else {
			Nonce++
		}
	}
	return hash,Nonce
}
func(pow *ProofOfWork) JoinHash(nonce uint64)[]byte{
 var block =pow.block
 block.Merkle=block.MerkleRootHash()
	blockInfo := bytes.Join([][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.Merkle,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(nonce),
		//block.Data,
	}, []byte{})
	hash:=sha256.Sum256(blockInfo)
	return hash[:]
}

