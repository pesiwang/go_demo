package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MongoDBHosts = "mongodb://root:hBsYnVjga6lg@dds-wz984ef1e5307f541538-pub.mongodb.rds.aliyuncs.com:3717,dds-wz984ef1e5307f542851-pub.mongodb.rds.aliyuncs.com:3717/admin?replicaSet=mgset-76035551"
	MaxCon       = 10
)

type UserMailbox struct {
	MessageID uint64 `json:"msgid" bson:"msgid"`
	Mid       int64  `json:"mid" bson:"mid"` // CAUTION not the same
	SessionID uint64 `json:"session_id" bson:"session_id"`
	NewCt     int64  `json:"new_ct,omitempty" bson:"new_ct,omitempty"`
}

func JsonStr(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// 实现 mongo.Logger 的接口
type MongoLog struct {
}

func (MongoLog) Output(calldepth int, s string) error {
	_, err := fmt.Printf("mgo log: %v\n", s)
	return err
}

func main() {

	session, err := mgo.DialWithTimeout(MongoDBHosts, 30*time.Second)
	if err != nil {
		log.Fatalf("CreateSession failed:%n", err)
	}
	mgo.SetDebug(true)           // 设置DEBUG模式
	mgo.SetLogger(new(MongoLog)) // 设置日志.

	//设置连接池的大小
	session.SetPoolLimit(MaxCon)
	t := time.Now() // time.Date(2024, time.August, 29, 17, 45, 9, 354194500, time.Local)
	query := bson.M{
		"new_ct": bson.M{"$lte": t.UnixMilli()},
	}

	// query := bson.M{
	// 	"mid": 17008610,
	// }
	mailboxList := make([]*UserMailbox, 0, 3)
	err = session.DB("chat").C("mailbox").Find(query).Skip(0).Limit(3).All(&mailboxList)
	if err != nil {
		fmt.Printf("err %v", err)
	} else {
		for i, v := range mailboxList {
			fmt.Printf("return [%v]=%v\n", i, JsonStr(*v))
		}
	}

	defer session.Close()
}
