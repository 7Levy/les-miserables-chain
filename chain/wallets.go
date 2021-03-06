package chain

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

//钱包数据文件
const walletDataFile = "wallets.data"

//钱包集
type Wallets struct {
	WalletMap map[string]*Wallet
}

//创建钱包集合
func NewWallets() (*Wallets, error) {
	//获取文件属性，如果文件不存在，那么创建一个钱包集合
	if _, err := os.Stat(walletDataFile); os.IsNotExist(err) {
		wallets := &Wallets{}
		wallets.WalletMap = make(map[string]*Wallet)
		return wallets, err
	}
	fileContent, err := os.ReadFile(walletDataFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	return &wallets, nil
}

//创建一个钱包
func (w *Wallets) CreateNewWallet() {
	wallet := NewWallet()
	walletAddress := wallet.GetAddress()
	fmt.Printf("新钱包地址为: %s\n", walletAddress)

	w.WalletMap[string(walletAddress)] = wallet
	fmt.Println("钱包数据写入磁盘中...")
	w.SaveWallet()
	time.Sleep(3 * time.Second)
	fmt.Println("钱包数据写入完成")
}

//钱包持久化
func (w *Wallets) SaveWallet() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panic(err)
	}
	err = os.WriteFile(walletDataFile, content.Bytes(), 0744)
	if err != nil {
		log.Panic(err)
	}
}
