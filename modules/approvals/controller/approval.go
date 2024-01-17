package controller

import (
	"github.com/gofiber/fiber/v2"

	"approval-service/modules/entities/models"
)

type approvalHandler struct {
	approvalSrv models.ApprovalUsecase
}

func NewApprovalController(router fiber.Router, approvalSrv models.ApprovalUsecase) {
	controllers := &approvalHandler{
		approvalSrv: approvalSrv,
	}
	_ = controllers//gfffffffffffffffffffffffffffffffffffffdsssssssssssssssssssssssssssssssssssssssssssssssssss
	// router.Post("/profile", controllers.newProfile)
	// //router.Post("/hr-profile", controllers.PullData)
	// router.Get("/profiles", controllers.getAllProfileData)
	// router.Get("/profile/:profileId", controllers.GetProfileByID)
	// router.Get("/profiles/:role", controllers.GetProfileByRole)
	// router.Get("/profiles/:email", controllers.GetProfileByEmail)
	// router.Put("/profile/:profileId", controllers.Update)
	// router.Delete("/profile/:profileId", controllers.Delete)
	// router.Get("/test", controllers.Test)

}
