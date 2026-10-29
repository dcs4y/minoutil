package leveldbclient

import (
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// 写优于读，有文件合并，无事务，数据量大了性能有影响。读多的场景使用boltdb。
var (
	client *leveldbClient
	lock   sync.RWMutex
)

func NewClient(path string) *leveldbClient {
	lock.RLock()
	defer lock.RUnlock()
	if client != nil {
		return client
	}
	// https://github.com/syndtr/goleveldb
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Println("leveldb打开文件失败：" + err.Error())
		return nil
	}
	client = &leveldbClient{db: db}
	return client
}

func GetClient() *leveldbClient {
	return client
}

type leveldbClient struct {
	db   *leveldb.DB
	lock sync.RWMutex
}

// Close 关闭数据库
func (client *leveldbClient) Close() {
	client.db.Close()
}

// Get 根据key查询值
func (client *leveldbClient) Get(key string) ([]byte, error) {
	return client.db.Get([]byte(key), nil)
}

// GetString 根据key查询string值
func (client *leveldbClient) GetString(key string) (string, error) {
	value, err := client.db.Get([]byte(key), nil)
	return string(value), err
}

// GetMap 查询数据库所有数据
func (client *leveldbClient) GetMap() (map[string][]byte, error) {
	client.lock.RLock()
	defer client.lock.RUnlock()
	result := make(map[string][]byte)
	it := client.db.NewIterator(nil, nil)
	for it.Next() {
		result[string(it.Key())] = append(result[string(it.Key())], it.Value()...)
	}
	it.Release()
	err := it.Error()
	return result, err
}

// GetMapString 查询数据库所有数据，值为string
func (client *leveldbClient) GetMapString() (map[string]string, error) {
	client.lock.RLock()
	defer client.lock.RUnlock()
	result := make(map[string]string)
	it := client.db.NewIterator(nil, nil)
	for it.Next() {
		result[string(it.Key())] = string(it.Value())
	}
	it.Release()
	err := it.Error()
	return result, err
}

// GetMapPrefix 根据前缀查询值
func (client *leveldbClient) GetMapPrefix(prefix string) (map[string][]byte, error) {
	client.lock.RLock()
	defer client.lock.RUnlock()
	result := make(map[string][]byte)
	it := client.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for it.Next() {
		result[string(it.Key())] = append(result[string(it.Key())], it.Value()...)
	}
	it.Release()
	err := it.Error()
	return result, err
}

// GetMapPrefixString 根据前缀查询值，值为string
func (client *leveldbClient) GetMapPrefixString(prefix string) (map[string]string, error) {
	client.lock.RLock()
	defer client.lock.RUnlock()
	result := make(map[string]string)
	it := client.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for it.Next() {
		result[string(it.Key())] = string(it.Value())
	}
	it.Release()
	err := it.Error()
	return result, err
}

// GetMapRange 根据范围 [Start,Limit) 查询值
func (client *leveldbClient) GetMapRange(start, end string) (map[string][]byte, error) {
	client.lock.RLock()
	defer client.lock.RUnlock()
	result := make(map[string][]byte)
	it := client.db.NewIterator(&util.Range{Start: []byte(start), Limit: []byte(end)}, nil)
	for it.Next() {
		result[string(it.Key())] = append(result[string(it.Key())], it.Value()...)
	}
	it.Release()
	err := it.Error()
	return result, err
}

// GetMapRangeString 根据范围 [Start,Limit) 查询值，值为string
func (client *leveldbClient) GetMapRangeString(start, end string) (map[string]string, error) {
	client.lock.RLock()
	defer client.lock.RUnlock()
	result := make(map[string]string)
	it := client.db.NewIterator(&util.Range{Start: []byte(start), Limit: []byte(end)}, nil)
	for it.Next() {
		result[string(it.Key())] = string(it.Value())
	}
	it.Release()
	err := it.Error()
	return result, err
}

// Save 保存数据
func (client *leveldbClient) Save(key string, value []byte) error {
	return client.db.Put([]byte(key), value, nil)
}

// SaveString 保存string数据
func (client *leveldbClient) SaveString(key string, value string) error {
	return client.db.Put([]byte(key), []byte(value), nil)
}

// SaveBatch 批量保存数据
func (client *leveldbClient) SaveBatch(batch map[string][]byte) error {
	bat := new(leveldb.Batch)
	for k, v := range batch {
		bat.Put([]byte(k), v)
	}
	return client.db.Write(bat, nil)
}

// SaveBatchString 批量保存string数据
func (client *leveldbClient) SaveBatchString(batch map[string]string) error {
	bat := new(leveldb.Batch)
	for k, v := range batch {
		bat.Put([]byte(k), []byte(v))
	}
	return client.db.Write(bat, nil)
}

// Delete 删除数据
func (client *leveldbClient) Delete(key string) error {
	return client.db.Delete([]byte(key), nil)
}
