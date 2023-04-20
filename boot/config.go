package boot

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v3"
)

var (
	Clients = []string{"gcs", "dos", "s3"}
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
	Consumer                    struct {
		Network       string `yaml:"network"`
		BrokerAddress string `yaml:"brokerAddress"`
		Topic         string `yaml:"topic"`
		GroupID       string `yaml:"groupID"`
		Partition     int    `yaml:"partition"`
	} `yaml:"consumer"`
}

type Sess struct {
	SessionMap map[string]*session.Session
	GCSSession *session.Session
	DOSSession *session.Session
	S3Session  *session.Session
}

type Client struct {
	ClientMap map[string]*s3.S3
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
