package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/lovelycbm/bongcoin/db"
	"github.com/lovelycbm/bongcoin/utils"
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

func (b *blockchain) restore(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utils.HandleError(decoder.Decode(b))
}

func (b *blockchain) persist(){
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	// 새로운 블록을 저장할때 data, blocks 버켓 두군데에다가 저장.	
	block := cretaeBlock(data,b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()	
}


func BlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}		
			fmt.Printf("NewestHash %s\nHeight %d\n", b.NewestHash, b.Height)
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				fmt.Println("Loaded blockchain from db:")
				b.restore(checkpoint)
			}
		})
	}
	fmt.Printf("NewestHash %s\nHeight:%d\n", b.NewestHash, b.Height)
	return b
}



