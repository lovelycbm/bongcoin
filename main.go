package main

import (
	"github.com/lovelycbm/bongcoin/cli"
	"github.com/lovelycbm/bongcoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}