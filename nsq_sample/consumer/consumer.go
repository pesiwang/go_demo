package main

import (
	"flag"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

var lookupAddr string

func init() {
	flag.StringVar(&lookupAddr, "lookupAddr", "127.0.0.1:4161", "nsqlookupd addr")

	flag.Parse()
}

// 消费者
func startConsumer() {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 40
	cfg.LookupdPollInterval = time.Second * 10
	consumer, err := nsq.NewConsumer("test1", "sensor01", cfg)
	if err != nil {
		log.Fatal(err)
	}
	// 设置消息处理函数
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Println(string(message.Body))
		return nil
	}))
	// nsqlookupd
	//[]string
	if err := consumer.ConnectToNSQLookupds([]string{lookupAddr}); err != nil {
		log.Fatal(err)
	}
	<-consumer.StopChan
}

func main() {
	startConsumer()
}
