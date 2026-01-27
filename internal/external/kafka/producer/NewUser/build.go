package NewUser

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

type NewUserTopicData struct {
	Id   int64
	Role int
}

func Build(logger *zap.SugaredLogger) (*Producer, error) {
	baseProducer, err := producer.New(logger, os.Getenv("NEW_USER_TOPIC"))
	if err != nil {
		return nil, err
	}
	return &Producer{baseProducer}, nil
}

func (p *Producer) Produce(data NewUserTopicData) error {
	message := kafka.Message{Value: []byte(fmt.Sprintf("{id:%d, role:%d}", data.Id, data.Role))}
	return p.BaseProducer.Send(&message)
}
