package cache

import (
	"errors"
	"time"
)

var (
	//ErrNotFound not found
	ErrNotFound = errors.New("not found")
)

//Entity 缓存实体
type Entity struct {
	Key        string        `json:"key"`
	Value      []byte        `json:"value"`
	Expiration time.Duration `json:"expiration"`
}

//Cache 缓存接口
type Cache interface {
	Set(e ...*Entity) error
	Get(key ...string) ([]*Entity, error)
	Del(key ...string) error
	List() ([]*Entity, error)
}
