package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

func Load(structure interface{}) error {
	return cleanenv.ReadEnv(structure)
}
