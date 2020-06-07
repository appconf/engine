package memkv

import (
	"strings"
	"sync"
)

//MemKv 一个内存式的kv数据存储
type MemKv struct {
	funcMap map[string]interface{}
	store   map[string]KVPair
	lock    *sync.RWMutex
}

//Get 获取指定key的KVPair, 而不是对应key的直接的值
func (kv *MemKv) Get(key string) (value KVPair, err error) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	value, ok := kv.store[key]
	if !ok {
		err = &KeyError{Key: key, Err: ErrNotExist}
	}
	return
}

//GetValue 获取指定key的值
func (kv *MemKv) GetValue(key string) (value interface{}, err error) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	v, ok := kv.store[key]
	if !ok || v.Value == nil {
		err = &KeyError{Key: key, Err: ErrNotExist}
	}
	return v.Value, err
}

//Gets 批量获取
func (kv *MemKv) Gets(keys []string) (values []KVPair, err error) {
	values = make([]KVPair, 0)
	for _, key := range keys {
		value, err := kv.Get(key)
		if err != nil {
			return values, err
		}
		values = append(values, value)
	}
	return
}

//GetPrefix 根据key的前缀进行批量获取
func (kv *MemKv) GetPrefix(prefix string) (values []KVPair) {
	values = make([]KVPair, 0)
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	for key, value := range kv.store {
		if strings.HasPrefix(key, prefix) {
			values = append(values, value)
		}
	}
	return
}

//Set 存储kv
func (kv *MemKv) Set(key string, value interface{}) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	kv.store[key] = KVPair{Key: key, Value: value}
}

//Exist 判断key是否存在，并且要求key的值不为nil
func (kv *MemKv) Exist(key string) bool {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	value, ok := kv.store[key]
	return value.Value != nil && ok
}

//Equal 判断key的值是否与value相等
func (kv *MemKv) Equal(key string, value interface{}) bool {
	kv.lock.RLock()
	defer kv.lock.RUnlock()

	v, ok := kv.store[key]
	if !ok {
		return false
	}

	return v.Value == value
}

//Delete 删除key
func (kv *MemKv) Delete(key string) {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	delete(kv.store, key)
}

//Clear 清空存储
func (kv *MemKv) Clear() {
	kv.lock.Lock()
	defer kv.lock.Unlock()

	kv.store = make(map[string]KVPair)
}

//FuncMap 获取注册的函数。主要用template向对应的MemKv实例查询数据
func (kv *MemKv) FuncMap() map[string]interface{} {
	return kv.funcMap
}

//New 新建MemKv
func New() *MemKv {
	m := &MemKv{
		funcMap: make(map[string]interface{}),
		store:   make(map[string]KVPair),
		lock:    new(sync.RWMutex),
	}

	m.funcMap["get"] = m.Get
	m.funcMap["getv"] = m.GetValue
	m.funcMap["gets"] = m.Gets
	m.funcMap["getp"] = m.GetPrefix
	m.funcMap["exist"] = m.Exist

	return m
}

var defaultKv = New()

//GetKv 获取默认的MemKv实例
func GetKv() *MemKv {
	return defaultKv
}
