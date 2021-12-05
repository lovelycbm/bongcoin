package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
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

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height  int	`json:"height "`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func (b *Block) getHash() {
	hash := sha256.Sum256([]byte(b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)

	if totalBlocks == 0 {
		return ""
	}

	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(),len(GetBlockChain().blocks)+1 }
	newBlock.getHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}

var ErrBlockNotFound = errors.New("block not found")

func (b *blockchain) GetBlock( height int) (*Block,error){
	if height > len(b.blocks) {
		return nil, ErrBlockNotFound
	}

	return b.blocks[height-1] ,nil
}