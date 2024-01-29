package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"approval-service/configs"
	"approval-service/logs"
	"approval-service/modules/servers"
	databases "approval-service/pkg/databases/postgres"
	redis "approval-service/pkg/databases/redis"
)

func main() {

	if err := godotenv.Load("../config.env"); err != nil {
		logs.Warn("Error loading .env file: %v", zap.Error(err))
	}

	cfg := new(configs.Config)

	configs.LoadConfigs(cfg)

	db, err := databases.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Append port to Kafka server addresses
	var kafkaServersWithPort []string
	for _, server := range cfg.Kafkas.Servers {
		serverWithPort := fmt.Sprintf("%s:%s", server, cfg.Kafkas.Port)
		kafkaServersWithPort = append(kafkaServersWithPort, serverWithPort)
	}
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.ClientID = cfg.Kafkas.ClientID

	producer, err := sarama.NewSyncProducer(kafkaServersWithPort, producerConfig)
	if err != nil {
		log.Fatal(err)
	}
	logs.Info(fmt.Sprintf("Connent kafka with server %v",kafkaServersWithPort))
	defer producer.Close()

	// Configure Kafka consumer
	consumerConfig := sarama.NewConfig()
	consumerConfig.ClientID = cfg.Kafkas.ClientID

	consumer, err := sarama.NewConsumerGroup(kafkaServersWithPort, cfg.Kafkas.Group, consumerConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	redis := redis.NewRedisClient(cfg)

	server := servers.NewServer(cfg, db, consumer, producer, redis)
	server.Start()

}
