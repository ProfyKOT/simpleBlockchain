package main

import "fmt"

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) SetID() {

	lastBlockchain := NewBlockchain().Iterator()
	lastBlock := lastBlockchain.Next()
	if lastBlock.Transactions[0].ID {
		tx.ID = lastBlock.Transactions[0].ID
	} else {
		tx.ID = []byte{1}
	}
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{12000, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{{txout}}}
	tx.SetID()
	return &tx
}
