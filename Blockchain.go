package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type Blockchain struct {
	tips []byte
	db   *bolt.DB
}

func (bc *Blockchain) addBlock(data string) {
	var lastHash []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	block := newBlock(data, lastHash)

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

func NewBlockchain() *Blockchain {
	var tips []byte
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b == nil {
			genesis := NewGenesisBlock()
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
