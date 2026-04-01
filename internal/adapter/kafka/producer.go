package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

func NewKafkaProducer(brokerURLs []string, topic string, logger *slog.Logger) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokerURLs,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: 1,
		Logger:       kafka.LoggerFunc(logger.Debug),
		ErrorLogger:  kafka.LoggerFunc(logger.Error),
	})

	return &KafkaProducer{
		writer: writer,
		logger: logger,
	}
}

func (p *KafkaProducer) Notify(shoot model.Shoot) error {
	shootJSON, err := json.Marshal(shoot)
	if err != nil {
		p.logger.Error("error serializing shoot to JSON", "error", err)
		return errors.Wrap(err, errors.ErrCodeKafkaProduce, "failed to serialize shoot data for Kafka")
	}

	msg := kafka.Message{
		Key:   []byte("shoot_notification"),
		Value: shootJSON,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error("error sending message to Kafka", "error", err, "topic", p.writer.Topic, "key", string(msg.Key))
		return errors.Wrap(err, errors.ErrCodeKafkaProduce, "failed to send shoot notification to Kafka")
	}

	p.logger.Info("shoot notification successfully sent to Kafka", "shoot_id", shoot.Id, "topic", p.writer.Topic)
	return nil
}

func (p *KafkaProducer) NotifyMessage(message string) error {
	msg := kafka.Message{
		Key:   []byte("general_message"),
		Value: []byte(message),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Error("error sending message to Kafka", "error", err, "topic", p.writer.Topic, "key", string(msg.Key))
		return errors.Wrap(err, errors.ErrCodeKafkaProduce, "failed to send text notification to Kafka")
	}

	p.logger.Info("text notification successfully sent to Kafka", "message", message, "topic", p.writer.Topic)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
