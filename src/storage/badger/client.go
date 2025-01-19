package badger

import (
	"encoding/json"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
)

type Client struct {
	db *badger.DB
}

func NewClient(dbPath string) (*Client, error) {
	// DB 디렉토리가 없으면 생성
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, err
	}

	// BadgerDB 옵션 설정
	options := badger.DefaultOptions(absPath)
	options.Logger = nil // 기본 로깅 비활성화

	// DB 열기
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Set stores a key-value pair
func (c *Client) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

// Get retrieves a value by key
func (c *Client) Get(key string, value interface{}) error {
	return c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, value)
		})
	})
}

// Delete removes a key-value pair
func (c *Client) Delete(key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// List returns all keys with a given prefix
func (c *Client) List(prefix string) ([]string, error) {
	var keys []string
	err := c.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixBytes := []byte(prefix)

		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			key := string(item.Key())
			keys = append(keys, key)
		}
		return nil
	})
	return keys, err
}
