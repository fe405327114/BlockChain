package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"crypto/sha256"
	"BlockChain/base58"
	"golang.org/x/crypto/ripemd160"
)

//1，定义一个结构，包含公钥私钥两个变量
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	//由两个坐标点拼接而成的临时公钥，便于传输，校验时进行拆分，还原成原始的公钥
	PublicKey []byte
}
//goland.org/x/crypto/ripemd160
//2，提供一个方法：生成公钥私钥
//3，提供一个方法：由公钥生成地址
func NewWallet() *Wallet {
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//生成公钥
	publicKeyRaw := privateKey.PublicKey
	//拼接
	publickey := append(publicKeyRaw.X.Bytes(), publicKeyRaw.Y.Bytes()...)
	return &Wallet{privateKey, publickey}

}

//由公钥生成地址
func (w *Wallet) getAddress() string {
	// 1,对公钥做哈希处理
	//2，拼接版本信息
	//3,两次哈希运算，取前四个字符作为校验码
	//4,拼接2与3
	//5,base58
	var version=byte(00)
	originPubHash := sha256.Sum256(w.PublicKey)
	riper:=ripemd160.New()
	riper.Write(originPubHash[:])
	ripeHash:=riper.Sum(nil)

	pubVersion := append([]byte{version}, ripeHash...)
	//生成校验码
	checkNum := checkNum(pubVersion)
	//拼接
	pubHash := append(pubVersion, checkNum...)
	//base58编码
	// go get et github.com/btcsuite/btcutil/base58
	address := base58.Encode(pubHash)
	return address
}

//生成校验码
func checkNum(pubVersion []byte) []byte {
	//3,两次哈希运算，取前四个字符作为校验码
	firstHash := sha256.Sum256(pubVersion)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:4]
}
