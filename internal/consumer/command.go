package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/intermediate-service-ta/boot"
	"github.com/intermediate-service-ta/internal/repository"
)

var errLogFile *os.File

type Message struct {
	Command       string
	Token         string
	AbsPathSource string
	AbsPathDest   string
	FileMode      uint64
	Buffer        []byte
}

type Consumer struct {
	fileRepo repository.FileRepository
	errorLog *log.Logger
}

func NewConsumer(fileRepo repository.FileRepository) *Consumer {
	var errFile error
	errLogFile, errFile = os.OpenFile("error-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errFile != nil {
		log.Fatalf("error opening file: %v", errFile)
	}

	return &Consumer{fileRepo: fileRepo}
}

func (con *Consumer) ConsumeCommand(c context.Context, dep *boot.Dependencies, sigchan chan os.Signal) {
	kafkaLogFile, err := os.OpenFile("kafka-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer kafkaLogFile.Close()

	commandLogFile, err := os.OpenFile("command-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer commandLogFile.Close()

	kafkaLog := log.New(kafkaLogFile, "kafka reader: ", 0)
	commandLog := log.New(commandLogFile, "kafka reader: ", 0)
	con.errorLog = log.New(errLogFile, "error: ", 0)

	consumerConf := dep.Config().Consumer
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{consumerConf.BrokerAddress},
		GroupID:     consumerConf.GroupID,
		Topic:       consumerConf.Topic,
		Partition:   consumerConf.Partition,
		MinBytes:    1,
		MaxBytes:    1e9,
		ErrorLogger: kafkaLog,
	})
	defer reader.Close()
	defer errLogFile.Close()

	ctx := context.Background()
	for {
		select {
		case <-sigchan:
			kafkaLog.Println("Shutting down consumer...")
			return
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				kafkaLog.Println("Error reading message from Kafka:", err)
			}

			var msg Message
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				kafkaLog.Println("failed to unmarshal:", err)
			}

			con.AuthQueue(c, msg, commandLog)
			time.Sleep(300 * time.Millisecond)
		}
	}
}
