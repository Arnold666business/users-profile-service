package producer

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ProducerConfig struct {
	Brokers          []string `env:"KAFKA_BROKERS" env-separator:","`
	NewUserTopic     string   `env:"NEW_USER_TOPIC"`
	DeletedUserTopic string   `env:"DELETED_USER_TOPIC"`
	RetryAttempts    int      `env:"KAFKA_RETRY_ATTEMPTS"`
	RetryInterval    int      `env:"KAFKA_RETRY_INTERVAL_SECONDS"`
	Timeout          int      `env:"KAFKA_SEND_TIMEOUT_SECONDS"`
}

type BaseProducer struct {
	writer *kafka.Writer
	cfg    *ProducerConfig
	logger *zap.SugaredLogger
	topic  string
}

func New(logger *zap.SugaredLogger, topic string) (*BaseProducer, error) {
	l := logger.Named("kafka.producer").With(zap.String("topic", topic))
	cfg, err := NewProducerConfig()
	if err != nil {
		return nil, err
	}
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		BatchSize:              100,
		BatchTimeout:           50 * time.Millisecond,
		BatchBytes:             1048576,
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: false,
		Async:                  false,
	}
	return &BaseProducer{
		writer: writer,
		cfg:    cfg,
		logger: l,
		topic:  topic,
	}, nil
}

func (p *BaseProducer) Send(ctx context.Context, msg *kafka.Message) error {
	err := ProduceWithRetry(p.logger, p.cfg, func(ctx context.Context) error {
		return p.writer.WriteMessages(ctx, *msg)
	})

	if err != nil {
		p.logger.Warnf("error writing to kafka topic %s: %s", p.topic, err)
		return err
	}
	return nil
}

func (p *BaseProducer) Close() error {
	return p.writer.Close()
}
