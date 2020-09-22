package mq

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

type KafkaConsumer struct {
	*Kafka
	wg sync.WaitGroup
}

func (*KafkaConsumer) config() *sarama.Config {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_10_2_0
	return config
}

func (c *KafkaConsumer) Init() {
	log.Println("[kafka-consumer-init]")
}

func (c *KafkaConsumer) Consume(topic string, handler HandlerFunc) error {
	consumer, err := sarama.NewConsumer(c.Addrs, c.config())
	if err != nil {
		log.Println("[kafka-consumer-error]", err.Error())
		panic(err.Error())
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Println("[kafka-consumer-error]", err.Error())
		return err
	}

	for _, partition := range partitions {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			log.Println("[kafka-consumer-error]", err.Error())
			continue
		}

		c.wg.Add(1)
		go c.message(pc, handler)
	}

	c.wg.Wait()

	return nil
}

func (c *KafkaConsumer) message(pc sarama.PartitionConsumer, handler HandlerFunc) {
	defer c.wg.Done()
	defer pc.Close()

	for {
		select {
		case msg := <-pc.Messages():
			handler(msg.Value)
			log.Println("[kafka-consume-success]",
				fmt.Sprintf("partitions=%d topic=%s offset=%d key=%s value=%s",
					msg.Partition, msg.Topic, msg.Offset, string(msg.Key), string(msg.Value)))

		case err := <-pc.Errors():
			log.Println("[kafka-consumer-error]", err.Error())

		case <-c.Context.Done():
			return
		}
	}
}
