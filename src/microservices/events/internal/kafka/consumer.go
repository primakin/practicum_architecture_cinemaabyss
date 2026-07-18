package kafka

import (
	"context"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (c *Consumer) Start() error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		c.wg.Add(1)
		go c.consumePartition(pc)
	}

	log.Printf("Consumer started for topic: %s", c.topic)
	return nil
}

func (c *Consumer) consumePartition(pc sarama.PartitionConsumer) {
	defer c.wg.Done()
	defer pc.Close()

	for {
		select {
		case msg := <-pc.Messages():
			c.handleMessage(msg)
		case err := <-pc.Errors():
			log.Printf("Consumer error: %v", err)
		case <-c.ctx.Done():
			log.Println("Consumer stopped")
			return
		}
	}
}

func (c *Consumer) handleMessage(msg *sarama.ConsumerMessage) {
	log.Printf("Message to handle: %s", msg.Value)
}

func (c *Consumer) Stop() {
	c.cancel()
	c.wg.Wait()
	c.consumer.Close()
	log.Println("Consumer stopped")
}
