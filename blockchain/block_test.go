package blockchain

import (
	"reflect"
	"testing"

	"github.com/lovelycbm/bongcoin/utils"
)



func TestCreatedBlock(t *testing.T) {
	dbStorage = testDB{}
	Mempool().Txs["test"] = &Tx{}
	b := cretaeBlock("x",1,1)

	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}){
		t.Error("createBlock() should return an intance of block")
	}
}

func TestFindBlock(t *testing.T){
	t.Run("Block is found", func(t *testing.T){
		dbStorage = testDB{
			fakeFindBlock: func() []byte {
				b:= &Block{
					Height : 1,
				}
				return utils.ToBytes(b)
			},
		}
		block, _  := FindBlock("xx")
		if block.Height != 1{
			t.Error("Block should be found")
		}
	})
}