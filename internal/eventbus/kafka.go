package eventbus

import (
	"context"
	"encoding/json"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

type Bus struct {
	writer *kafka.Writer
}

func NewBus(brokers []string) *Bus {
	return &Bus{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			RequiredAcks: kafka.RequireAll,
			Balancer:     &kafka.LeastBytes{},
		},
	}
}

func (b *Bus) Publish(ctx context.Context, topic string, evt Event) error {
	data, _ := json.Marshal(evt)
	return b.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(evt.Type),
		Value: data,
	})
}

func Consume(ctx context.Context, brokers []string, topic, group string, fn func(Event) error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: group,
		Topic:   topic,
	})

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Printf("kafka fetch: %v", err)
			continue
		}
		var evt Event
		if err := json.Unmarshal(m.Value, &evt); err == nil {
			if err := fn(evt); err == nil {
				_ = r.CommitMessages(ctx, m)
			}
		}
	}
}
