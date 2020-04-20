package main

import (
	"encoding/json"
	"errors"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os"
)

func setupDB(path string) (db *bolt.DB, err error) {
	/*_, err = os.Stat(path)
	if err != nil {
		log.Fatalf("sync db is not found on the path: %s", path)
	}*/

	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("can't open db file %s: %s\n", path, err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("Words"))
		if err != nil {
			return errors.New("could not create Words bucket")
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v\n", err)
	}

	return db, nil
}

func cleanupDB(db *bolt.DB) (err error) {
	var fn string

	fn = db.Path()

	err = db.Close()
	if err != nil {
		return err
	}

	err = os.Remove(fn)
	if err != nil {
		return err
	}

	return nil
}

func storeWord(db *bolt.DB, word WordStat) (err error) {
	var data []WordStat
	w, err := getWord(db, word.Text)
	if err != nil {
		if err.Error() != "not found" {
			return err
		} else {
			data = []WordStat{word}
		}
	} else {
		data = w
		data = append(data, word)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = db.Batch(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("Words")).Put([]byte(word.Text), dataBytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func getWord(db *bolt.DB, word string) (w []WordStat, err error) {
	w = make([]WordStat, 1)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Words"))
		v := b.Get([]byte(word))
		if v != nil {
			err = json.Unmarshal(v, &w)
			if err != nil {
				return err
			}
		} else {
			return errors.New("not found")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return w, nil
}
