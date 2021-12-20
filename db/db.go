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

func getDbName() string{
	port:= os.Args[1][6:]	
	return fmt.Sprintf("%s_%s.db", dbName,port)
	// for i, a := range os.Args {
	// 	fmt.Println(i,a)
	// }
}

func DB() *bolt.DB{
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
	return db
}

func Close() {
	DB().Close()
}



func SaveBlock(hash string , data []byte)  {
	// fmt.Printf("Saveing Block %s\nData: %b\n", hash, data)
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		err := bucket.Put([]byte(hash), data)	
		return err
	})
	utils.HandleError(err)
}

func SaveCheckPoint(data []byte)  {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucketName))
		err := bucket.Put([]byte(checkpoint), data)	
		return err
	})
	utils.HandleError(err) 
}

func Checkpoint() []byte{
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucketName))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})	
	return data
}

func Block(hash string) []byte{
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucketName))
		data = bucket.Get([]byte(hash))
		return nil
	})	
	return data
}

func EmptyBlocks() {
	DB().Update(func(tx *bolt.Tx) error {
		utils.HandleError(tx.DeleteBucket([]byte(blocksBucketName)))
		_ ,err := tx.CreateBucket([]byte(blocksBucketName))
		utils.HandleError(err)
		return nil
	})	
}