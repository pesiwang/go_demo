package main

import (
	"flag"
	"fmt"
	"time"
	"util"

	"github.com/nsqio/go-nsq"
)

var lookupAddr string

func init() {
	flag.StringVar(&lookupAddr, "lookupAddr", "127.0.0.1:4161", "nsqlookupd addr")

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
		fmt.Println(err)
	}
	// 发布消息
	for {
		if err := producer.Publish("test2", []byte("test message")); err != nil {
			fmt.Println("publish error: " + err.Error())
		}
		fmt.Println("send test")
		time.Sleep(1 * time.Second)
	}
}

func main() {

	startProducer()
}
