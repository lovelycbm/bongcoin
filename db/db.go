package db

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/lovelycbm/bongcoin/utils"
)

const (
	dbName = "blockchain.db"
	dataBucketName = "data"
	blocksBucketName = "blocks"
)
var db *bolt.DB

func DB() *bolt.DB{
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		utils.HandleError(err)
		
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucketName))
			utils.HandleError(err) 
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucketName))
			return err
		})
		utils.HandleError(err)	
	}
	return db
}

func SaveBlock(hash string , data []byte)  {
	fmt.Printf("Saveing Block %s\n Data: %b\n", hash, data)
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		err := bucket.Put([]byte(hash), data)	
		return err
	})
	utils.HandleError(err) 
}

func SaveBlockChain(data []byte)  {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucketName))
		err := bucket.Put([]byte("checkpoint"), data)	
		return err
	})
	utils.HandleError(err) 
}
