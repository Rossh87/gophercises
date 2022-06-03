package database

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

type PathData struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func SetupDB() (*bolt.DB, error) {
	db, err := bolt.Open("path.db", 0600, nil)

	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("PATH"))

		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}

	fmt.Println("DB Setup Done")
	return db, nil
}

func SetPathData(db *bolt.DB, pathData []PathData) error {
	pathBytes, err := json.Marshal(pathData)

	if err != nil {
		return fmt.Errorf("could not marshal JSON: %v", pathData)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		err = tx.Bucket([]byte("PATH")).Put([]byte("PATHS"), pathBytes)

		if err != nil {
			return fmt.Errorf("could not set paths: %v", err)
		}

		return nil
	})

	fmt.Println("set paths")
	return err
}

func GetPathData(db *bolt.DB) ([]byte, error) {
	var result []byte

	err := db.View(func(t *bolt.Tx) error {
		result = t.Bucket([]byte("PATH")).Get([]byte("PATHS"))
		return nil
	})

	if err != nil {
		return result, fmt.Errorf("failed to retrieve path data from db")
	}

	return result, nil
}
