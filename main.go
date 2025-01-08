package main

func main() {
	ourBlockchain := NewBlockchain("12121")

	cli := CLI{ourBlockchain}
	cli.Run()
}
