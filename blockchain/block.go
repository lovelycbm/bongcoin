package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/lovelycbm/bongcoin/db"
	"github.com/lovelycbm/bongcoin/utils"
)



type Block struct {	
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
	// hash 의 앞부분의 0이 몇개인지를 나타냄 
	Difficulty int  `json:"difficulty"`
	// 채굴자들이 값을 바꿔가며 hash를 변경하여 앞자리를 찾기 위해 쓰는 변수
	Nonce     int   `json:"nonce"`
	Timestamp int `json:"timestamp"`
	Transactions []*Tx `json:"transactions"`
}

// 블록 한개에 대한 구조 
// block struct를 싱글턴으로 해서 
// 값을 다양한 함수에서 직접 추가 및 수정하도록 함. 

func (b *Block) restore(data []byte) {
	utils.FromBytes(b,data)
	
}

func persistBlock(b *Block)  {
	// hash 를 key 로 하는 block struct를 db에 저장함.
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

var ErrBlockNotFound = errors.New("Block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil,ErrBlockNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block,nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)				
		// fmt.Printf("\n\n\nTarget:%s\nHash:%s\n,Nonce:%d\n\n\n",target,hash,b.Nonce)
		if strings.HasPrefix(hash, target) {			
			b.Hash = hash 
			break
		} else {
			b.Nonce++
		}

	}
}

func cretaeBlock( prevHash string, height,diff int) *Block {
	block := &Block{		
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
		Difficulty: diff,
		Nonce:    0,		
	}
	
	
	block.mine()
	block.Transactions = Mempool().TxToConfirm()
	persistBlock(block)
	return block
}