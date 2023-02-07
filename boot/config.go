package boot

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	GoogleAccessKeyID           string `yaml:"googleAccessKeyID"`
	GoogleAccessKeySecret       string `yaml:"googleAccessKeySecret"`
	DigitalOceanAccessKeyID     string `yaml:"digitalOceanAccessKeyID"`
	DigitalOceanAccessKeySecret string `yaml:"digitalOceanAccessKeySecret"`
	AmazonAccessKeyID           string `yaml:"amazonAccessKeyID"`
	AmazonAccessKeySecret       string `yaml:"amazonAccessKeySecret"`
}

func LoadConfig(file string) (Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return Config{}, fmt.Errorf("error loading config file: %w", err)
	}

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		err = fmt.Errorf("error decoding config file: %w", err)
	}
	return cfg, err
}
