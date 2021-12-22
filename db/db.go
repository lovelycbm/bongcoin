package db

import (
	"fmt"
	"os"

	"github.com/lovelycbm/bongcoin/utils"
	bolt "go.etcd.io/bbolt"
)

const (
	dbName = "blockchain"
	dataBucketName = "data"
	blocksBucketName = "blocks"
	checkpoint = "checkpoint"
)
var db *bolt.DB

type DB struct {}

func (DB) FindBlock(hash string) []byte {
	return findBlock(hash)
}
func (DB) LoadChain() []byte{
	return loadChain()
}
func (DB) SaveBlock(hash string, data []byte) {
	saveBlock(hash, data)
}
func (DB) SaveChain(data []byte) {
	saveChain(data)
}
func (DB) DeleteAllBlocks() {
	emptyBlocks();
}

func getDbName() string{
	port:= os.Args[1][6:]	
	return fmt.Sprintf("%s_%s.db", dbName,port)
	
}

func InitDB() {
	if db == nil {
		// fmt.Println((getDbName()));
		dbPointer, err := bolt.Open(getDbName(), 0600, nil)
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
	// return db
}

func Close() {
	db.Close()
}

func saveBlock(hash string , data []byte)  {
	// fmt.Printf("Saveing Block %s\nData: %b\n", hash, data)
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		err := bucket.Put([]byte(hash), data)	
		return err
	})
	utils.HandleError(err)
}

func saveChain(data []byte)  {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucketName))
		err := bucket.Put([]byte(checkpoint), data)	
		return err
	})
	utils.HandleError(err) 
}

func loadChain() []byte{
	var data []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucketName))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})	
	return data
}

func findBlock(hash string) []byte{
	var data []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		data = bucket.Get([]byte(hash))
		return nil
	})	
	return data
}

func emptyBlocks() {
	db.Update(func(tx *bolt.Tx) error {
		utils.HandleError(tx.DeleteBucket([]byte(blocksBucketName)))
		_ ,err := tx.CreateBucket([]byte(blocksBucketName))
		utils.HandleError(err)
		return nil
	})	
}