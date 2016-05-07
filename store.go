package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

var (
	BucketNotFound = errors.New("Bucket Not Found")
)

func init() {
	var err error
	db, err = bolt.Open("./db/ranker.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	sigch := make(chan os.Signal)
	go func() {
		signal.Notify(sigch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGKILL, syscall.SIGTERM)
		<-sigch
		db.Close()
		os.Exit(0)
	}()
}

type Persister func(users <-chan User)

func storeUser(users <-chan User) {
	for user := range users {
		if err := StoreUser(user.LT, &user); err != nil {
			log.Printf("store user error: %s, user: %v\n", err, user)
		}
	}
}

// 写入到指定bucket的指定key, val中
func Store(bucket string, key, val []byte) error {
	return db.Update(func(bx *bolt.Tx) error {
		bk, err := bx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return bk.Put(key, val)
	})
}

// StoreUser store user into the bucket
func StoreUser(bucket string, user *User) error {
	key := []byte(user.ID)
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return Store(bucket, key, val)
}

func Delete(bucket string, key []byte) error {
	return db.Update(func(bx *bolt.Tx) error {
		bt := bx.Bucket([]byte(bucket))
		if bt == nil {
			return BucketNotFound
		}
		return bt.Delete(key)
	})
}

func DeleteUser(bucket string, user *User) error {
	key := []byte(user.ID)
	return Delete(bucket, key)
}

// Load read key from provided bucket
func Load(bucket string, key []byte) ([]byte, error) {
	var data []byte
	err := db.View(func(bx *bolt.Tx) error {
		bt := bx.Bucket([]byte(bucket))
		if bt == nil {
			return BucketNotFound
		}
		data = bt.Get(key)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func LoadUser(bucket, id string) *User {
	data, err := Load(bucket, []byte(id))
	if err != nil {
		log.Printf("Load user error: %s, bucket: %s, id: %s\n", err, bucket, id)
		return nil
	}
	val := new(User)
	err = json.Unmarshal(data, val)
	if err != nil {
		log.Printf("UnMashal user error: %s, user: %s\n", err, string(data))
		return nil
	}
	return val
}

//
func LoadBucket(bucket string) ([][]byte, error) {
	var data [][]byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return BucketNotFound
		}
		b.ForEach(func(k, v []byte) error {
			data = append(data, v)
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
