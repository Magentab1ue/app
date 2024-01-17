package servers

// import (
// 	"time"
// 	_addRepo "profile-service/modules/additional/repository"
// 	_additonUse "profile-service/modules/additional/usecase"
// 	_acctiveRepo "profile-service/modules/externalDB/tcct_active_employees/repository"
// 	_deleteRepo "profile-service/modules/externalDB/tcct_deleted_employees/repository"
// 	"profile-service/modules/externalDB/tcct_orgInfo/repository"
// 	"profile-service/modules/organization_chart/controller"
// 	_orgRepo "profile-service/modules/organization_chart/repository"
// 	"profile-service/modules/organization_chart/usecase"
// 	_posiRepo "profile-service/modules/position/repository"
// 	_posiUse "profile-service/modules/position/usecase"
// 	_pubHan "profile-service/modules/producer/handlers"
// 	_produceUse "profile-service/modules/producer/usecase"
// 	_conPro "profile-service/modules/profile/controller"
// 	_profileRepo "profile-service/modules/profile/repository"
// 	_profileUsecase "profile-service/modules/profile/usecase"
//
//

// 	"github.com/gofiber/fiber/v2"
// )

// func (s *server) Handlers() error {

// 	// Group a version
// 	v1 := s.App.Group("/v1")

// 	//repo
// 	profileRepository := _profileRepo.NewprofileRepositoryDB(s.Db)
// 	AdditionalRepository := _addRepo.NewAdditionalRepositoryDB(s.Db)
// 	positionRepository := _posiRepo.NewPositionRepositoryDB(s.Db)
// 	orgRepo := _orgRepo.NewOrganizationChartRepositoryDB(s.Db)

// 	//externaldata
// 	// accRepo := _acctiveRepo.NewTcctActiveEmployeesDB(s.DbX)
// 	// deletedRepo := _deleteRepo.NewTcctDeletedEmployeesDB(s.DbX)
// 	// orgInfoRepo := repository.NewTcctOrgInfoDB(s.DbX)

// 	accMock := _acctiveRepo.NewTcctActiveEmployeesMock()
// 	deMock := _deleteRepo.NewTcctDeletedEmployeesMock()
// 	orgMock := repository.NewTcctOrgInfoMock()

// 	publisherHandler := _pubHan.NewEventProducer(s.SyncProducer)

// 	//usecase
// 	producerUsecase := _produceUse.NewProducerServiceUsers(publisherHandler)
// 	additonUse := _additonUse.NewAdditionalService(AdditionalRepository, s.Redis)
// 	posiUse := _posiUse.NewPositionService(positionRepository, s.Redis)
// 	orgUse := usecase.NewOrgService(orgRepo, s.Redis)
// 	profileUsecase := _profileUsecase.NewProfileService(
// 		profileRepository,
// 		producerUsecase,
// 		posiUse,
// 		s.Redis,
// 		additonUse,
// 		orgUse,
// 		s.Minio,
// 		accMock,
// 		deMock,
// 		orgMock)

// 	go func() {
// 		for {
// 			profileUsecase.UpdateData()
// 			// รอสักครู่ก่อนที่จะทำงานในรอบถัดไป
// 			time.Sleep(20 * time.Minute)
// 		}
// 	}()
// 	_conPro.NewProfileController(v1, profileUsecase)
// 	controller.NeworgController(v1, orgUse)

// 	// End point not found response
// 	s.App.Use(func(c *fiber.Ctx) error {
// 		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
// 			"status":      fiber.ErrInternalServerError.Message,
// 			"status_code": fiber.ErrInternalServerError.Code,
// 			"message":     "error, end point not found",
// 			"result":      nil,
// 		})
// 	})

// 	return nil

// }
