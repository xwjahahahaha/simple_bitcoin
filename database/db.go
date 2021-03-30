package database

import (
	"github.com/boltdb/bolt"
	"log"
)

type Dao struct {
	DataBase *bolt.DB
	Bucket *bolt.Bucket
}

// Open bolt database
func OpenDatabase(dbName string) *Dao {
	db, err := bolt.Open(dbName, 0666, nil)
	if err != nil {
		log.Panic(err)
		return nil
	}
	return &Dao{
		DataBase: db,
		Bucket:   nil,
	}
}

// create a bucket in a blot database
func (dao *Dao) CreateBucket(bucket string) {
	if dao.DataBase == nil {
		log.Panic("DataBase need created!")
	}
	err := dao.DataBase.Update(func(tx *bolt.Tx) error {
		NewBucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Panic(err)
		}
		dao.Bucket = NewBucket
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// put or update an object
// parameter3 is about serialize val function
func (dao *Dao) ObjectPut(key, value string, fn func(val string) []byte) {
	if dao.Bucket == nil {
		log.Panic("Bucket need assigned!")
	}
	err := dao.Bucket.Put([]byte(key), fn(value))
	if err != nil {
		log.Panic(err)
	}
}

func (dao *Dao) ObjectPurePut(key, value []byte) {
	if dao.Bucket == nil {
		log.Panic("Bucket need assigned!")
	}
	err := dao.Bucket.Put(key, value)
	if err != nil {
		log.Panic(err)
	}
}


// delete an object
func (dao *Dao) ObjectDel(key string) {
	if dao.Bucket == nil {
		log.Panic("Bucket need assigned!")
	}
	err := dao.Bucket.Delete([]byte(key))
	if err != nil {
		log.Panic(err)
	}
}

// query an object
// parameter 2 is about deserialize val function
func (dao *Dao) ObjectGet(key string, fn func(val []byte) string) string {
	if dao.Bucket == nil {
		log.Panic("Bucket need assigned!")
	}
	bytes := dao.Bucket.Get([]byte(key))
	return fn(bytes)
}

func (dao *Dao) ObjectPureGet(key []byte) []byte {
	if dao.Bucket == nil {
		log.Panic("Bucket need assigned!")
	}
	return dao.Bucket.Get(key)
}