package DeletedUser

import (
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
	baseProducer, err := producer.New(logger, os.Getenv("DELETED_USER_TOPIC"))
	if err != nil {
		return nil, err
	}
	return &Producer{baseProducer}, nil
}

func (p *Producer) Produce(id int64) error {
	message := kafka.Message{Value: []byte(fmt.Sprintf("id:%d", id))}
	return p.BaseProducer.Send(&message)
}
