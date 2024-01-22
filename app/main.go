package main

import (
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

	if err := godotenv.Load(); err != nil {

		logs.Warn("Error loading .env file: %v", zap.Error(err))
	}

	cfg := new(configs.Config)

	configs.LoadConfigs(cfg)

	db, err := databases.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	kafkaConfig := sarama.NewConfig()
    kafkaConfig.ClientID = cfg.Kafkas.ClientID

	producer, err := sarama.NewSyncProducer(cfg.Kafkas.Hosts, kafkaConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	consumer, err := sarama.NewConsumerGroup(cfg.Kafkas.Hosts, cfg.Kafkas.Group, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	redis := redis.NewRedisClient(cfg)

	server := servers.NewServer(cfg, db, consumer, producer, redis)
	server.Start()

}
