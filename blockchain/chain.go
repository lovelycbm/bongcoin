package blockchain

import (
	"sync"
)

// 블록체인의 기본 구조
// B1
// 	b1Hash  = (data+ "")
// B2
// 	b2Hash = (data + b1Hash)
// B3
// 	b3Hash = (data + b2Hash)

/// func 앞에 모양은 구조체에서 쓰이는 method라고 선언하는걸로 이해하면 됨.


type blockchain struct {
	NewestHash	string `json:"newestHash"`
	Height  int	`json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := cretaeBlock(data,b.NewestHash, b.Height)
	b.NewestHash = block.Hash
	b.Height = b.Height
	// block := Block{data, "", b.NewestHash, b.Height + 1}
	
}


func BlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis")
		})
	}
	return b
}



