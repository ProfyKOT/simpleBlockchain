package main

import "github.com/boltdb/bolt"

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *BlockchainIterator) Next() *Block {
	var block *Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		block = DeserializeBlock(b.Get(bc.currentHash))
		return nil
	})
	if err != nil {
		panic(err)
	}
	bc.currentHash = block.PrevHash
	return block
}
