package mongoclient

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

// https://www.cnblogs.com/qiniu/p/13492504.html

var clients = make(map[string]*mongoClient)

type MongoConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type mongoClient struct {
	database *qmgo.Database
}

func GetClientByName(name string) *mongoClient {
	return clients[name]
}

func GetClient() *mongoClient {
	return clients[""]
}

func NewClient(name string, config MongoConfig) *mongoClient {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mgClient, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: config.Uri, Auth: &qmgo.Credential{Username: config.Username, Password: config.Password}})
	if err != nil {
		fmt.Println("NewMongoClient:", err)
	}
	client := &mongoClient{
		mgClient.Database(config.Database),
	}
	clients[name] = client
	return client
}

func (mc *mongoClient) GetCollection(collectionName string) *qmgo.Collection {
	return mc.database.Collection(collectionName)
}

func (mc *mongoClient) Save(collection string, m bson.M) string {
	result, err := mc.GetCollection(collection).InsertOne(context.TODO(), m)
	if err != nil {
		log.Println(err.Error())
	}
	id := result.InsertedID
	return id.(primitive.ObjectID).Hex()
}
