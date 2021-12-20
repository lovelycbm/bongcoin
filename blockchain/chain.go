package blockchain

import (
	"encoding/json"
	"net/http"
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

// 블록체인 전체에 대한 구조
// blockchain struct를 b라는 값에다가 싱글턴으로 해서
// DB에서 가져온 값을 b 에다가 넣은 후
// 그 b 값을 다양한 함수에서 직접 추가 및 수정하도록 함.

const (
	defaultDifficulty int = 2
	difficultyAdjustmentInterval int = 5
	blockInterval int = 2
	allowedRange int = 2
)

type blockchain struct {
	NewestHash	string `json:"newestHash"`
	Height  int	`json:"height"`
	CurrentDifficulty int `json:"currentDifficulty"`
	m sync.Mutex
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b,data)
	
}
func (b *blockchain) AddBlock() *Block{
	// 새로운 블록을 저장할때 data, blocks 버켓 두군데에다가 저장.	
	block := cretaeBlock(b.NewestHash, b.Height+1,getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
	return block
}

func persistBlockchain(b *blockchain){
	// blockchain struct를 byte로 변환하여 db에 저장.	
	db.SaveCheckPoint(utils.ToBytes(b))
}

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _ , block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}


func FindTx(b *blockchain,targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}
// db를 어째서 이렇게 찾아야 하는가 추후 고민 해볼것.
// 데이터 hash 의 정합성 문제?
func Blocks(b *blockchain) []*Block{
	b.m.Lock()
	defer b.m.Unlock()
	var blocks []*Block
	hashCursor:= b.NewestHash
	for {		
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			// 해시커서를 찾은 블록의 prevHash로 변경
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func recalcuateDifficulty(b *blockchain) int {
	// 이전 블록의 난이도를 가져와서 현재 난이도를 계산하는 함수
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyAdjustmentInterval - 1]
	actualTIme := (newestBlock.Timestamp/60) - (lastRecalculatedBlock.Timestamp/60)
	expectedTime := blockInterval * difficultyAdjustmentInterval

	if actualTIme <= (expectedTime-allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTIme >= (expectedTime+allowedRange) {
		return b.CurrentDifficulty - 1
	} 	
	return b.CurrentDifficulty	
}

func getDifficulty(b *blockchain) int {
	if b.Height ==0 {
		return defaultDifficulty
	} else if b.Height % difficultyAdjustmentInterval == 0 {
		// 난이도 재조정
		return recalcuateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func UTxOutsByAddress(address string,b *blockchain) []*UTxOut {

	var uTxOuts []*UTxOut

	creatorTxs := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE"{
					break
				}
				if address == FindTx(b,input.TxID).TxOuts[input.Index].Address {
					creatorTxs[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				if address == output.Address {
					if 	_, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID,index,output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}												
					}
				}			
				
			}
		}
	}
	return uTxOuts	
}

func BalanceByAddress(address string,b *blockchain) int {
	var balance int
	txOuts := UTxOutsByAddress(address,b)
	for _, txOut := range txOuts {
		balance += txOut.Amount
	}
	return balance
}

func BlockChain() *blockchain {	
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {				
			b.restore(checkpoint)
		}
	})
	
	return b
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	json.NewEncoder(rw).Encode(b)
}




func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	db.EmptyBlocks()

	for _, block := range newBlocks {
		// persistBlock(block)
		persistBlock(block)

	}
}

func (b *blockchain)AddPeerBlock(newBlock *Block){
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height +=1
	b.CurrentDifficulty = newBlock.Difficulty
	b.NewestHash = newBlock.Hash
	persistBlockchain(b)
	persistBlock(newBlock)

	// mempool 동기화가 아직 안됨. 
	// 가령 채굴이 진행되면 어떤 mem에는 txs값이 있고 어떤 node에는 없을수 있음
	// 모든 노드의 mem을 초기화 할 필요있음.

	for _, tx := range newBlock.Transactions {
		_, ok := m.Txs[tx.ID]
		if ok {
			delete(m.Txs, tx.ID)
		}
	}
}
