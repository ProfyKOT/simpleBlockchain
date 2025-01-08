package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) getBalance(address string) {
	outputLists := cli.bc.FindUTXO(address)
	var balance int

	for _, output := range outputLists {
		balance += output.Value

	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from string, to string, amount int) {
	newTx := NewUTXOTransaction(from, to, amount, cli.bc)

	cli.bc.mineBlock([]*Transaction{newTx})

	fmt.Printf("Send %d\n", amount)
}

func (cli *CLI) printData() {
	blockItem := cli.bc.Iterator()

	for {
		block := blockItem.Next()
		fmt.Println("Block: ------------------")
		fmt.Printf("Block Time: %s\n", time.Unix(
			block.Timestamp,
			0).UTC().Format("2006-01-02"))
		fmt.Printf("%s %x\n", "Block hash:", block.Hash)
		fmt.Printf("%s %x\n", "Block Previous Hash:", block.PrevHash)
		fmt.Printf("%s\n", "Block Data:")
		for _, transaction := range block.Transactions {
			fmt.Printf("\t%s%x\n", "Transaction ID: ", transaction.ID)
			fmt.Printf("\t%s\n", "TXInputs:")
			for _, vin := range transaction.Vin {
				fmt.Printf("\t\t%s%x\n", "TXInput ID:", vin.Txid)
				fmt.Printf("\t\t%s%x\n", "Vout ID:", vin.Vout)
				fmt.Printf("\t\t%s%x\n", "ScriptSig:", vin.ScriptSig)
				fmt.Println("")
			}
			fmt.Printf("\t%s\n", "TXOutputs:")
			for _, vout := range transaction.Vout {
				fmt.Printf("\t\t%s%d\n", "Value:", vout.Value)
				fmt.Printf("\t\t%s%s\n", "ScriptPubKey:", vout.ScriptPubKey)
				fmt.Println("")
			}
			fmt.Println("")
		}
		fmt.Println("-------------------------")

		if len(block.PrevHash) == 0 {
			break
		}
	}

}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t addblock -data [data] - add a new block")
	fmt.Println("\t printBlocks - print all blocks")

}

func (cli *CLI) Run() {
	cli.validateArgs()

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("printBlocks", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)

	balanceAddress := balanceCmd.String("address", "", "Address to getBalance")
	sendAmount := sendCmd.Int("amount", 0, "Amount to add")
	sendTo := sendCmd.String("to", "", "Where to add")
	sendFrom := sendCmd.String("from", "", "Address from send")

	switch os.Args[1] {
	case "send":
		sendCmd.Parse(os.Args[2:])
	case "getBalance":
		balanceCmd.Parse(os.Args[2:])
	case "printBlocks":
		printBlockCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if sendCmd.Parsed() {
		if *sendAmount != 0 {
			cli.send(*sendFrom, *sendTo, *sendAmount)
		} else {
			cli.printUsage()
		}

	}
	if balanceCmd.Parsed() {
		if *balanceAddress != "" {
			cli.getBalance(*balanceAddress)
		}
	}
	if printBlockCmd.Parsed() {
		cli.printData()
	}

}
