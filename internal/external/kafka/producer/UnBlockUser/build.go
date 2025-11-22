package UnBlockUser

import (
	"os"
	"time"
	"users-profile-service/internal/external/kafka/producer"

	"go.uber.org/zap"
)

type Producer struct {
	*producer.BaseProducer
}

type UnblockUserTopicData struct {
	UserId    int64     `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

func Build(logger *zap.SugaredLogger) (*Producer, error) {
	baseProducer, err := producer.New(logger, os.Getenv("UN_BLOCK_USER_TOPIC"))
	if err != nil {
		return nil, err
	}
	return &Producer{baseProducer}, nil
}

func (p *Producer) Produce(data UnblockUserTopicData) error {
	//message := kafka.Message{}
}
