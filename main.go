package main

func main() {
	ourBlockchain := NewBlockchain()

	cli := CLI{ourBlockchain}
	cli.Run()

	//ourBlockchain.addBlock("Ura mi sdelali blockchain")
	//block := ourBlockchain.Iterator().Next()
	//
	//fmt.Println(block.Hash)
	//fmt.Println(block.Data)

	//ourBlockchain.addBlock("Second block")
	//ourBlockchain.addBlock("Third block")
	//ourBlockchain.addBlock("Fourth block")

	//for _, item := range ourBlockchain.blocks {
	//	fmt.Printf("%s %x\n", "Block hash:", item.Hash)
	//	fmt.Printf("%s %d\n", "Block Timestamp:", item.Timestamp)
	//	fmt.Printf("%s %x\n", "Block Previous Hash:", item.PrevHash)
	//	fmt.Printf("%s %x\n", "Block Data:", item.Data)
	//	fmt.Println("---------")
	//}
	//target := big.NewInt(1)
	//target.Lsh(target, uint(232))
	//targetByte := []byte(target.String())
	//fmt.Printf("%d\n", len(sha256.Sum256(targetByte)))
	//fmt.Printf("%s\n", target)
	//fmt.Println("Hello, world!")
	//Block{2121,{"d"},{},{}} f;
	//newBlock := newBlock("Hello world!", []byte{})
	//newBlock.SetHash()
}
