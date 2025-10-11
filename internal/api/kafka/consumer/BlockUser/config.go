package BlockUser

import "users-profile-service/internal/config"

type ConsumerConfig struct {
	Brokers []string `env:"KAFKA_BROKERS" env-separate:","`
	GroupId string   `env:"KAFKA_CONSUMER_GROUP"`
	Topic   string   `env:"BLOCK_UNBLOCK_USER_TOPIC"`
}

func NewConsumerConfig() (*ConsumerConfig, error) {
	cfg := &ConsumerConfig{}
	err := config.Load(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
