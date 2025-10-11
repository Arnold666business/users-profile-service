package Audit

import (
	"context"
	"os"
	"users-profile-service/internal/external/kafka/producer"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	*producer.BaseProducer
}

func Build(logger *zap.SugaredLogger) (*Producer, error) {
	baseProducer, err := producer.New(logger, os.Getenv("AUDIT_TOPIC"))
	if err != nil {
		return nil, err
	}
	return &Producer{baseProducer}, nil
}

func (p *Producer) Produce(ctx context.Context, id int64) error {
	message := kafka.Message{} ????
	return p.BaseProducer.Send(ctx, &message)
}
