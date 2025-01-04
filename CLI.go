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

func (cli *CLI) addBlock(data string) {
	cli.bc.addBlock(data)
	fmt.Printf("Added block %s\n", data)
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
		fmt.Printf("%s %s\n", "Block Data:", block.Data)
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

	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printBlockCmd := flag.NewFlagSet("printBlocks", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Data to add")

	switch os.Args[1] {
	case "addBlock":
		addBlockCmd.Parse(os.Args[2:])
	case "printBlocks":
		printBlockCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData != "" {
			cli.addBlock(*addBlockData)
		} else {
			cli.printUsage()
		}

	}
	if printBlockCmd.Parsed() {
		cli.printData()
	}

}
