package BlockUser

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BlockedInfo struct {
	Id          int64     `json:"id"`
	UnblockData time.Time `json:"unblock_data"`
	Forever     bool      `json:"forever_flag"`
}

func Build(logger *zap.SugaredLogger) error {
	cfg, err := NewConsumerConfig()
	if err != nil {
		return err
	}
	l := logger.Named("kafka.consumer").With(
		zap.String("topic", cfg.Topic),
		zap.String("group_id", cfg.GroupId))

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupId,
	})

	for {
		var msg kafka.Message
		msg, err = reader.ReadMessage(context.Background())
		if err != nil {
			l.Errorf("Error reading message from kafka: %v", err)
			break
		}
		var data BlockedInfo
		err = json.Unmarshal(msg.Value, &data)
		if err != nil {
			l.Errorf("Error unmarshalling message from kafka: %v", err)
			continue
		}
	}
	reader.Close()
	return nil
}
