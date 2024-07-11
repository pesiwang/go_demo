package main

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

func SyncProduce() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	msg := &sarama.ProducerMessage{}
	msg.Topic = "wanglei_test_topic"
	content := "this is a sync message"
	msg.Value = sarama.StringEncoder(content)

	client, err := sarama.NewSyncProducer([]string{"alikafka-pre-cn-tk23q8puw001-1-vpc.alikafka.aliyuncs.com:9092", "alikafka-pre-cn-tk23q8puw001-2-vpc.alikafka.aliyuncs.com:9092", "alikafka-pre-cn-tk23q8puw001-3-vpc.alikafka.aliyuncs.com:9092"}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.Close()

	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}

	fmt.Printf("send msg partition/offset: %d/%d, value is: %s", partition, offset, content)

}

func AsyncProduce() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewAsyncProducer([]string{"alikafka-pre-cn-tk23q8puw001-1-vpc.alikafka.aliyuncs.com:9092", "alikafka-pre-cn-tk23q8puw001-2-vpc.alikafka.aliyuncs.com:9092", "alikafka-pre-cn-tk23q8puw001-3-vpc.alikafka.aliyuncs.com:9092"}, config)
	if err != nil {
		fmt.Println("error is:", err.Error())
		return
	}
	defer client.AsyncClose()

	go func(p sarama.AsyncProducer) {
		for {
			select {
			case msg := <-p.Successes():
				value, _ := msg.Value.Encode()
				fmt.Printf("send msg partition/offset: %d/%d, value is: %s", msg.Partition, msg.Offset, string(value))
				return
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
				return
			}
		}
	}(client)

	msg := &sarama.ProducerMessage{
		Topic: "wanglei_test_topic",
		Value: sarama.ByteEncoder("this is a async message"),
	}
	client.Input() <- msg
	time.Sleep(time.Second * 1)
}

func main() {
	AsyncProduce()
	// fmt.Println()
	// SyncProduce()
	time.Sleep(time.Second * 2)
}
