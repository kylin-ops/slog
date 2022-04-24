package kafka

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
)

func NewProduce(addrs []string, topics ...string) (*Produce, error) {
	w := &kafka.Writer{
		Addr:         kafka.TCP(addrs...),
		Async:        true,
		Balancer:     &kafka.Hash{},
		RequiredAcks: 0,
	}

	if len(topics) == 0 {
		return nil, errors.New("没有指定topic")
	}

	return &Produce{
		writer: w,
		topic:  topics[0],
		topics: topics,
	}, nil
}

type Produce struct {
	writer *kafka.Writer
	topic  string
	topics []string
}

func (p *Produce) SendSingleTopicMessage(key string, value []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: p.topic,
		Key:   []byte(key),
		Value: value,
	})
}

func (p *Produce) WriteMultipleTopicsMessage(key string, value []byte) error {
	var messages []kafka.Message
	for _, topic := range p.topics {
		messages = append(messages, kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: []byte(value),
		})
	}
	return p.writer.WriteMessages(context.Background(), messages...)
}
