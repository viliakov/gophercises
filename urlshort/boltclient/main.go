package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("../main/urlShortener.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("mappings"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("mappings"))
		err := b.Put([]byte("/boltdb"), []byte("https://github.com/boltdb/bolt"))
		return err
	})
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("mappings"))
		err := b.Put([]byte("/google"), []byte("https://google.com"))
		return err
	})
}
