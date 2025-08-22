package kafkaservice

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaServiceStruct struct {
	ctx       context.Context
	reader    *kafka.Reader
	listeners []Listener
}

// Listen implements KafkaService.
func (k *KafkaServiceStruct) Listen() error {
	log.Println("Staring listening")
	for {
		m, err := k.reader.ReadMessage(k.ctx)

		if err != nil {
			return err
		}

		log.Printf("Received new message, calling %v listeners", len(k.listeners))
		for _, listener := range k.listeners {
			listener(m)
		}
	}
}

// RegisterListener implements KafkaService.
func (k *KafkaServiceStruct) RegisterListener(listener Listener) error {
	k.listeners = append(k.listeners, listener)
	return nil
}

type Listener func(msg kafka.Message)

type KafkaService interface {
	RegisterListener(listener Listener) error
	Listen() error
}

func ConnectToKafks(brokers []string, topic string, groupId string, ctx context.Context) (KafkaService, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupId,
	})
	service := KafkaServiceStruct{ctx, reader, make([]Listener, 0, 2)}
	return &service, nil
}
