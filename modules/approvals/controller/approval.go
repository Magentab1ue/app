package controller

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"approval-service/logs"
	"approval-service/modules/entities/models"
)

type approvalHandler struct {
	approvalSrv models.ApprovalUsecase
}

func NewApprovalController(router fiber.Router, approvalSrv models.ApprovalUsecase) {
	controllers := &approvalHandler{
		approvalSrv: approvalSrv,
	}
	_ = controllers //gfffffffffffffffffffffffffffffffffffffdsssssssssssssssssssssssssssssssssssssssssssssssssss
	// router.Post("/profile", controllers.newProfile)swwrfwaefaewrfaewrf
	// //router.Post("/hr-profile", controllers.PullData)fffffffffffffffffffffffdfdf
	// router.Get("/profiles", controllers.getAllProfileData)
	// router.Get("/profile/:profileId", controllers.GetProfileByID)
	// router.Get("/profiles/:role", controllers.GetProfileByRole)
	// router.Get("/profiles/:email", controllers.GetProfileByEmail)
	// router.Put("/profile/:profileId", controllers.Update)
	// router.Delete("/profile/:profileId", controllers.Delete)
	// router.Get("/test", controllers.Test)

}

func (h *approvalHandler) UpdateStatus(c *fiber.Ctx) error {
	logs.Info("Attempting to update approval status")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Error("Error parsing approval ID:", zap.Error(err))
	}

	req := new(models.UpdateStatusReq)
	err = c.BodyParser(req)
	if err != nil {
		logs.Error("Error parsing update status approval body:", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	apprrovalUpdated, err := h.approvalSrv.Update(id,req)
	if err != nil {
		logs.Error("Error update status approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}
	logs.Info("update approval status successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       apprrovalUpdated,
		},
	)
}


func (h *approvalHandler) ReceiveRequest(c *fiber.Ctx) error {
	logs.Info("Attempting to update approval status")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Error("Error parsing approval ID:", zap.Error(err))
	}
	optional := map[string]interface{}{}

	//Optional
	requestUser := c.Query("requestUser")
	if requestUser != "" {
		optional["requestUser"] = requestUser
	}
	
	apprrovalUpdated, err := h.approvalSrv.ReceiveRequest(id,optional)
	if err != nil {
		logs.Error("Error update status approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("get Receive approval successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       apprrovalUpdated,
		},
	)
}
