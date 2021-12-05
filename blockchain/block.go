package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/lovelycbm/bongcoin/db"
	"github.com/lovelycbm/bongcoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) toByte() []byte{
	var blockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&blockBuffer)
	utils.HandleError(encoder.Encode(b))
	return blockBuffer.Bytes()
}


func(b *Block) persist() {
	db.SaveBlock(b.Hash, b.toByte())
}

func cretaeBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)	
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte( payload)))	
	block.persist()
	return block
}