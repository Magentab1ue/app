package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Event interface {
	String() string
}

var SubscribedTopics = []string{
	UserProfile{}.TopicCreate(),
	UserProfile{}.TopicUpdate(),
	UserProfileDeleted{}.TopicDelete(),
	ProjectEvent{}.TopicCreate(),
	ProjectEvent{}.TopicUpdate(),
	ProjectEventDeleted{}.TopicDelete(),
	TaskEvent{}.TaskEventCreated(),
	TaskEvent{}.TaskEventUpdated(),
	TaskEvent{}.TaskEventDeleted(),
}

type UserProfile struct {
	ProfileId uint   `json:"profileId" validate:"required"`
	Name      string `json:"Name" validate:"required"`
	UserId    uint   `json:"userId" validate:"required"`
}

func (UserProfile) TopicCreate() string {
	return "tcchub.profile.created"
}
func (UserProfile) TopicUpdate() string {
	return "tcchub.profile.updated"
}

type UserProfileDeleted struct {
	UserId uint `json:"userId" validate:"required"`
}

func (UserProfileDeleted) TopicDelete() string {
	return "tcchub.profile.deleted"
}

type RequestCreatedEvent struct {
	ID           uint           `json:"id"`
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

func (RequestCreatedEvent) String() string {
	return "tcchub.approval.created"
}

type ApprovalUpdatedEvent struct {
	ID           uint           `json:"id"`
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

func (ApprovalUpdatedEvent) String() string {
	return "tcchub.approval.updated"
}

type ApprovalDeletedEvent struct {
	Task datatypes.JSON `json:"task"`
}

func (ApprovalDeletedEvent) String() string {
	return "tcchub.approval.deleted"
}

type ProjectEvent struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	TeamLeads []struct {
		ID int `json:"id"`
	} `json:"teamleads"`
	Approvers []struct {
		ID   int `json:"id"`
		Role []string
	} `json:"approvers"`
	Members []struct {
		ID int `json:"id"`
	} `json:"members"`
}

func (ProjectEvent) TopicCreate() string {
	return "tcchub.project.created"
}
func (ProjectEvent) TopicUpdate() string {
	return "tcchub.project.updated"
}

type ProjectEventDeleted struct {
	ID uint `json:"id"`
}

func (ProjectEventDeleted) TopicDelete() string {
	return "tcchub.project.deleted"
}

type TaskEvent struct {
	ID             int    `json:"taskId"`
	UserID         string `json:"userId"`
	ProjectId      string `json:"projectId"`
	Status         int    `json:"status"`
	ApprovalStatus int    `json:"approvalStatus"`
	Detail         string ` json:"detail"`
}

func (TaskEvent) TaskEventCreated() string {
	return "tcchub.task.task.created"
}
func (TaskEvent) TaskEventUpdated() string {
	return "tcchub.task.task.updated"
}
func (TaskEvent) TaskEventDeleted() string {
	return "tcchub.task.task.deleted"
}
