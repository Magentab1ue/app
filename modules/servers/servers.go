package servers

// import (
// 	"profile-service/configs"
// 	"profile-service/logs"
// 	"profile-service/pkg/utils"
//
//

// 	"github.com/IBM/sarama"
// 	"github.com/go-redis/redis/v8"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/minio/minio-go/v7"
// 	"go.uber.org/zap"
// 	"gorm.io/gorm"
// )

// type server struct {
// 	App                  *fiber.App
// 	Db                   *gorm.DB
// 	Cfg                  *configs.Config
// 	ConsumerGroup        sarama.ConsumerGroup
// 	SyncProducer         sarama.SyncProducer
// 	consumerGroupHandler sarama.ConsumerGroupHandler
// 	Redis                *redis.Client
// 	Minio                *minio.Client
// 	DbX                  *gorm.DB
// }

// func NewServer(cfg *configs.Config,
// 	db *gorm.DB,
// 	consumerGroup sarama.ConsumerGroup,
// 	syncProducer sarama.SyncProducer,
// 	redis *redis.Client,
// 	Minio *minio.Client,
// 	dbX *gorm.DB,
// ) *server {
// 	return &server{
// 		App:           fiber.New(),
// 		Db:            db,
// 		Cfg:           cfg,
// 		ConsumerGroup: consumerGroup,
// 		SyncProducer:  syncProducer,
// 		Redis:         redis,
// 		Minio:         Minio,
// 		DbX:           dbX,
// 	}
// }

// func (s *server) Start() {
// 	if err := s.Handlers(); err != nil {
// 		logs.Error(err)
// 		panic(err.Error())
// 	}

// 	fiberConnURL, err := utils.UrlBuilder("fiber", s.Cfg)
// 	if err != nil {
// 		logs.Error(err)
//

// 	port := s.Cfg.App.Port
// 	logs.Info("server started on localhost:", zap.String("port", port))

// 	if err := s.App.Listen(fiberConnURL); err != nil {
// 		logs.Error(err)
// 		panic(err.Error())
// 	}
// }
