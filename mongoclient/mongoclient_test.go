package mongoclient

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func Test_GetCollection(t *testing.T) {
	collectionName := "DeviceMessage"
	c := GetClient().GetCollection(collectionName)
	t.Log("collectionName=", c.GetCollectionName())
	m := bson.M{}
	m["name"] = "d2"
	m["a"] = "a"
	m["b"] = "b"
	m["c"] = "c"
	m["d"] = 0
	ctx := context.Background()
	result, err := c.InsertOne(ctx, m) // 保存后会自动创建库和集合
	if err != nil {
		fmt.Println(err)
	}
	t.Log("id=", result)
	// 查询
	filter := make(map[string]interface{})
	filter["name"] = "d"
	q := c.Find(ctx, filter)
	var list []map[string]interface{}
	err = q.All(&list)
	if err != nil {
		fmt.Println(err)
	}
	t.Log("find:", list)
}
