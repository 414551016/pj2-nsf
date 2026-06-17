package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
	berror "go.etcd.io/bbolt/errors"
)

const (
	serviceBucket = "services"
)

type bboltDb struct {
	db *bbolt.DB
}

func newBboltDb(dbPath string) (*bboltDb, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create DB directory: %v", err)
	}

	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %v", err)
	}

	return &bboltDb{
		db: db,
	}, nil
}

func (b *bboltDb) Release() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

func (b *bboltDb) Exist(key string) (bool, error) {
	exist := false
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(serviceBucket))
		if bucket == nil {
			exist = false
			return nil
		}
		v := bucket.Get([]byte(key))
		exist = v != nil
		return nil
	})
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (b *bboltDb) Save(key, value string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(serviceBucket))
		if err != nil {
			return fmt.Errorf("failed to create or access bucket: %v", err)
		}
		return bucket.Put([]byte(key), []byte(value))
	})
}

func (b *bboltDb) SaveAll(data map[string]string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(serviceBucket))
		if err != nil {
			return fmt.Errorf("failed to create or access bucket: %v", err)
		}
		for key, value := range data {
			if err := bucket.Put([]byte(key), []byte(value)); err != nil {
				return fmt.Errorf("failed to save key %s: %v", key, err)
			}
		}
		return nil
	})
}

func (b *bboltDb) Load(key string) (string, error) {
	var value []byte
	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(serviceBucket))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		v := bucket.Get([]byte(key))
		value = make([]byte, len(v))
		copy(value, v)
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (b *bboltDb) LoadAll() (map[string]string, error) {
	result := make(map[string]string)

	err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(serviceBucket))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			value := make([]byte, len(v))
			copy(value, v)
			result[string(k)] = string(value)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *bboltDb) Remove(key string) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(serviceBucket))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}
		return bucket.Delete([]byte(key))
	})
}

func (b *bboltDb) RemoveAll() error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(serviceBucket))
	})
	if errors.Is(err, berror.ErrBucketNotFound) {
		return nil
	}
	return err
}
