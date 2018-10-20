package main

import (
	"encoding/gob"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"crypto/elliptic"
)

//定义一个wallets容器，以key-value结构存储每一个wallet
type Wallets struct {
	//hey为地址，value为wallet
	WalletAddress map[string]*Wallet
}

//生成wallets对象的方法
func NewWallets() *Wallets {
	//生成一个wallets对象，这个对象的属性为map，需要遍历存储map的文件
	var ws Wallets
	//map必须初始化才可以使用
	ws.WalletAddress=make(map[string]*Wallet)
	ws.loadFile()
	return &ws
}

//创建新的地址账户
func (ws *Wallets) NewAccount() string{
	//1，生成一个wallet对象
	wallet := NewWallet()
	//2,调用wallet方法，生成钱包地址
	address:=wallet.getAddress()
	//3,调用ws方法，将地址保存到map中，并存入文件
	ws.WalletAddress[address]=wallet
	ws.saveToFile()
	return address
}
//将map存储至文件
const WalletFile  = "wallet.dat"
func (ws *Wallets)saveToFile(){
	//1，创建文件
	//2,编码
	var buffer bytes.Buffer
	gob.Register(elliptic.P224())
	encoder:=gob.NewEncoder(&buffer)

	err:=encoder.Encode(ws)
	if err!=nil{
		log.Panic(err)
	}
	//3,写入文件
	ioutil.WriteFile(WalletFile,buffer.Bytes(),0600)
}
//载入文件
func (ws *Wallets)loadFile(){
	_,err:=os.Stat(WalletFile)
	if os.IsNotExist(err){
    return
	}
	//1，打开文件
	content,err:=ioutil.ReadFile(WalletFile)
	if err!=nil{
		return
	}
	//2，解码
	//注意解码的参数
	//这里编码和下面的解码该curve类型必须注册后才可以使用
	gob.Register(elliptic.P224())
	decoder:=gob.NewDecoder(bytes.NewReader(content))
	var wsTmp Wallets
	decoder.Decode(&wsTmp)
	//3,给ws对象赋值
	ws.WalletAddress=wsTmp.WalletAddress
}
//钱包列表
 func(ws *Wallets)AllAddress()[]string{
 	//遍历map，取出所有的key
 	var addressLidt []string
 	for address:=range ws.WalletAddress{
 		addressLidt=append(addressLidt,address)
	}
	return addressLidt
 }