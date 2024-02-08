package controller

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"approval-service/logs"
	"approval-service/modules/entities/models"
)

type approvalHandler struct {
	approvalSrv models.ApprovalUsecase
}

var validate = validator.New()

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
	router.Get("/approvals/:requestId", controllers.GetApprovalByRequestID)
	router.Delete("/approval/:id", controllers.DeleteApproval)
	router.Post("/approval/create", controllers.CreateRequest)
	//mockData
	router.Post("/approval/Add/project", controllers.CreateProjectMock)
	router.Post("/approval/Add/user", controllers.CreateUserMock)
	router.Post("/approval/Add/task", controllers.CreateTaskMock)
	router.Get("/approval/project/projects", controllers.GetAllProject)
	router.Get("/approval/task/tasks", controllers.GetAllTask)
}

func (h *approvalHandler) UpdateStatus(c *fiber.Ctx) error {
	logs.Info("Put : Attempting to update approval status")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Warn("Error parsing approval ID:", zap.Error(err))
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
		logs.Warn("Error parsing update status approval body:", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	if err := validate.Struct(req); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
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
		logs.Warn(err.Error())
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
	logs.Info("Get : Attempting to ReceiveRequest approval ")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Warn("Error parsing approval ID:", zap.Error(err))
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
	requestUser := c.Query("senderId")
	if requestUser != "" {
		optional["sender_id"] = requestUser
	}
	status := c.Query("status")
	if status != "" {
		optional["status"] = status
	}
	projectId := c.Query("projectId")
	if projectId != "" {
		projectId, err := strconv.Atoi(projectId)
		if err != nil {
			logs.Warn("Error parsing approval ID:", zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(
				models.ResponseData{
					Message:    "project option is number",
					Status:     fiber.ErrBadRequest.Message,
					StatusCode: fiber.ErrBadRequest.Code,
				},
			)
		}
		optional["project_id"] = projectId
	}

	apprrovalReceive, err := h.approvalSrv.GetReceiveRequest(uint(id), optional)
	if err != nil {
		logs.Warn("Error can't get Receive approval ", zap.Error(err))
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
	logs.Info("Get : Attempting to get request approval")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Warn("Error parsing approval ID:", zap.Error(err))
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
	projectId := c.Query("projectId")
	if projectId != "" {
		projectId, err := strconv.Atoi(projectId)
		if err != nil {
			logs.Warn("Error parsing approval ID:", zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(
				models.ResponseData{
					Message:    "project option is number",
					Status:     fiber.ErrBadRequest.Message,
					StatusCode: fiber.ErrBadRequest.Code,
				},
			)
		}
		optional["project_id"] = projectId
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
		logs.Warn("Error get send request approval ", zap.Error(err))
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
	logs.Info("Delete : Attempting to delete approval")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Warn("Error parsing approval ID:", zap.Error(err))
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
		logs.Warn("Error delete approval ", zap.Error(err))
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
	logs.Info("PUT : Attempting to teamlead sent request to HR or Approver")

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
	if err := validate.Struct(requestBody); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}

	apprrovalUpdated, err := h.approvalSrv.SentRequest(uint(idApprove), requestBody)
	if err != nil {
		logs.Warn("Error can't send approval", zap.Error(err))
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
	logs.Info("GET : Attempting to GET approval By id")

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
	optional := map[string]interface{}{}

	//Optional
	projectId := c.Query("projectId")
	if projectId != "" {
		projectId, err := strconv.Atoi(projectId)
		if err != nil {
			logs.Warn("Error parsing approval ID:", zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(
				models.ResponseData{
					Message:    "project option is number",
					Status:     fiber.ErrBadRequest.Message,
					StatusCode: fiber.ErrBadRequest.Code,
				},
			)
		}
		optional["project_id"] = projectId
	}
	status := c.Query("status")
	if status != "" {
		optional["status"] = status
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
	logs.Info("GET : Attempting to GET approval By id successfully")
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
	logs.Info("GET : Attempting to get all approval")

	optional := map[string]interface{}{}
	//Optional
	projectId := c.Query("projectId")
	if projectId != "" {
		projectId, err := strconv.Atoi(projectId)
		if err != nil {
			logs.Warn("Error parsing approval ID:", zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(
				models.ResponseData{
					Message:    "project option is number",
					Status:     fiber.ErrBadRequest.Message,
					StatusCode: fiber.ErrBadRequest.Code,
				},
			)
		}
		optional["project_id"] = projectId
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

	apprrovalReceive, err := h.approvalSrv.GetAll(optional)
	if err != nil {
		logs.Warn("Error can't get all approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("get all approval successfully")
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
	logs.Info("GET : Attempting approval by id")

	id, err := c.ParamsInt("id")
	if err != nil {
		logs.Warn("Error parsing approval ID:", zap.Error(err))
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
	projectId := c.Query("projectId")
	if projectId != "" {
		projectId, err := strconv.Atoi(projectId)
		if err != nil {
			logs.Warn("Error parsing approval ID:", zap.Error(err))
			return c.Status(fiber.StatusNotFound).JSON(
				models.ResponseData{
					Message:    "project option is number",
					Status:     fiber.ErrBadRequest.Message,
					StatusCode: fiber.ErrBadRequest.Code,
				},
			)
		}
		optional["project_id"] = projectId
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

	apprrovalReceive, err := h.approvalSrv.GetByUserID(uint(id), optional)
	if err != nil {
		logs.Warn("Error can't get approval by id", zap.Error(err))
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

func (h *approvalHandler) CreateRequest(c *fiber.Ctx) error {
	logs.Info("POST : Attempting to create approval")

	approvalBody := new(models.CreateReq)
	if err := c.BodyParser(&approvalBody); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	if err := validate.Struct(approvalBody); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}

	apprrovalCreate, err := h.approvalSrv.CreateRequest(approvalBody)
	if err != nil {
		logs.Warn("Error can't Create approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("Create approval successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       apprrovalCreate,
		},
	)
}

func (h *approvalHandler) GetApprovalByRequestID(c *fiber.Ctx) error {
	logs.Info("GET : Attempting to Approva bu request id")

	id := c.Params("requestId")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ResponseData{
				Message:    "Require parameters",
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	// Convert id to UUID
	uuid, err := uuid.Parse(id)
	if err != nil {
		// Handle the error when the conversion fails
		return c.Status(fiber.StatusBadRequest).JSON(
			models.ResponseData{
				Message:    "Invalid UUID format",
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	res, err := h.approvalSrv.GetByRequestID(uuid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}
	logs.Info("GET : Attempting to Approva bu request id successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       res,
		},
	)
}

func (h *approvalHandler) CreateProjectMock(c *fiber.Ctx) error {
	logs.Info("Post : Attempting to by pass create project")
	ProjectReq := new(models.ProjectJson)
	if err := c.BodyParser(ProjectReq); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	if err := validate.Struct(ProjectReq); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	res, err := h.approvalSrv.CreateProject(ProjectReq)
	if err != nil {
		logs.Warn("Error can't Create approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}
	logs.Info("GET : Attempting to Approva bu request id successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       res,
		},
	)
}

func (h *approvalHandler) CreateUserMock(c *fiber.Ctx) error {
	logs.Info("Post : Attempting to by pass create project")
	req := new(models.UserProfile)
	if err := c.BodyParser(req); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	if err := validate.Struct(req); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	res, err := h.approvalSrv.CreateUser(req)
	if err != nil {
		logs.Warn("Error can't Create approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("GET : Attempting to Approva bu request id successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       res,
		},
	)
}

func (h *approvalHandler) CreateTaskMock(c *fiber.Ctx) error {
	logs.Info("Post : Attempting to by pass create task")
	req := new(models.Task)
	if err := c.BodyParser(req); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	if err := validate.Struct(req); err != nil {
		logs.Warn("Invalid request", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrBadRequest.Message,
				StatusCode: fiber.ErrBadRequest.Code,
			},
		)
	}
	res, err := h.approvalSrv.AddTask(req)
	if err != nil {
		logs.Warn("Error can't Create task ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("GET : Attempting to task bu request id successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       res,
		},
	)
}

func (h *approvalHandler) GetAllProject(c *fiber.Ctx) error {
	logs.Info("GET : Attempting to get all approval")

	apprrovalReceive, err := h.approvalSrv.GetAllProject()
	if err != nil {
		logs.Warn("Error can't get all approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("get all approval successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       apprrovalReceive,
		},
	)
}

func (h *approvalHandler) GetAllTask(c *fiber.Ctx) error {
	logs.Info("GET : Attempting to get all approval")

	apprrovalReceive, err := h.approvalSrv.GetAllTask()
	if err != nil {
		logs.Warn("Error can't get all approval ", zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(
			models.ResponseData{
				Message:    err.Error(),
				Status:     fiber.ErrNotFound.Message,
				StatusCode: fiber.ErrNotFound.Code,
			},
		)
	}

	logs.Info("get all approval successfully")
	return c.Status(fiber.StatusOK).JSON(
		models.ResponseData{
			Message:    "Succeed",
			Status:     "OK",
			StatusCode: fiber.StatusOK,
			Data:       apprrovalReceive,
		},
	)
}
