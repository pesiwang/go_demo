package main

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// kafka consumer

func main() {
	topic := "wanglei_test_topic"
	brokers := []string{"alikafka-pre-cn-tk23q8puw001-1-vpc.alikafka.aliyuncs.com:9092", "alikafka-pre-cn-tk23q8puw001-2-vpc.alikafka.aliyuncs.com:9092", "alikafka-pre-cn-tk23q8puw001-3-vpc.alikafka.aliyuncs.com:9092"}
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for {
		for partition := range partitionList { // 遍历所有的分区
			// 针对每个分区创建一个对应的分区消费者
			pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
				return
			}
			defer pc.AsyncClose()
			// 异步从每个分区消费信息
			go func(sarama.PartitionConsumer) {
				for msg := range pc.Messages() {
					fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
				}
			}(pc)
		}

		time.Sleep(time.Second)
	}
}
