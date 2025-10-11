package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var once sync.Once

func Load(structure interface{}) error {
	var err error
	once.Do(func() {
		err = cleanenv.ReadEnv(structure)
	})
	return err
}
