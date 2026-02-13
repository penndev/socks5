package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
)

type Storage struct {
	ctx context.Context
	db  *badger.DB
}

// Wails 启动钩子
func (s *Storage) Startup(ctx context.Context) {
	s.ctx = ctx
}

//
// ===== 初始化 =====
//

func NewStorage() (*Storage, error) {
	upath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(upath, "Socks5Desktop")

	opts := badger.DefaultOptions(dbPath).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

//
// ========== Set (JS 可调用) ==========
//

func (s *Storage) Set(key string, value any) error {
	if key == "" {
		return errors.New("key不能为空")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

//
// ========== Get (JS 可调用) ==========
//

func (s *Storage) Get(key string) (any, error) {
	if key == "" {
		return nil, errors.New("key不能为空")
	}

	var result any

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				result = nil
				return nil
			}
			return err
		}

		data, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &result); err != nil {
			// 兼容旧版 gob 格式数据，解析失败时返回空
			result = nil
		}
		return nil
	})

	return result, err
}

//
// ========== Delete (JS 可调用) ==========
//

func (s *Storage) Delete(key string) error {
	if key == "" {
		return errors.New("key不能为空")
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

//
// ========== Close（可选） ==========
//

func (s *Storage) Close() error {
	return s.db.Close()
}
