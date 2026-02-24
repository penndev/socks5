package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

const bucketName = "data"

type Storage struct {
	db *bolt.DB
}

func NewStorage() (*Storage, error) {
	upath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dbDir := filepath.Join(upath, "Socks5Desktop")
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dbDir, "data.db")
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Set(key string, value any) error {
	if key == "" {
		return errors.New("key不能为空")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Put([]byte(key), data)
	})
}

func (s *Storage) Get(key string) (any, error) {
	if key == "" {
		return nil, errors.New("key不能为空")
	}
	var result any
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(key))
		if v == nil {
			result = nil
			return nil
		}
		data := make([]byte, len(v))
		copy(data, v)
		if err := json.Unmarshal(data, &result); err != nil {
			result = nil
		}
		return nil
	})
	return result, err
}

func (s *Storage) Delete(key string) error {
	if key == "" {
		return errors.New("key不能为空")
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete([]byte(key))
	})
}

func (s *Storage) Close() error {
	return s.db.Close()
}
