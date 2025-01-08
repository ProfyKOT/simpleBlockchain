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
	Timestamp    int64
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
	headers := bytes.Join([][]byte{Timestamp, b.HashTransactions(), b.PrevHash}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}
func newBlock(transactions []*Transaction, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, []byte{}, prevHash, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}
func NewGenesisBlock(coinbase *Transaction) *Block {
	return newBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
