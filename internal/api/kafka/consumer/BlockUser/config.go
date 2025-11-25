package BlockUser

import (
	"sync"
	"users-profile-service/internal/config"
)

type ConsumerConfig struct {
	Brokers []string `env:"KAFKA_BROKERS" env-separate:","`
	GroupId string   `env:"KAFKA_CONSUMER_GROUP"`
	Topic   string   `env:"BLOCK_UNBLOCK_USER_TOPIC"`
}

var (
	once        sync.Once
	cgfInstance *ConsumerConfig
)

func NewConsumerConfig() (*ConsumerConfig, error) {
	var err error
	once.Do(func() {
		cfgInstance := &ConsumerConfig{}
		err = config.Load(cfgInstance)
	})
	return cgfInstance, err
}
