package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"

	"github.com/intermediate-service-ta/boot"
)

type Message struct {
	Command string
	Buffer  []byte
}

func ConsumeCommand(ctx context.Context, dep *boot.Dependencies) {
	kafkaLogFile, err := os.OpenFile("kafka-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer kafkaLogFile.Close()

	commandLogFile, err := os.OpenFile("command-log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer kafkaLogFile.Close()

	kafkaLog := log.New(kafkaLogFile, "kafka reader: ", 0)
	commandLog := log.New(commandLogFile, "kafka reader: ", 0)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{dep.Config().Consumer.BrokerAddress},
		GroupID:     dep.Config().Consumer.GroupID,
		Topic:       dep.Config().Consumer.Topic,
		Partition:   dep.Config().Consumer.Partition,
		MinBytes:    1,
		MaxBytes:    1e5,
		ErrorLogger: kafkaLog,
		MaxAttempts: 3,
	})

	for {
		select {
		case <-sigchan:
			kafkaLog.Println("Shutting down consumer...")
			return
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				kafkaLog.Println("Error reading message from Kafka:", err)
				continue
			}

			var msg Message
			if err := json.Unmarshal(message.Value, &msg); err != nil {
				kafkaLog.Println("failed to unmarshal:", err)
				continue
			}

			commandLog.Println(msg.Command, msg.Buffer, message.Offset)
			fmt.Println(msg.Command, msg.Buffer)
		}
	}
}
