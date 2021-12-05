package main

import (
	"github.com/lovelycbm/bongcoin/explorer"
	"github.com/lovelycbm/bongcoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
