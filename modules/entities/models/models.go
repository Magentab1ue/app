package models

import (
	"time"

	"gorm.io/gorm"
)

// db
type Approval struct {
	gorm.Model
	ID           uint `gorm:"primaryKey" json:"id"`
	RequestID    uint `json:"request_id"`
	To           string `json:"to"`
	Approver     interface{} `json:"approver"`
	Status       string `json:"status"`
	Project      interface{} `json:"project"`
	CreationDate time.Time `json:"creation_date"`
	RequestUser  uint `json:"request_user"`
	Task         interface{} `json:"task"`
}
