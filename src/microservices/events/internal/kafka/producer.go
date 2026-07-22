package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *Producer) SendMessage(key string, value interface{}) (partition int32, offset int64, err error) {
	msgBytes, err := json.Marshal(value)
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(msgBytes),
	}

	partition, offset, err = p.producer.SendMessage(msg)
	if err != nil {
		return
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", p.topic, partition, offset)
	return
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
