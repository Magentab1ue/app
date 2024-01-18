package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// request
type RequestSentRequest struct {
	To           pq.Int64Array `json:"to"`
	CreationDate time.Time     `json:"creation_date"`
	RequestUser  uint          `json:"request_user"`
	IsSignature  bool          `json:"is_signature"`
}

// db
type Approvals struct {
	gorm.Model
	RequestID    uuid.UUID       `json:"request_id"`
	To           pq.Int64Array   `json:"to" gorm:"type:integer[]"`
	Approver     uint            `json:"approver"`
	Status       string          `json:"status"`
	Project      json.RawMessage `json:"project" gorm:"type:jsonb"` // Assuming your database supports JSONB
	CreationDate time.Time       `json:"creation_date"`
	RequestUser  uint            `json:"request_user"`
	IsSignature  bool            `json:"is_signature"`
	Task         json.RawMessage `json:"task" gorm:"type:jsonb"` // Assuming your database supports JSONB
}

type UpdateStatusReq struct {
	Status   string `json:"status"`
	Approver uint
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
