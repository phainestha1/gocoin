package main

import (
	"gocoin/cli"
	"gocoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}