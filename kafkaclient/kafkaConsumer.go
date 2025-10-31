package kafkaclient

import (
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaHandler interface {
	Handler(pc sarama.PartitionConsumer)
}

type kafkaConsumer struct {
	consumer        sarama.Consumer
	kafkaHandlerMap map[string]KafkaHandler
}

// NewKafkaConsumer 创建消息者
func NewKafkaConsumer(servers []string) kafkaConsumer {
	consumer, err := sarama.NewConsumer(servers, nil)
	if err != nil {
		fmt.Println("KafkaConsumer初始化失败：", err.Error())
	}
	return kafkaConsumer{consumer: consumer, kafkaHandlerMap: make(map[string]KafkaHandler)}
}

func (kafkaConsumer kafkaConsumer) AddHandler(topic string, kafkaHandler KafkaHandler) {
	kafkaConsumer.kafkaHandlerMap[topic] = kafkaHandler
}

func (kafkaConsumer kafkaConsumer) Run() {
	for topic, handler := range kafkaConsumer.kafkaHandlerMap {
		partitionList, err := kafkaConsumer.consumer.Partitions(topic)
		if err != nil {
			panic(err.Error())
		}
		for partition := range partitionList {
			pc, err := kafkaConsumer.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				panic(err.Error())
			}
			go handler.Handler(pc)
		}
	}
}
