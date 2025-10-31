package kafkaclient

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

var client sarama.SyncProducer

// NewClient 创建生产者客户端
func NewClient(servers []string) {
	config := sarama.NewConfig()
	// 发送完数据需要leader和follow都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 3
	// 连接kafka
	producer, err := sarama.NewSyncProducer(servers, config)
	if err != nil {
		fmt.Println("KafkaProducer初始化失败：", err.Error())
	}
	client = producer
}

func SendKafkaMessage(topic string, key string, value string) {
	// 构造一个消息
	message := &sarama.ProducerMessage{}
	message.Topic = topic
	message.Key = sarama.StringEncoder(key)
	message.Value = sarama.StringEncoder(value)
	message.Timestamp = time.Now()
	// 发送消息
	partition, offset, err := client.SendMessage(message)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("已发送消息到Kafka：{partition:%d,offset:%d}\n", partition, offset)
}
