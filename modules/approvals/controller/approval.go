package controller

import (
	"fmt"
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
	router.Get("/approval/:profileId", controllers.GetApprovalByID)
	router.Get("/approvals", controllers.GetAllApproval)
	router.Put("/approval/update-status/:id", controllers.UpdateStatus)
	router.Get("/approval/user/:id", controllers.GetByUserID)
	router.Put("/approval/send-request/:id", controllers.RequestSent)
	router.Get("/approval/user/:id/receive-request", controllers.ReceiveRequest)
	router.Get("/approval/user/:id/send-request", controllers.SendRequest)
	router.Delete("/approval/:id", controllers.DeleteApproval)
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
		optional["request_user"] = requestUser
	}
	projectId := c.Query("project")
	if projectId != "" {
		projectId, err := strconv.Atoi(projectId)
		if err != nil {
			logs.Error("Error parsing approval ID:", zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(
				models.ResponseData{
					Message:    "project option is number",
					Status:     fiber.ErrBadRequest.Message,
					StatusCode: fiber.ErrBadRequest.Code,
				},
			)
		}
		optional["project"] = fmt.Sprintf(`{"id":%d}`, projectId)
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
	projectId := c.Query("project")
	if projectId != "" {
		projectId, err := strconv.Atoi(projectId)
		if err != nil {
			logs.Error("Error parsing approval ID:", zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(
				models.ResponseData{
					Message:    "project option is number",
					Status:     fiber.ErrBadRequest.Message,
					StatusCode: fiber.ErrBadRequest.Code,
				},
			)
		}
		optional["project"] = fmt.Sprintf(`{"id":%d}`, projectId)
	}

	to := c.Query("to")
	if to != "" {
		to, _ := strconv.Atoi(to)
		optional["to"] = to
	}
	status := c.Query("status")
	if status != "" {
		optional["status"] = status
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

func (h *approvalHandler) RequestSent(c *fiber.Ctx) error {
	logs.Info("Attempting to teamlead sent request to HR or Approver")

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ResponseData{
				Message:    "Require parameters",
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	idApprove, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	var requestBody *models.RequestSentRequest
	if err := c.BodyParser(&requestBody); err != nil {
		logs.Info("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	apprrovalUpdated, err := h.approvalSrv.SentRequest(uint(idApprove), requestBody)
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

	logs.Info("get send approval successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       apprrovalUpdated,
		},
	)
}

func (h *approvalHandler) GetApprovalByID(c *fiber.Ctx) error {

	id := c.Params("profileId")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ResponseData{
				Message:    "Require parameters",
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	idProfile, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ResponseData{
				Message:    "Invalid parameter type",
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	res, err := h.approvalSrv.GetByID(uint(idProfile))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       res,
		},
	)
}

func (h *approvalHandler) GetAllApproval(c *fiber.Ctx) error {
	logs.Info("Attempting to update approval status")

	optional := map[string]interface{}{}
	//Optional
	status := c.Query("status")
	if status != "" {
		optional["status"] = status
	}
	to := c.Query("to")
	if to != "" {
		to, _ := strconv.Atoi(to)
		optional["to"] = []uint{uint(to)}
	}

	apprrovalReceive, err := h.approvalSrv.GetAll(optional)
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

func (h *approvalHandler) GetByUserID(c *fiber.Ctx) error {
	logs.Info("Attempting to update approval status")

	optional := map[string]interface{}{}

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

	//Optional
	status := c.Query("status")
	if status != "" {
		optional["status"] = status
	}
	to := c.Query("to")
	if to != "" {
		to, _ := strconv.Atoi(to)
		optional["to"] = []uint{uint(to)}
	}

	apprrovalReceive, err := h.approvalSrv.GetByUserID(uint(id), optional)
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
