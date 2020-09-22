package mq

var __mq iMQ

func Init(mq iMQ) {
	__mq = mq
	__mq.Init()
}

func SendMessage(topic string, message []byte) error {
	return __mq.Producer().SendMessage(topic, message)
}

func Consume(topic string, handler HandlerFunc) error {
	return __mq.Consumer().Consume(topic, handler)
}

func ConsumeGroup(groupId string, topics []string, handler HandlerFunc) error {
	return __mq.ConsumerGroup(groupId).Consume(topics, handler)
}
