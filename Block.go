package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Timestamp int64
	//Data      []byte
	Transactions []*Transaction
	Hash         []byte
	PrevHash     []byte
	Nonce        int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	enc := gob.NewEncoder(&result)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()

}
func DeserializeBlock(d []byte) *Block {
	var block Block
	enc := gob.NewDecoder(bytes.NewReader(d))
	err := enc.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
func (b *Block) SetHash() {
	Timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{Timestamp, b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}
func newBlock(transactions []*Transaction, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), []byte{}, prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
func NewGenesisBlock() *Block {
	return newBlock("Genesis Block", []byte{})
}
