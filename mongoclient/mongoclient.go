package mongoclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// https://pkg.go.dev/go.mongodb.org/mongo-driver/v2@v2.9.2#section-readme
// https://www.mongodb.com/zh-cn/docs/drivers/go/current/

var clients = make(map[string]*mongoClient)

type MongoConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type mongoClient struct {
	database *mongo.Database
}

func GetClientByName(name string) *mongoClient {
	return clients[name]
}

func GetClient() *mongoClient {
	return clients[""]
}

func NewClient(name string, config MongoConfig) *mongoClient {
	mgClient, err := mongo.Connect(options.Client().SetTimeout(10 * time.Second).ApplyURI(config.Uri).SetAuth(options.Credential{Username: config.Username, Password: config.Password}))
	if err != nil {
		fmt.Println("NewMongoClient:", err)
	}
	client := &mongoClient{
		mgClient.Database(config.Database),
	}
	clients[name] = client
	return client
}

func (mc *mongoClient) GetCollection(collectionName string) *mongo.Collection {
	return mc.database.Collection(collectionName)
}

func (mc *mongoClient) SaveM(collection string, m bson.M) any {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := mc.GetCollection(collection).InsertOne(ctx, m)
	if err != nil {
		log.Println(err.Error())
	}
	return result.InsertedID
}

func (mc *mongoClient) SaveD(collection string, d bson.D) any {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := mc.GetCollection(collection).InsertOne(ctx, d)
	if err != nil {
		log.Println(err.Error())
	}
	return result.InsertedID
}

func (mc *mongoClient) Disconnect() {
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := mc.database.Client().Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func (mc *mongoClient) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = mc.database.Client().Ping(ctx, readpref.Primary())
}
