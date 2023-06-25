package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func initEngine() {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://root:1jst7QTWIHCgIypT@dds-uf64ea78c8dea6241704-pub.mongodb.rds.aliyuncs.com:3717")

	// 连接到MongoDB
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
func GetMgoCli() *mongo.Client {
	if mongoClient == nil {
		initEngine()
	}
	return mongoClient
}

func main() {
	var (
		client     = GetMgoCli()
		db         *mongo.Database
		collection *mongo.Collection
	)
	//2.选择数据库 my_db
	db = client.Database("project_id")

	//选择表 my_collection
	collection = db.Collection("client")

	query := bson.M{"_id": "clientID_xxxx"}
	var result bson.M
	err := collection.FindOne(context.TODO(), query).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments || err.Error() == mongo.ErrNoDocuments.Error() {
			fmt.Printf("no document, error:%v", err)
		} else {
			fmt.Printf("unknown error:%v", err)
		}
	}

	fmt.Println(result)
}
