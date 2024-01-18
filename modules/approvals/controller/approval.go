package controller

import (
	"strconv"

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
	// //router.Post("/hr-profile", controllers.PullData)
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
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}

	req := new(models.UpdateStatusReq)
	err = c.BodyParser(req)
	if err != nil {
		logs.Error("Error parsing update status approval body:", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}

	apprrovalUpdated, err := h.approvalSrv.UpdateStatus(uint(id), req)
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
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}

	optional := map[string]interface{}{}

	//Optional
	requestUser := c.Query("requestUser")
	if requestUser != "" {
		optional["requestUser"] = requestUser
	}

	apprrovalReceive, err := h.approvalSrv.GetReceiveRequest(uint(id), optional)
	if err != nil {
		logs.Error("Error can't get Receive approval ", zap.Error(err))
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
			Data:       apprrovalReceive,
		},
	)
}

func (h *approvalHandler) SendRequest(c *fiber.Ctx) error {
	logs.Info("Attempting to get request approval")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Error("Error parsing approval ID:", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	optional := map[string]interface{}{}

	//Optional
	project := c.Query("project")
	if project != "" {
		optional["project"] = project
	}
	to := c.Query("to")
	if to != "" {
		to, _ := strconv.Atoi(to)
		optional["to"] = []uint{uint(to)}
	}

	apprrovalSend, err := h.approvalSrv.GetSendRequest(uint(id), optional)
	if err != nil {
		logs.Error("Error get send request approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("get send approval successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       apprrovalSend,
		},
	)
}

func (h *approvalHandler) DeleteApproval(c *fiber.Ctx) error {
	logs.Info("Attempting to delete approval")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Error("Error parsing approval ID:", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}

	err = h.approvalSrv.DeleteApproval(uint(id))
	if err != nil {
		logs.Error("Error delete approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("delete approval successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Deleted Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
		},
	)
}
