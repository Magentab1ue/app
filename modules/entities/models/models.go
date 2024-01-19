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
	To           pq.Int64Array `json:"to"`
	CreationDate time.Time     `json:"creation_date"`
	RequestUser  uint          `json:"request_user"`
	IsSignature  bool          `json:"is_signature"`
}

type RequestReq struct {
	To           pq.Int64Array  `json:"to" gorm:"type:integer[]"`
	Project      datatypes.JSON `json:"project" gorm:"type:jsonb"`
	CreationDate time.Time      `json:"creation_date"`
	RequestUser  uint           `json:"request_user"`
	Task         datatypes.JSON `json:"task" gorm:"type:jsonb"`
}

// db
type Approvals struct {
	gorm.Model
	RequestID    uuid.UUID      `json:"request_id"`
	To           pq.Int64Array  `json:"to" gorm:"type:integer[]"`
	Approver     uint           `json:"approver"`
	Status       string         `json:"status"`
	Project      datatypes.JSON `json:"project" gorm:"type:jsonb"` // Assuming your database supports JSONB
	CreationDate time.Time      `json:"creation_date"`
	RequestUser  uint           `json:"request_user"`
	IsSignature  bool           `json:"is_signature"`
	Task         datatypes.JSON `json:"task" gorm:"type:jsonb"` // Assuming your database supports JSONB
}

type UpdateStatusReq struct {
	IsSignature bool   `json:"is_signature"`
	Status      string `json:"status"`
	Approver    uint
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
