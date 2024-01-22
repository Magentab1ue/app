package servers

import (
	"github.com/gofiber/fiber/v2"

	"approval-service/modules/approvals/controller"
	"approval-service/modules/approvals/repository"
	_appSrv "approval-service/modules/approvals/usecase"
	consumerHandler "approval-service/modules/consumer/handlers"
	_consumerUsecase "approval-service/modules/consumer/usecase"
	_handlerProducer "approval-service/modules/producer/handlers"
)

func (s *server) Handlers() error {

	// Group a version
	v1 := s.App.Group("/v1")
	//repo
	approveRepo := repository.NewapprovalRepositoryDB(s.Db)

	// consumer
	consumeUsecase := _consumerUsecase.NewConsumerUsecase(approveRepo)
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
