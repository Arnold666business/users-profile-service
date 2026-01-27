package BlockUser

import (
	"time"
	"users-profile-service/internal/user/block"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BlockUser struct {
	conn      *kafka.Reader
	logger    *zap.SugaredLogger
	processor *block.BlockUserProcessor
}

type BlockedInfo struct {
	UserID      int64     `json:"id"`
	BlockTypeID int       `json:"block_type_id"`
	Forever     bool      `json:"forever_flag"`
	BlockAt     time.Time `json:"block_at"`
}

func Build(logger *zap.SugaredLogger, processor *block.BlockUserProcessor) (*BlockUser, error) {
	cfg, err := NewConsumerConfig()
	if err != nil {
		return nil, err
	}
	l := logger.Named("kafka.consumer").With(
		zap.String("topic", cfg.Topic),
		zap.String("group_id", cfg.GroupId))

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupId,
		Topic:   cfg.Topic,
	})

	return &BlockUser{
		conn:      reader,
		logger:    l,
		processor: processor,
	}, nil
}
