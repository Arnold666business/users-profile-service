package NewUser

import (
	"context"
	"fmt"
	"os"
	"users-profile-service/internal/external/kafka/producer"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	*producer.BaseProducer
}

func Build(logger *zap.SugaredLogger) (*Producer, error) {
	baseProducer, err := producer.New(logger, os.Getenv("NEW_USER_TOPIC"))
	if err != nil {
		return nil, err
	}
	return &Producer{baseProducer}, nil
}

func (p *Producer) produce(id int64, role int) error {
	message := kafka.Message{Value: []byte(fmt.Sprintf("{id:%d, role:%d}", id, role))}
	return p.BaseProducer.Send(context.Background(), &message)
}
