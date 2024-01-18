package servers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"approval-service/configs"
	"approval-service/logs"
	"approval-service/modules/entities/events"
	"approval-service/pkg/utils"
)

type server struct {
	App                  *fiber.App
	Db                   *gorm.DB
	Cfg                  *configs.Config
	ConsumerGroup        sarama.ConsumerGroup
	SyncProducer         sarama.SyncProducer
	consumerGroupHandler sarama.ConsumerGroupHandler
	Redis                *redis.Client
}

func NewServer(cfg *configs.Config,
	db *gorm.DB,
	consumerGroup sarama.ConsumerGroup,
	syncProducer sarama.SyncProducer,
	redis *redis.Client,
) *server {
	return &server{
		App:           fiber.New(),
		Db:            db,
		Cfg:           cfg,
		ConsumerGroup: consumerGroup,
		SyncProducer:  syncProducer,
		Redis:         redis,
	}
}

func (s *server) Start() {
	if err := s.Handlers(); err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	fiberConnURL, err := utils.UrlBuilder("fiber", s.Cfg)
	if err != nil {
		logs.Error(err)
		panic(err.Error())
	}

	// Start consumer
	go func() {
		log.Info("Connect to kafa server:", s.Cfg.Kafkas.Hosts, ",Group:", s.Cfg.Kafkas.Group)
		log.Info("Subscribed topics:", events.SubscribedTopics)
		for {
			s.ConsumerGroup.Consume(context.Background(), events.SubscribedTopics, s.consumerGroupHandler)
		}
	}()

	port := s.Cfg.App.Port
	logs.Info("server started on localhost:", zap.String("port", port))

	if err := s.App.Listen(fiberConnURL); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}
