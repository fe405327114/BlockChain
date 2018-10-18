package main

import (
	"bytes"
	"encoding/gob"
	"crypto/sha256"
	"fmt"
	"os"
)

//定义付款信息结构
type TxInput struct {
	//output所在的交易id
	Txid []byte
	//付款的output索引，用于剔除UTXO中的已花费输出
	PayIndex int
	//私钥签名
	ScriptSig string
}

//定义收款结构信息
type TxOutput struct {
	//接收的金额
	ReValue float64
	//公钥信息，用于解锁付款人的签名
	ScriptPubKey string
}

//定义交易结构信息
type Transaction struct {
	//交易id
	TxId []byte
	//付款端
	Inputs []*TxInput
	//收款端
	Outputs []*TxOutput
}

//生成交易id
func (tx *Transaction) SetTxId() {
	//字节流缓存器
	var buffer bytes.Buffer
	//编码器
	encoder := gob.NewEncoder(&buffer)
	//编码
	err := encoder.Encode(tx)
	if err != nil {
		panic(err)
	}
	//对交易数据进行hash
	TxHash := sha256.Sum256(buffer.Bytes())
	tx.TxId = TxHash[:]
}

//创建挖矿交易
const reward = 12.5
func NewCoinBase(address, data string) *Transaction {
	//挖矿交易没有input，只有output
	//比特币系统，对于这个input的id填0，对索引填0xffff，data由矿工填写，一般填所在矿池的名字
	input := TxInput{[]byte{}, -1, data}
	output := TxOutput{reward, address}
	TransCoinBase := Transaction{[]byte{}, []*TxInput{&input}, []*TxOutput{&output}}
	//生成交易id
	TransCoinBase.SetTxId()
	return &TransCoinBase
}
//检查是否为该地址的input，
func (input *TxInput) UnlockedBy(unlockData string) bool {
	//此方法用于验证传入的数据是否可以解开
	return input.ScriptSig == unlockData
}
//检查是否为该地址的output,
func (output *TxOutput) PayBy(unlockData string) bool {
	return output.ScriptPubKey == unlockData
}
//判断是否为挖矿交易
func (tx *Transaction) IsCoinBase() bool {
	//挖矿交易没有input
	if len(tx.Inputs) == 1 {
		//2. 交易id为空
		//3. 交易的index 为 -1
		if len(tx.TxId)==0 && tx.Inputs[0].PayIndex == -1  {
			return true
		}
	}
	return false
}
//创建交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction{
	//第一部分：找到所有需要的UTXO（有可能余额不足，余额正好，有剩余）
	//第二部分：根据找到UTXO生成Input
	//第三部分：生成Output，包括给收款人以及找零（如果有剩余）
	//获取交易用到的utxo以及其包含余额
	var validUtxos map[string][]int//key为交易id，value为output索引值
	 var total float64
	validUtxos,total=bc.AllUtxoOfTrans(from,amount)
	fmt.Printf("找到用于交易的utxo：%f\n",total)
	if total<amount{
		fmt.Println("余额不足")
		os.Exit(1)
	}
	//将validUtxos转化为input和output
	var inputs []*TxInput
	var outputs []*TxOutput
	for txId,utxos:=range validUtxos{
     //遍历utxos，将得到的值添加进已经花费的input.payindex中
     for _,noPayIndex:=range utxos{
     	input:=TxInput{[]byte(txId),noPayIndex,from}
     	inputs=append(inputs,&input)
	 }
	}
   //创建output
   output:=TxOutput{amount,to}
  outputs=append(outputs,&output)
   //找零
   if total>amount{
   	//output=TxOutput{total-amount,from}
   	//此处有个坑，这里不可以重新给output赋值
   	outputs=append(outputs,&TxOutput{total-amount,from})
   }
   fmt.Println("amount:",amount)
   transation:=Transaction{[]byte{},inputs,outputs}
   transation.SetTxId()
   return &transation
}
