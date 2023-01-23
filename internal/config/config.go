package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Token struct {
	PrivateKey string        `yaml:"private_key"`
	PublicKey  string        `yaml:"public_key"`
	ExpiresIn  time.Duration `yaml:"expires_in"`
	MaxAge     int           `yaml:"max_age"`
}

type Config struct {
	Port         string `yaml:"port"`
	MongoURI     string `yaml:"mongo_uri"`
	RedisUri     string `yaml:"redis_uri"`
	AccessToken  Token  `yaml:"access"`
	RefreshToken Token  `yaml:"refresh"`
}

func LoadConfig(path string) (*Config, error) {
	var conf *Config
	yfile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yfile, &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
