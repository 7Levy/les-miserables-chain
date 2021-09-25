package chain

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"les-miserables-chain/persistence"
	"log"
	"math/big"
)

//链结构体
type Chain struct {
	LastHash []byte   //链的最新高度区块hash
	DB       *bolt.DB //数据库对象
}

//创世区块链
func NewBlockChain() *Chain {
	var lastHash []byte

	db, err := bolt.Open(persistence.DbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(persistence.BlockBucket))
		//判断bucket是否存在
		if b == nil {
			fmt.Println("Creating the genesis block.....")
			//创世区块集成交易
			coinbaseTx := NewCoinBaseTX("levy", "In a soldier's stance, I aimed my hand at the mongrel dogs who teach")
			genesisBlock := NewGenesisBlock(coinbaseTx)
			//bucket不存在，创建一个桶
			b, err := tx.CreateBucket([]byte(persistence.BlockBucket))
			if err != nil {
				log.Panic(err)
			}
			//创世区块存储到bucket中
			err = b.Put(genesisBlock.BlockCurrentHash, Serialize(genesisBlock))
			if err != nil {
				log.Panic(err)
			}
			//存储最新的出块hash
			err = b.Put([]byte("last"), genesisBlock.BlockCurrentHash)
			if err != nil {
				log.Panic(err)
			}
			lastHash = genesisBlock.BlockCurrentHash
		} else {
			lastHash = b.Get([]byte("last"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return &Chain{
		LastHash: lastHash,
		DB:       db,
	}
}

//区块派生
func (chain *Chain) AddBlock(transactions []*Transaction) {
	//1.创建区块
	newBlock := NewBlock(transactions, chain.LastHash)
	//2.区块bucket更新
	err := chain.DB.Update(func(tx *bolt.Tx) error {
		//获取当前表
		b := tx.Bucket([]byte(persistence.BlockBucket))
		//存储区块数据
		err := b.Put(newBlock.BlockCurrentHash, Serialize(newBlock))
		if err != nil {
			log.Panic(err)
		}
		//存储最新出块的hash
		err = b.Put([]byte("last"), newBlock.BlockCurrentHash)
		if err != nil {
			log.Panic(err)
		}
		//更新最新出块的hash
		chain.LastHash = newBlock.BlockCurrentHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

}

//查询地址下的未花费输出的交易集合
func (chain *Chain) FindUnspentTransactions(address string) []Transaction {
	//存储未花费的交易
	var unspentTxs []Transaction
	spentTxs := make(map[string][]int)
	blockchainIterator := chain.Iterator()
	var hashInt big.Int

	for {
		err := blockchainIterator.DB.View(func(tx *bolt.Tx) error {
			//获取当前区块
			b := tx.Bucket([]byte(persistence.BlockBucket))
			blockBytes := b.Get(blockchainIterator.CurrentHash)
			block := DeserializeBlock(blockBytes)

			for _, transaction := range block.Transactions {
				fmt.Printf("TransactionHash:%x\n", transaction.Index)
				//将交易ID转换为16进制
				index := hex.EncodeToString(transaction.Index)
				//Outputs的label
			Outputs:
				for outIdx, out := range transaction.Outputs {
					if spentTxs[index] != nil {
						for _, spentOut := range spentTxs[index] {
							if spentOut == outIdx {
								continue Outputs
							}
						}
					}
					if out.UnlockOutput(address) {
						unspentTxs = append(unspentTxs, *transaction)
					}
				}
				if transaction.IsCoinbase() == false {
					for _, in := range transaction.Inputs {
						if in.UnlockInput(address) {
							inTxID := hex.EncodeToString(in.TxID)
							spentTxs[inTxID] = append(spentTxs[inTxID], in.OutputIndex)
						}
					}
				}

			}
			fmt.Println()
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		blockchainIterator = blockchainIterator.Next()
		hashInt.SetBytes(blockchainIterator.CurrentHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return unspentTxs
}
