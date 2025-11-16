package kafka

import (
	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/kafka"
	"github.com/nnieie/golanglab5/pkg/logger"
)

func InitKafka() {
	brokers := config.Kafka.Brokers
	topic := "interaction-events"

	if config.Kafka != nil {
		brokers = config.Kafka.Brokers
		topic = config.Kafka.Topic
	} else {
		logger.Warnf("Kafka config not found, using default values")
	}

	err := kafka.InitProducer(kafka.ProducerConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	if err != nil {
		logger.Errorf("Failed to initialize Kafka producer: %v", err)
	}
}
