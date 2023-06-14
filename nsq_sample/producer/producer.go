package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"util"

	"github.com/nsqio/go-nsq"
)

var lookupAddr string

func init() {
	flag.StringVar(&lookupAddr, "lookupAddr", "101.132.227.177:4161", "nsqlookupd addr")

	flag.Parse()
}

// 生产者
func startProducer() {
	nsqdHelper := util.NewNsqdHelper()

	nodeIps, err := nsqdHelper.GetNodesIps([]string{lookupAddr})
	if err != nil {
		fmt.Printf("get nqsd node failed:%v\n", err)
		return
	}

	fmt.Printf("nodeIps:%v\n", nodeIps)
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(nodeIps[0], cfg)
	if err != nil {
		log.Fatal(err)
	}
	// 发布消息
	for {
		if err := producer.Publish("test", []byte("test message")); err != nil {
			log.Fatal("publish error: " + err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {

	startProducer()
}
