package kafka

import "github.com/nnieie/golanglab5/pkg/kafka"

var KafkaInstance *kafka.Kafka

func InitKafka() {
	instance := kafka.NewKafkaInstance()
	KafkaInstance = instance
}

func CloseKafka() {
	KafkaInstance.Close()
}
