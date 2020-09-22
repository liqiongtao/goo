package mq

import (
	"context"
	"log"
	"sync"
)

type Kafka struct {
	context.Context
	Addrs    []string
	producer *KafkaProducer
	mu       sync.Mutex
}

func (k *Kafka) Init() {
	log.Println("[kafka-init]")
}

func (k *Kafka) Producer() iProducer {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.producer == nil {
		k.producer = &KafkaProducer{Kafka: k}
		k.producer.Init()
	}

	return k.producer
}

func (k *Kafka) Consumer() iConsumer {
	return &KafkaConsumer{Kafka: k}
}

func (k *Kafka) ConsumerGroup(groupId string) iConsumerGroup {
	return &KafkaConsumerGroup{Kafka: k, GroupId: groupId}
}
