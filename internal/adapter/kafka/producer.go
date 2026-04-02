package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/Coiiap5e/photographer/internal/errors"
	"github.com/Coiiap5e/photographer/internal/model"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer                   *kafka.Writer
	logger                   *slog.Logger
	meter                    metric.Meter
	kafkaMessagesSentCounter metric.Int64Counter
}

func NewKafkaProducer(brokerURLs []string, topic string, logger *slog.Logger, meter metric.Meter) (*KafkaProducer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokerURLs,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: 1,
		Logger:       kafka.LoggerFunc(logger.Debug),
		ErrorLogger:  kafka.LoggerFunc(logger.Error),
	})

	producer := &KafkaProducer{
		writer: writer,
		logger: logger,
		meter:  meter,
	}

	kafkaMessagesSentCounter, err := meter.Int64Counter(
		"kafka_messages_sent_total",
		metric.WithDescription("Total number of messages sent to Kafka"),
		metric.WithUnit("1"),
	)
	if err != nil {
		logger.Error("failed to create Kafka messages sent counter", "error", err)
		return nil, errors.Wrap(err, errors.ErrCodeKafkaProduce, "failed to create Kafka messages sent counter")
	}
	producer.kafkaMessagesSentCounter = kafkaMessagesSentCounter

	return producer, nil
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

	p.kafkaMessagesSentCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("topic", p.writer.Topic), attribute.String("message_type", "shoot_notification")))
	p.logger.Info("shoot notification successfully sent to Kafka", "shoot_id", shoot.Id, "topic", p.writer.Topic, "metric_count", 1)
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

	p.kafkaMessagesSentCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("topic", p.writer.Topic), attribute.String("message_type", "general_message")))
	p.logger.Info("text notification successfully sent to Kafka", "message", message, "topic", p.writer.Topic, "metric_count", 1)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
