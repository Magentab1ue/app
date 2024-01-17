package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// db
type Approval struct {
	gorm.Model
	ID           uint            `gorm:"primaryKey" json:"id"`
	RequestID    uuid.UUID       `json:"request_id"`
	To           string          `json:"to"`
	Approver     json.RawMessage `json:"approver"`
	Status       string          `json:"status"`
	Project      json.RawMessage `json:"project"`
	CreationDate time.Time       `json:"creation_date"`
	RequestUser  uint            `json:"request_user"`
	Task         json.RawMessage `json:"task"`
}
