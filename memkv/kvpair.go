package memkv

import (
	"errors"
)

//KVPair 键值对数据结构
type KVPair struct {
	Key   string
	Value interface{}
}

var (
	//ErrNotExist MemKv key不存在错误
	ErrNotExist = errors.New("key does not exists")
)

//KeyError 用于MemKv的key相关错误
type KeyError struct {
	Key string
	Err error
}

//Error error实现
func (e *KeyError) Error() string {
	return e.Err.Error() + ":" + e.Key
}
