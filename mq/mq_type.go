package mq

type HandlerFunc func(data []byte) bool

type iMQ interface {
	Init()
	Producer() iProducer
	Consumer() iConsumer
	ConsumerGroup(groupId string) iConsumerGroup
}

type iProducer interface {
	SendMessage(topic string, message []byte) error
}

type iConsumer interface {
	Consume(topic string, handler HandlerFunc) error
}

type iConsumerGroup interface {
	Consume(topics []string, handler HandlerFunc) error
}
