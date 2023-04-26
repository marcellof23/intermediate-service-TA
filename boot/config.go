package boot

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v3"
)

var (
	Clients []string
	Backup  = "backup"
)

type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	GoogleAccessKeyID           string   `yaml:"googleAccessKeyID"`
	GoogleAccessKeySecret       string   `yaml:"googleAccessKeySecret"`
	GCSRegion                   string   `yaml:"gcsRegion"`
	GCSEndpoint                 string   `yaml:"gcsEndpoint"`
	DigitalOceanAccessKeyID     string   `yaml:"digitalOceanAccessKeyID"`
	DigitalOceanAccessKeySecret string   `yaml:"digitalOceanAccessKeySecret"`
	DOSRegion                   string   `yaml:"dosRegion"`
	DOSEndpoint                 string   `yaml:"dosEndpoint"`
	AmazonAccessKeyID           string   `yaml:"amazonAccessKeyID"`
	AmazonAccessKeySecret       string   `yaml:"amazonAccessKeySecret"`
	S3Region                    string   `yaml:"s3Region"`
	S3Endpoint                  string   `yaml:"s3Endpoint"`
	JWTSecretKey                string   `yaml:"jwtSecretKey"`
	BucketName                  string   `yaml:"bucketName"`
	Clients                     []string `yaml:"clients"`
	Consumer                    `yaml:"consumer"`
	DatabaseConfig              `yaml:"database"`
}

type Consumer struct {
	Network       string `yaml:"network"`
	BrokerAddress string `yaml:"brokerAddress"`
	Topic         string `yaml:"topic"`
	GroupID       string `yaml:"groupID"`
	Partition     int    `yaml:"partition"`
}

type DatabaseConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
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
