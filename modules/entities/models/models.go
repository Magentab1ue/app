package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
//request
type RequestSentRequest struct{
	To           []uint          `json:"to"`
	CreationDate time.Time       `json:"creation_date"`
	RequestUser  uint            `json:"request_user"`
}

// db
type Approval struct {
	gorm.Model
	ID           uint            `gorm:"primaryKey" json:"id"`
	RequestID    uuid.UUID       `json:"request_id"`
	To           []uint          `json:"to"`
	Approver     uint            `json:"approver"`
	Status       string          `json:"status"` //
	Project      json.RawMessage `json:"project"`
	CreationDate time.Time       `json:"creation_date"`
	RequestUser  uint            `json:"request_user"`
	IsSignature  bool            `json:"is_signature"`
	Task         json.RawMessage `json:"task"`
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
