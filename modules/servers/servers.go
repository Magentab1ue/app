package servers

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
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
	logs.Info(fmt.Sprintf("Starting comsuming kafa server: %v Port : %v Group: %v", s.Cfg.Kafkas.Servers, s.Cfg.Kafkas.Port, s.Cfg.Kafkas.Group))
	logs.Info(fmt.Sprintf("Subscribed topics: %s", events.SubscribedTopics))
	// go func() {
	// 	for {
	// 		err := s.ConsumerGroup.Consume(context.Background(), events.SubscribedTopics, s.consumerGroupHandler)
	// 		if err != nil {
	// 			logs.Error(fmt.Sprintf("Can't consuming with : server: %v Port : %v Group: %v", s.Cfg.Kafkas.Servers, s.Cfg.Kafkas.Port, s.Cfg.Kafkas.Group),zap.Error(err))
	// 		}
	// 	}
	// }()

	logs.Info(fmt.Sprintf("server started on %s ", fiberConnURL))

	if err := s.App.Listen(fiberConnURL); err != nil {
		logs.Error(err)
		panic(err.Error())
	}
}
