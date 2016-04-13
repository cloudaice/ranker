package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/boltdb/bolt"
)

var (
	db *bolt.DB
)

func init() {
	var err error
	db, err = bolt.Open("ranker.db", 0600, nil)
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

func StoreWorker() {
	ch := make(chan User, 100)
	stop := make(chan struct{})
	up := NewUserUpdater("world", ch)
	go func() {
		up.Start()
		stop <- struct{}{}
	}()

	go func() {
		for i := 1; i < 11; i++ {
			log.Println(i)
			users := SearchUser("world", i, "followers")
			for _, user := range users {
				ch <- user
			}
		}
		close(ch)
	}()
	<-stop
}

type UserUpdater struct {
	input  chan User
	bucket string
}

func NewUserUpdater(bucket string, input chan User) *UserUpdater {
	return &UserUpdater{
		input:  input,
		bucket: bucket,
	}
}

func (u *UserUpdater) Start() {
	for user := range u.input {
		key := []byte(user.ID)
		val, err := json.Marshal(user)
		if err != nil {
			log.Println("Mashal error", err)
			continue
		}
		if err := Store(u.bucket, key, val); err != nil {
			log.Println("Store error", err)
		}
	}
}

func Store(bucket string, key, val []byte) error {
	return db.Update(func(bx *bolt.Tx) error {
		bk, err := bx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return bk.Put(key, val)
	})
}

func Load(bucket string, key []byte) ([]byte, error) {
	var data []byte
	err := db.View(func(bx *bolt.Tx) error {
		bt := bx.Bucket([]byte(bucket))
		data = bt.Get(key)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// StoreCityMap store city dict
func StoreCityMap() error {
	val, err := json.Marshal(CityList)
	if err != nil {
		return err
	}
	return Store("meta", []byte("cityMap"), val)
}

// LoadCityMap ...
func LoadCityMap(val map[string]string) error {
	data, err := Load("meta", []byte("cityMap"))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}

// StoreCountryMap store country dict
func StoreCountryMap() error {
	val, err := json.Marshal(CountryMap)
	if err != nil {
		return err
	}
	return Store("meta", []byte("countryMap"), val)
}

// LoadCountryMap
func LoadCountryMap(val map[string]string) error {
	data, err := Load("meta", []byte("countryMap"))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, val)
}

//
func LoadBucket(bucket string) ([][]byte, error) {
	var data [][]byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
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
