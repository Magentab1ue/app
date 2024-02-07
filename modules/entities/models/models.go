package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// request
type RequestSentRequest struct {
	SenderID    uint   `json:"senderId" validate:"required,numeric,min=1"`
	IsSignature bool   `json:"isSignature"`
	Name        string `json:"name" validate:"required"`
	Detail      string `json:"detail" `
	ToRole      string `json:"toRole" validate:"required"`
}

type TaskRequest struct {
	Status         string `json:"status"`
	ApprovalStatus string `json:"approvalStatus"`
	UserID         uint   `json:"userId"`
	ProjectId      uint   `json:"projectId"`
	Detail         string `db:"detail" json:"detail"`
}

type CreateReq struct {
	ProjectId uint           `json:"projectId"  validate:"required,numeric,min=1"`
	SenderID  uint           `json:"senderId" validate:"required,numeric,min=1"`
	Task      datatypes.JSON `json:"task" gorm:"type:jsonb" validate:"required"`
	Name      string         `json:"name" validate:"required"`
	Detail    string         `json:"detail" `
}

type UpdateStatusReq struct {
	IsSignature bool   `json:"isSignature" `
	Status      string `json:"status" validate:"required"`
	Approver    uint   `json:"approver" validate:"required"`
}

type ResponseData struct {
	Message    string      `json:"message"`
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

type ProduceReq struct {
	ID   uint
	Data interface{}
}

const (
	Pending string = "pending"
	Approve string = "approved"
	Reject  string = "reject"
)

type ConsumerOffset struct {
	gorm.Model
	Topic     string
	Offset    int64
	Partition int32
}

// db
type Approvals struct {
	gorm.Model
	RequestID    uuid.UUID      `json:"requestId"`
	To           pq.Int64Array  `json:"to" gorm:"type:integer[]"`
	Approver     uint           `json:"approver"`
	Status       string         `json:"status"`
	CreationDate time.Time      `json:"creationDate"`
	IsSignature  bool           `json:"isSignature"`
	Task         datatypes.JSON `json:"task" gorm:"type:jsonb"` // Assuming your database supports JSONB
	Name         string         `json:"name"`                   // name timesheet
	Detail       string         `json:"detail"`                 //detail timesheet
	ToRole       string         `json:"toRole"`
	SenderID     uint           `json:"senderId"`
	ProjectID    uint           `json:"projectId"`
}

type UserProfile struct {
	gorm.Model
	ProfileId uint        `json:"profileId" validate:"required"`
	Name      string      `json:"name" validate:"required"`
	Approvals []Approvals `gorm:"foreignKey:SenderID"`
}

type Project struct {
	gorm.Model
	Project   datatypes.JSON `json:"project" gorm:"type:jsonb"`
	Approvals []Approvals    `gorm:"foreignKey:ProjectID"`
}
type ProjectJson struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	TeamLeads []struct {
		ID int `json:"id"`
	} `json:"teamleads"`
	Approvers []struct {
		ID   int      `json:"id"`
		Role []string `json:"roles"`
	} `json:"approvers"`
	Members []struct {
		ID int `json:"id"`
	} `json:"members"`
}

type Task struct {
	gorm.Model
	Status         string `json:"status"`
	ApprovalStatus string `json:"approvalStatus"`
	UserID         uint   `json:"userId"`
	ProjectId      uint   `json:"projectId"`
	Detail         string `db:"detail" json:"detail"`
}

var TaskAppproveStatusMap = map[int]string{
	0: "open",
	1: "waiting",
	2: "approved",
	3: "reject",
}

var TaskStatusMap = map[int]string{
	0: "OPEN",
	1: "DONE",
}

type ReqTask struct {
	ID             uint   `json:"id"`
	Status         int    `json:"status"`
	ApprovalStatus int    `json:"approvalStatus"`
	Detail         string `db:"detail" json:"detail"`
}
