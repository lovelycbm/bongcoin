package blockchain

import (
	"reflect"
	"sync"
	"testing"

	"github.com/lovelycbm/bongcoin/utils"
)

type testDB struct {
	fakeLoadChain func() []byte
	fakeFindBlock func() []byte
}

func (t testDB) FindBlock(hash string) []byte {
	return t.fakeFindBlock()
}
func (t testDB) LoadChain() []byte {
	return t.fakeLoadChain()
}
func (testDB) SaveBlock(hash string, data []byte) {}
func (testDB) SaveChain(data []byte)              {}
func (testDB) DeleteAllBlocks()                   {}

func TestBlockchain(t *testing.T) {
	t.Run("Should create blockchain",func(t *testing.T){
		dbStorage = testDB{
			fakeLoadChain: func() []byte {
				return nil
			},
		}
		bc:= BlockChain()
		if bc.Height != 1 {
			t.Error("Blockchain() should create a blockchain with height 1")
		}
	})
	t.Run("Should restore blockchain",func(t *testing.T){
		once = *new(sync.Once)
		dbStorage = testDB{
			fakeLoadChain: func() []byte {
				bc := &blockchain{Height:2, NewestHash:"xxx",CurrentDifficulty:1}
				return utils.ToBytes(bc)
			},
		}
		bc:= BlockChain()
		if bc.Height != 2 {
			t.Errorf("Blockchain() should restore a blockchain with height of %d, got %d",2, bc.Height)
		}
	})
}

func TestBlocks(t *testing.T){
	blocks:= []*Block{		
		{PrevHash: "x",},
		{PrevHash: "",},
	}
	fakeBlocks := 0 
	dbStorage = testDB{
		fakeFindBlock: func() []byte {
			defer func(){
				fakeBlocks++
			}()
			return utils.ToBytes(blocks[fakeBlocks])
		},
	}
	bc:= &blockchain{}
	blocksResult := Blocks(bc)

	if reflect.TypeOf(blocksResult) != reflect.TypeOf([]*Block{}){
		t.Error("Blocks() should return a slice of blocks")
	}
}

func TestFindTx(t *testing.T) {
	t.Run("Tx not found" , func(t *testing.T){

		dbStorage = testDB{
			fakeFindBlock: func() []byte {			
				b := &Block{
					Height : 2,
					// PrevHash : "x",
					Transactions: []*Tx{},
				}	
				return utils.ToBytes(b)
			},
		}

		tx := FindTx(&blockchain{NewestHash:"x"}, "test")
		if tx != nil {
			t.Error("Tx should be not found")
		}
	})
	t.Run("Tx should be found" , func(t *testing.T){

		dbStorage = testDB{
			fakeFindBlock: func() []byte {			
				b := &Block{
					Height : 2,					
					Transactions: []*Tx{
						{ID: "test"},
					},					
				}	
				return utils.ToBytes(b)
			},
		}

		tx := FindTx(&blockchain{NewestHash:"x"}, "test")
		if tx == nil {
			t.Error("Tx should be found")
		}
	})
}

func TestGetDifficulty(t *testing.T){
	blocks:= []*Block{
		{PrevHash: "x",},
		{PrevHash: "x",},
		{PrevHash: "x",},
		{PrevHash: "x",},
		{PrevHash: "",},
	}
	fakeBlock := 0
	dbStorage = testDB{
		fakeFindBlock: func() []byte {			
			defer func(){
				fakeBlock++
			}()
			return utils.ToBytes(blocks[fakeBlock])
		},
	}
	type test struct{
		height int
		want int		
	}
	tests:= []test{
		{height : 0, want:defaultDifficulty},
		{height : 2, want:defaultDifficulty},
		{height : 5, want:3},
	}

	for _, tc := range tests{
		bc := &blockchain{Height:tc.height, CurrentDifficulty:defaultDifficulty}
		got:= getDifficulty(bc)
		if got!= tc.want {
			t.Errorf("getDifficulty() should be %d, got %d",tc.want,got)
		}
	}
}

func TestAddPeerBlock(t *testing.T){
	bc:= &blockchain{
		Height:1,
		CurrentDifficulty:1,
		NewestHash:"xx",
	}
	m.Txs["test"] = &Tx{}
	nb:= &Block{
		Difficulty: 2, 
		Hash : "test",
		Transactions: []*Tx{
			{ID:"test"},
		},
	}
	bc.AddPeerBlock(nb)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("AddPeerBlock() should mutate blockchain")
	}	
	
}

func TestReplace (t *testing.T){
	bc:= &blockchain{
		Height:1,
		CurrentDifficulty:1,
		NewestHash:"xx",
	}

	blocks := []*Block{
		{Difficulty : 2, Hash:"test"},
		{Difficulty : 2, Hash:"test"},
	}

	bc.Replace(blocks)

	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "test" {
		t.Error("AddPeerBlock() should mutate blockchain")
	}	

}