package main

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

type Blockchain struct {
	tips []byte
	db   *bolt.DB
}

func (bc *Blockchain) mineBlock(transactions []*Transaction) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	block := newBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		err := b.Put(block.Hash, block.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			return err
		}
		bc.tips = block.Hash
		return nil
	})
}

func NewBlockchain(address string) *Blockchain {
	var tips []byte
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b == nil {
			cbtx := NewCoinbaseTX(address, 24000000000000000, "")
			genesis := NewGenesisBlock(cbtx)
			b, err = tx.CreateBucket([]byte("blocks"))
			err := b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}
			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				return err
			}
			tips = genesis.Hash
		} else {
			tips = b.Get([]byte("l"))
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	result := Blockchain{tips, db}

	return &result
}
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tips, bc.db}
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	spentTXOs := make(map[string][]int)
	var unspentTXs []Transaction
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, input := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if input.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.isCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTXs
}
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
