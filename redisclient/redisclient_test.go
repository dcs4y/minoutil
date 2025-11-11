package redisclient

import (
	"fmt"
	"testing"
)

func init() {
	NewClient("", RedisConfig{
		Host:     "192.168.4.99",
		Port:     6379,
		Database: 0,
		Password: "123456",
	})
}

func TestRedis(t *testing.T) {
	rc := GetClient()
	key := "zsettest"
	t.Log(rc.ZSetAdd(key, 1, "11", 20, "0a0", 1.5, "test", 10, "很多"))
	t.Log(rc.ZSetLength(key))
	t.Log(rc.ZSetGetScore(key, "test"))
	t.Log(rc.ZSetGetByScore(key, 0, 20, false))
	t.Log(rc.ZSetDelete(key, "test", "0a0"))
	t.Log(rc.ZSetDeleteByScore(key, 10, 20))
}

func TestClient_ChannelSubscribe(t *testing.T) {
	rc := GetClient()
	ch := rc.ChannelSubscribe("chan1")
	for {
		fmt.Println("等等通道消息...")
		rm, ok := <-ch
		if ok {
			fmt.Println(rm)
		}
	}
}

func TestClient_ChannelSubscribeEvent(t *testing.T) {
	rc := GetClient()
	// 订阅88数据库所有事件
	ch := rc.ChannelSubscribe("__key*@88__:*")
	for {
		fmt.Println("等等通道消息...")
		rm, ok := <-ch
		if ok {
			fmt.Println(rm)
		}
	}
}

func TestClient_ChannelPublish(t *testing.T) {
	rc := GetClient()
	fmt.Println(rc.ChannelPublish("chan1", "消息1"))
	fmt.Println(rc.ChannelPublish("chan1", "消息2"))
	fmt.Println(rc.ChannelPublish("chan1", "消息3"))
}

func TestClient_LuaScript(t *testing.T) {
	rc := GetClient()
	// 内置方法：redis.call("cmd","key","param")
	fmt.Println(rc.LuaExecute("return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}", []string{"key1", "key2"}, []string{"first", "second"}))
	fmt.Println(rc.LuaLoad("return {'key',22}"))
	fmt.Println(rc.LuaExists("e0e1f9fabfc9d4800c877a703b823ac0578ff8db"))
	fmt.Println(rc.LuaExecuteWithSha1("e0e1f9fabfc9d4800c877a703b823ac0578ff8db", []string{}, []string{}))
	//fmt.Println(rc.LuaKill())
	//fmt.Println(rc.LuaFlush())
}

func TestClient_GeoAdd(t *testing.T) {
	rc := GetClient()
	key := "geotest"
	fmt.Println(rc.GeoAdd(key, 107.123, 38.987653, "地址1", 107.98797, 38.2234, "地址2", 100, 45, "又是一个地址"))
	posList, err := rc.GeoGetPosList(key, "地址1")
	if err != nil {
		fmt.Println(err)
	}
	for _, pos := range posList {
		if pos == nil {
			break
		}
		fmt.Println(pos[0], pos[1])
	}
	fmt.Println(rc.GeoGetPos(key, "地址1"))
	fmt.Println(rc.GeoGetDistance(key, "地址1", "又是一个地址"))
	{
		positionList, err := rc.GeoGetRadius(key, 108, 39, 200000)
		if err != nil {
			fmt.Println(err)
		}
		for _, pos := range positionList {
			fmt.Println(*pos)
		}
	}
	{
		positionList, err := rc.GeoGetRadiusByAddr(key, "地址3", 10)
		if err != nil {
			fmt.Println(err)
		}
		for _, pos := range positionList {
			fmt.Println(*pos)
		}
	}
}
