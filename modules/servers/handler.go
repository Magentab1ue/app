package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"approval-service/modules/approvals/controller"
	approvalRepository "approval-service/modules/approvals/repository"
	_appSrv "approval-service/modules/approvals/usecase"
	consumerHandler "approval-service/modules/consumer/handlers"
	_consumerUsecase "approval-service/modules/consumer/usecase"
	_handlerProducer "approval-service/modules/producer/handlers"
	profileRepository "approval-service/modules/profile/repository"
	projectRepository "approval-service/modules/project/repository"
)

func (s *server) Handlers() error {

	//middleware cors allow
	s.App.Use(cors.New(cors.Config{
		AllowOrigins: s.Cfg.Cors.AllowOrigins,
		AllowHeaders: "Content-Type",
	}))

	// Group a version
	v1 := s.App.Group("/v1")
	//repo
	approveRepo := approvalRepository.NewapprovalRepositoryDB(s.Db)
	profileRepo := profileRepository.NewprofileRepositoryDB(s.Db)
	projectRepo := projectRepository.NewproProjectRepositoryDB(s.Db)
	//consumRepo := consumerRepository.NewConsumerRepository(s.Db)

	// consumer
	consumeUsecase := _consumerUsecase.NewConsumerUsecase(profileRepo, projectRepo)
	eventHandlerConsumer := consumerHandler.NewEventHandler(consumeUsecase)
	s.consumerGroupHandler = consumerHandler.NewHandlerConsumeGroup(eventHandlerConsumer)

	// producer
	handlerProducer := _handlerProducer.NewEventProducer(s.SyncProducer)
	//producerUsecase := _publisherUsecase.NewProducerServiceApprovals(handlerProducer)

	//service
	approveSrv := _appSrv.NewApprovalService(
		approveRepo,
		handlerProducer,
		s.Redis,
	)

	controller.NewApprovalController(v1, approveSrv)

	// End point not found response
	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "error, end point not found",
			"result":      nil,
		})
	})

	return nil

}
