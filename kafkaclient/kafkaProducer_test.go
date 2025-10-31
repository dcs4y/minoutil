package kafkaclient

import (
	"testing"
	"time"

	"github.com/IBM/sarama"
)

func TestKafkaProducer(t *testing.T) {
	config := sarama.NewConfig()
	// 发送完数据需要leader和follow都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在success channel返回
	config.Producer.Return.Successes = true

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "GPS_CollisionWarningData"
	msg.Key = sarama.StringEncoder(time.Now().String())
	msg.Value = sarama.StringEncoder("This is my first kafka message!" + time.Now().String())
	msg.Timestamp = time.Now()

	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"192.168.7.226:9092"}, config)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer client.Close()

	// 发送消息
	partition, offset, err := client.SendMessage(msg)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("{partition:%d,offset:%d}", partition, offset)
}

func TestKafkaConsumer(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"192.168.7.226:9092"}, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	topic := "GPS_CollisionWarningData"
	pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		t.Log(err.Error())
	}
	for {
		select {
		case msg := <-pc.Messages():
			t.Logf("{topic:%s,partition:%d,offset:%d,key:%v,value:%v}", msg.Topic, msg.Partition, msg.Offset, msg.Key, string(msg.Value))
		}
	}
	time.Sleep(time.Second * 10)
}
