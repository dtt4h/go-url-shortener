package service

import (
	"context"
	"encoding/json"

	"github.com/dtt4h/go-url-shortener/internal/model"
	"github.com/segmentio/kafka-go"
)

type EventService interface {
	PublishURLCreated(ctx context.Context, event model.URLEvent) error
	PublishURLDeleted(ctx context.Context, event model.URLEvent) error
	Close() error
}

type eventService struct {
	writer *kafka.Writer
}

func NewEventService(brokers []string, topic string) EventService {
	writer := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}
	return &eventService{writer: writer}
}
func (s *eventService) PublishURLCreated(ctx context.Context, event model.URLEvent) error {
	event.EventType = "URL_CREATED"
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.ShortCode),
		Value: data,
	}
	return s.writer.WriteMessages(ctx, msg)
}

func (s *eventService) PublishURLDeleted(ctx context.Context, event model.URLEvent) error {
	event.EventType = "URL_DELETED"

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event.ShortCode),
		Value: data,
	}

	return s.writer.WriteMessages(ctx, msg)
}

func (s *eventService) Close() error {
	return s.writer.Close()
}
