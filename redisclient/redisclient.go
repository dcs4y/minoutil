package redisclient

import (
	"errors"
	"github.com/dcs4y/minoutil/v2/logutil"
	"github.com/gomodule/redigo/redis"
	"math"
	"strconv"
	"strings"
)

/**
Redis 命令 https://www.runoob.com/redis/redis-commands.html
*/

var log = logutil.GetLog("redis")

var clients = make(map[string]*redisClient)

type RedisConfig struct {
	Host     string
	Port     int
	Database int
	Password string
}

type redisClient struct {
	pool *redis.Pool // 创建redis连接池
}

func NewClient(name string, config RedisConfig) *redisClient {
	// 初始化Redis连接池
	redisPool := &redis.Pool{ //实例化一个连接池
		MaxIdle:     8,   //初始连接数
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			options := make([]redis.DialOption, 1)
			setDB := redis.DialDatabase(config.Database)
			options[0] = setDB
			if config.Password != "" {
				setPassword := redis.DialPassword(config.Password)
				options = append(options, setPassword)
			}
			return redis.Dial("tcp", config.Host+":"+strconv.Itoa(config.Port), options...)
		},
	}
	client := &redisClient{pool: redisPool}
	clients[name] = client
	return client
}

func GetClientByName(name string) *redisClient {
	return clients[name]
}

func GetClient() *redisClient {
	return clients[""]
}

// ClosePool 关闭连接池，慎用。
func (rc *redisClient) ClosePool() {
	rc.pool.Close() //关闭连接池
}

//--------------------------------------------------key string----------------------------------------------------------

// SetObject 设置key的值
func (rc *redisClient) SetObject(key string, value interface{}, second int) error {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	if second > 0 {
		//为防止上锁后redis机器故障，使用set nx ex上锁同时设置过期时间：（原子操作）
		//_, err := rdConn.Do("set", key, value, "nx", "ex", second)
		_, err := rdConn.Do("SETEX", key, second, value)
		return err
	} else {
		_, err := rdConn.Do("Set", key, value)
		return err
	}
}

// SetExpire 单独设置过期时间(秒)
func (rc *redisClient) SetExpire(key string, second int) error {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	if second > 0 {
		_, err := rdConn.Do("expire", key, second)
		return err
	}
	return errors.New("过期时间必须大于0！")
}

// GetKeys 查找所有符合给定模式(pattern|*)的key。
func (rc *redisClient) GetKeys(pattern string) ([]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Strings(rdConn.Do("KEYS", pattern))
}

// Delete 删除KEY
func (rc *redisClient) Delete(key string) (int, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int(rdConn.Do("DEL", key))
}

// GetString 获取key的值_字符串
func (rc *redisClient) GetString(key string) (string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.String(rdConn.Do("Get", key))
}

func (rc *redisClient) GetInt(key string) (int, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int(rdConn.Do("Get", key))
}

func (rc *redisClient) GetInt64(key string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("Get", key))
}

func (rc *redisClient) GetFloat64(key string) (float64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Float64(rdConn.Do("Get", key))
}

func (rc *redisClient) GetBool(key string) (bool, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Bool(rdConn.Do("Get", key))
}

func (rc *redisClient) GetByte(key string) ([]byte, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Bytes(rdConn.Do("Get", key))
}

func (rc *redisClient) GetObject(key string) (interface{}, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return rdConn.Do("Get", key)
}

// Lock 简单锁。只有在 key 不存在时设置 key 的值。删除key即为解锁。
func (rc *redisClient) Lock(key string, value interface{}) (bool, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	m, err := rdConn.Do("SETNX", key, value)
	if m.(int64) == 1 {
		return true, nil
	}
	return false, err
}

// LockExpire 简单锁。只有在 key 不存在时设置 key 的值。删除key即为解锁，否则超时自动解锁。
func (rc *redisClient) LockExpire(key string, value interface{}, second int) (bool, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	m, err := rdConn.Do("SETNX", key, value)
	if m.(int64) == 1 {
		_, err = rdConn.Do("expire", key, second)
		return true, err
	}
	return false, err
}

// SetObjects 设置多key的值
func (rc *redisClient) SetObjects(kv ...interface{}) error {
	if len(kv)%2 != 0 {
		return errors.New("key-value必须成对！")
	}
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	_, err := rdConn.Do("MSET", kv...)
	return err
}

// GetStrings 获取多key的值
func (rc *redisClient) GetStrings(key ...interface{}) (map[interface{}]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	array, err := redis.Strings(rdConn.Do("MGET", key...))
	if err != nil {
		return nil, err
	}
	m := make(map[interface{}]string)
	for i, s := range key {
		m[s] = array[i]
	}
	return m, nil
}

//---------------------------------------------------------hash---------------------------------------------------------

// HashSet 设置hash的值
func (rc *redisClient) HashSet(redisKey string, hashKey string, value interface{}) error {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	_, err := rdConn.Do("HSet", redisKey, hashKey, value)
	return err
}

// HashSetAll 同时设置hash的多个值
func (rc *redisClient) HashSetAll(redisKey string, hashKeyValue ...interface{}) error {
	if len(hashKeyValue)%2 != 0 {
		return errors.New("key-value必须成对！")
	}
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	hashKeyValue = append([]interface{}{redisKey}, hashKeyValue...)
	_, err := rdConn.Do("HMSET", hashKeyValue...)
	return err
}

// HashGetString 获取hash的key对应的值。
func (rc *redisClient) HashGetString(redisKey string, hashKey string) (string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.String(rdConn.Do("HGet", redisKey, hashKey))
}

// HashGetAll 获取hash的所有值
func (rc *redisClient) HashGetAll(redisKey string) (map[string]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.StringMap(rdConn.Do("HGETALL", redisKey))
}

// HashDelete 删除hash的多个key
func (rc *redisClient) HashDelete(redisKey string, hashKey ...string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	hashKey = append([]string{redisKey}, hashKey...)
	return redis.Int64(rdConn.Do("HDEL", stringToInterface(hashKey...)...))
}

// HashLock hash锁。只有在 hashKey 不存在时设置 hashKey 的值。删除hashKey即为解锁。
func (rc *redisClient) HashLock(redisKey string, hashKey string, value interface{}) (bool, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	m, err := rdConn.Do("HSETNX", redisKey, hashKey, value)
	if m.(int64) == 1 {
		return true, nil
	}
	return false, err
}

//-------------------------------------------------------------list-----------------------------------------------------

// ListLeftPush 列表头部插入数据。返回列表长度。
func (rc *redisClient) ListLeftPush(key string, value ...string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	value = append([]string{key}, value...)
	return redis.Int64(rdConn.Do("LPUSH", stringToInterface(value...)...))
}

// ListRightPush 列表尾部插入数据。返回列表长度。
func (rc *redisClient) ListRightPush(key string, value ...string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	value = append([]string{key}, value...)
	return redis.Int64(rdConn.Do("RPUSH", stringToInterface(value...)...))
}

// ListAddBefore 在列表指定位置(pivot元素)之前插入值。返回列表长度，pivot不存在时返回-1。
func (rc *redisClient) ListAddBefore(key string, pivot string, value string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("LINSERT", key, "BEFORE", pivot, value))
}

// ListAddAfter 在列表指定位置(pivot元素)之后插入值。返回列表长度，pivot不存在时返回-1。
func (rc *redisClient) ListAddAfter(key string, pivot string, value string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("LINSERT", key, "AFTER", pivot, value))
}

// ListUpdate 更新列表的指定位置的数据
func (rc *redisClient) ListUpdate(key string, index int64, value string) error {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	_, err := rdConn.Do("LSET", key, index, value)
	return err
}

// ListGetString 通过索引获取列表元素值
func (rc *redisClient) ListGetString(key string, index int64) (string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.String(rdConn.Do("LINDEX", key, index))
}

// ListGetRange 获取指定索引范围内的数据。start开始位置，end结束位置。正向数为正数(从0开始)，负向数为负数(从-1开始)。
func (rc *redisClient) ListGetRange(redisKey string, start, end int64) ([]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Strings(rdConn.Do("LRANGE", redisKey, start, end))
}

// ListGetAll 获取列表的所有数据
func (rc *redisClient) ListGetAll(redisKey string) ([]string, error) {
	return rc.ListGetRange(redisKey, 0, -1)
}

// ListDelete 删除列表中的value
func (rc *redisClient) ListDelete(key string, value string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("LREM", key, 0, value))
}

// ListLeftPopString 移除并获取头部数据
func (rc *redisClient) ListLeftPopString(key string) (string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.String(rdConn.Do("LPOP", key))
}

// ListRightPopString 移除并获取尾部数据
func (rc *redisClient) ListRightPopString(key string) (string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.String(rdConn.Do("RPOP", key))
}

// ListCut 列表裁剪。删除指定范围以外的元素。
func (rc *redisClient) ListCut(key string, start int64, end int64) error {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	_, err := rdConn.Do("LTRIM", key, start, end)
	return err
}

// ListLength 获取列表长度
func (rc *redisClient) ListLength(key string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("LLEN", key))
}

//--------------------------------------------------------------set-----------------------------------------------------

// SetAdd 无序集合_增加元素。返回增加个数。
func (rc *redisClient) SetAdd(key string, value ...string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	value = append([]string{key}, value...)
	return redis.Int64(rdConn.Do("SADD", stringToInterface(value...)...))
}

// SetLength 无序集合_获取集合成员数
func (rc *redisClient) SetLength(key string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("SCARD", key))
}

// SetGetAll 无序集合_获取集合所有元素
func (rc *redisClient) SetGetAll(key string) ([]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Strings(rdConn.Do("SMEMBERS", key))
}

// SetInter 无序集合_获取集合交集
func (rc *redisClient) SetInter(keys ...string) ([]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Strings(rdConn.Do("SINTER", stringToInterface(keys...)...))
}

// SetUnion 无序集合_获取集合并集
func (rc *redisClient) SetUnion(keys ...string) ([]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Strings(rdConn.Do("SUNION", stringToInterface(keys...)...))
}

// SetDeleteRand 无序集合_删除集合中一个随机元素。返回被删除的元素。
func (rc *redisClient) SetDeleteRand(key string) (string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.String(rdConn.Do("SPOP", key))
}

// SetDelete 无序集合_删除集合中的元素。返回删除的元素个数。
func (rc *redisClient) SetDelete(key string, value ...string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	value = append([]string{key}, value...)
	return redis.Int64(rdConn.Do("SREM", stringToInterface(value...)...))
}

// SetIsMember 无序集合_判断 value 元素是否是集合(set) redisKey 的成员。此方法无效。
func (rc *redisClient) SetIsMember(key string, value string) (bool, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Bool(rdConn.Do("SISMEMBER", key, value))
}

// ZSetAdd 有序集合_增加元素。返回增加个数。
func (rc *redisClient) ZSetAdd(key string, scoreValue ...interface{}) (int64, error) {
	if len(scoreValue)%2 != 0 {
		return 0, errors.New("分数和值必须成对！")
	}
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	scoreValue = append([]interface{}{key}, scoreValue...)
	return redis.Int64(rdConn.Do("ZADD", scoreValue...))
}

// ZSetDelete 有序集合_删除元素。返回删除的元素个数。
func (rc *redisClient) ZSetDelete(key string, value ...string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	value = append([]string{key}, value...)
	return redis.Int64(rdConn.Do("ZREM", stringToInterface(value...)...))
}

// ZSetDeleteByScore 有序集合_按分数闭区间删除元素。返回删除的元素个数。
func (rc *redisClient) ZSetDeleteByScore(key string, minScore, maxScore float64) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("ZREMRANGEBYSCORE", key, minScore, maxScore))
}

// ZSetLength 有序集合_获取成员数
func (rc *redisClient) ZSetLength(key string) (int64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int64(rdConn.Do("ZCARD", key))
}

// ZSetGetScore 有序集合_获取元素分数
func (rc *redisClient) ZSetGetScore(key string, value string) (float64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Float64(rdConn.Do("ZSCORE", key, value))
}

// ZSetGetByScore 有序集合_按分数区间获取成员。descending是否倒序，false为正序。
func (rc *redisClient) ZSetGetByScore(key string, minScore, maxScore float64, descending bool) ([]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	if descending {
		// 按分数和按元素字典倒序
		return redis.Strings(rdConn.Do("ZREVRANGEBYSCORE", key, maxScore, minScore))
	} else {
		// 按分数和按元素字典正序
		return redis.Strings(rdConn.Do("ZRANGEBYSCORE", key, minScore, maxScore))
	}
}

// ZSetGetByScoreWithScore 有序集合_按分数区间获取成员和分数。descending是否倒序，false为正序。
func (rc *redisClient) ZSetGetByScoreWithScore(key string, minScore, maxScore float64, descending bool) (map[string]string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	if descending {
		// 按分数和按元素字典倒序
		return redis.StringMap(rdConn.Do("ZREVRANGEBYSCORE", key, maxScore, minScore, "WITHSCORES"))
	} else {
		// 按分数和按元素字典正序
		return redis.StringMap(rdConn.Do("ZRANGEBYSCORE", key, minScore, maxScore, "WITHSCORES"))
	}
}

//------------------------------------------------------------发布订阅---------------------------------------------------
// 发布订阅 (pub/sub) 可以分发消息，但无法记录历史消息。
// 当没有订阅者时，发送的消息将被丢弃。

type redisMessage struct {
	Action  string
	Channel string
	Message string
}

// ChannelSubscribe 发布订阅_订单频道。向stop发送数据则退出当前订阅。channel可使用*匹配。
// redis事件通知时channel格式：__[keyspace|keyevent]@<db>__:[prefix]
func (rc *redisClient) ChannelSubscribe(channel ...string) chan redisMessage {
	ch := make(chan redisMessage)
	go func(channel ...string) {
		defer close(ch)
		rdConn := rc.pool.Get()
		// 函数运行结束 ，把连接放回连接池
		defer rdConn.Close()
		for _, c := range channel {
			if strings.HasPrefix(c, "__key") {
				// 开启单一客户端的redis事件通知。
				log.Print("开启事件通知：")
				log.Println(rdConn.Do("config", "set", "notify-keyspace-events", "KEA"))
			}
		}
		cmd := "SUBSCRIBE" // 通道订阅
		for _, c := range channel {
			if strings.ContainsAny(c, "*") {
				cmd = "PSUBSCRIBE" // 通道名进行匹配订阅
				break
			}
		}
		rdConn.Send(cmd, stringToInterface(channel...)...)
		rdConn.Flush()
		for {
			vs, err := redis.Values(rdConn.Receive())
			if err != nil {
				log.Println("redis订阅失败！", err)
				break
			}
			action := string(vs[0].([]byte))
			switch action {
			case "message":
				rm := redisMessage{
					Action:  string(vs[0].([]byte)),
					Channel: string(vs[1].([]byte)),
					Message: string(vs[2].([]byte)),
				}
				ch <- rm
			case "pmessage":
				rm := redisMessage{
					Action:  string(vs[0].([]byte)),
					Channel: string(vs[2].([]byte)),
					Message: string(vs[3].([]byte)),
				}
				ch <- rm
			case "subscribe", "psubscribe":
				for _, v := range vs {
					switch v.(type) {
					case []byte:
						log.Print(string(v.([]byte)), " ")
					case int64:
						log.Println("订阅频道数：", v)
					default:
						log.Println("未知数据", v)
					}
				}
			default:
				for _, v := range vs {
					switch v.(type) {
					case []byte:
						log.Println("string=", string(v.([]byte)))
					case int64:
						log.Println("int64=", v)
					default:
						log.Println("default=", v)
					}
				}
			}
		}
	}(channel...)
	return ch
}

// ChannelPublish 发布订阅_向通道发送消息。返回通道的订阅数。
func (rc *redisClient) ChannelPublish(channel string, message interface{}) (int, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Int(rdConn.Do("PUBLISH", channel, message))
}

//-----------------------------------------------------------执行lua脚本--------------------------------------------------

// LuaExecute 执行lua脚本
func (rc *redisClient) LuaExecute(luaScript string, redisKeys []string, luaParams []string) (interface{}, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	args := append([]string{luaScript}, strconv.Itoa(len(redisKeys)))
	args = append(args, redisKeys...)
	args = append(args, luaParams...)
	return rdConn.Do("EVAL", stringToInterface(args...)...)
}

// LuaExecuteWithSha1 执行redis中的lua脚本
func (rc *redisClient) LuaExecuteWithSha1(sha1 string, redisKeys []string, luaParams []string) (interface{}, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	args := append([]string{sha1}, strconv.Itoa(len(redisKeys)))
	args = append(args, redisKeys...)
	args = append(args, luaParams...)
	return rdConn.Do("EVALSHA", stringToInterface(args...)...)
}

// LuaLoad 载入lua脚本到redis，但不执行。返回脚本sha1。
func (rc *redisClient) LuaLoad(luaScript string) (string, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.String(rdConn.Do("SCRIPT", "LOAD", luaScript))
}

// LuaExists 查看lua脚本是否存在redis中
func (rc *redisClient) LuaExists(sha1 ...string) (map[string]bool, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	args := append([]string{"EXISTS"}, sha1...)
	rs, err := redis.Ints(rdConn.Do("SCRIPT", stringToInterface(args...)...))
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	for i, r := range rs {
		if r == 1 {
			result[sha1[i]] = true
		} else {
			result[sha1[i]] = false
		}
	}
	return result, nil
}

// LuaKill 停止执行正在运行的lua脚本
func (rc *redisClient) LuaKill() error {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	_, err := rdConn.Do("SCRIPT", "KILL")
	return err
}

// LuaFlush 清空所有已载入的lua脚本
func (rc *redisClient) LuaFlush() error {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	_, err := rdConn.Do("SCRIPT", "FLUSH")
	return err
}

//--------------------------------------------------------------GEO-----------------------------------------------------

// GeoAdd lonLatAddr为组合：经度(longitude)、纬度(latitude)、位置名称(address)。返回添加成功的个数。
func (rc *redisClient) GeoAdd(key string, lonLatAddr ...interface{}) (int, error) {
	if len(lonLatAddr)%3 != 0 {
		return 0, errors.New("位置数据错误！")
	}
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	lonLatAddr = append([]interface{}{key}, lonLatAddr...)
	return redis.Int(rdConn.Do("GEOADD", lonLatAddr...))
}

// GeoGetPosList 获取多个地址的经纬度
func (rc *redisClient) GeoGetPosList(key string, addr ...string) ([]*[2]float64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	addr = append([]string{key}, addr...)
	return redis.Positions(rdConn.Do("GEOPOS", stringToInterface(addr...)...))
}

// GeoGetPos 获取某个地址的经纬度
func (rc *redisClient) GeoGetPos(key string, addr string) (*[2]float64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	pos, err := redis.Positions(rdConn.Do("GEOPOS", key, addr))
	return pos[0], err
}

// GeoGetDistance 计算两个位置之间的距离(米)
func (rc *redisClient) GeoGetDistance(key string, addr1, addr2 string) (float64, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	return redis.Float64(rdConn.Do("GEODIST", key, addr1, addr2, "m"))
}

// GeoGetRadius 通过坐标查找位置集合。经度(longitude)、纬度(latitude)、半径(radius)(米)。
func (rc *redisClient) GeoGetRadius(key string, longitude, latitude, radius float64) ([]*position, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	values, err := redis.Values(rdConn.Do("GEORADIUS", key, longitude, latitude, radius, "m", "WITHCOORD", "WITHDIST", "ASC"))
	if err != nil {
		return nil, err
	}
	return parsePosition(values), err
}

// GeoGetRadiusByAddr 通过地址查找位置集合。经度(longitude)、纬度(latitude)、半径(radius)(米)。
func (rc *redisClient) GeoGetRadiusByAddr(key string, address string, radius float64) ([]*position, error) {
	rdConn := rc.pool.Get()
	// 函数运行结束 ，把连接放回连接池
	defer rdConn.Close()
	values, err := redis.Values(rdConn.Do("georadiusbymember", key, address, radius, "m", "WITHCOORD", "WITHDIST", "ASC"))
	if err != nil {
		return nil, err
	}
	return parsePosition(values), err
}

// 通过地址查找位置集合

type position struct {
	Address   string  // 地址
	Distance  float64 // 距离(米)
	Longitude float64 // 经度
	Latitude  float64 // 纬度
}

func parsePosition(values []interface{}) (positionList []*position) {
	for _, value := range values {
		v := value.([]interface{})
		distance, err := strconv.ParseFloat(string(v[1].([]byte)), 64)
		if err != nil {
			log.Println(err)
		}
		pos := v[2].([]interface{})
		longitude, err := strconv.ParseFloat(string(pos[0].([]byte)), 64)
		if err != nil {
			log.Println(err)
		}
		latitude, err := strconv.ParseFloat(string(pos[1].([]byte)), 64)
		if err != nil {
			log.Println(err)
		}
		positionList = append(positionList, &position{
			Address:   string(v[0].([]byte)),
			Distance:  round(distance, -6),
			Longitude: round(longitude, -6),
			Latitude:  round(latitude, -6),
		})
	}
	return
}

//------------------------------------------------------------工具方法---------------------------------------------------

// 对数字进行四舍五入计算。precision=0时无小数部分，precision<0时为保留的小数位数，precision>0时低precision位为0。
func round(f float64, precision int) float64 {
	if precision == 0 {
		return math.Round(f)
	}
	n := math.Pow10(-precision)
	return math.Round(f*n) / n
}

// 字符串数组转接口数组
func stringToInterface(array ...string) (result []interface{}) {
	for _, s := range array {
		result = append(result, s)
	}
	return
}
