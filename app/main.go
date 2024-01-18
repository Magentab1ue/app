package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"

	"approval-service/configs"
	"approval-service/modules/servers"
	databases "approval-service/pkg/databases/postgres"
	redis "approval-service/pkg/databases/redis"
)

func main() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := new(configs.Config)

	configs.LoadConfigs(cfg)

	db, err := databases.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	producer, err := sarama.NewSyncProducer(cfg.Kafkas.Hosts, nil)
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
