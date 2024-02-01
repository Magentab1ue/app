package models

import (
	"time"

	"github.com/google/uuid"

	"approval-service/modules/entities/events"
)

type ApprovalUsecase interface {
	UpdateStatus(uint, *UpdateStatusReq) (*Approvals, error)
	GetReceiveRequest(uint, map[string]interface{}) ([]Approvals, error)
	GetSendRequest(uint, map[string]interface{}) ([]Approvals, error)
	DeleteApproval(uint) error
	GetByID(uint) (*Approvals, error)
	SentRequest(uint, *RequestSentRequest) (*Approvals, error)
	GetAll(map[string]interface{}) ([]Approvals, error)
	GetByUserID(uint, map[string]interface{}) ([]Approvals, error)
	CreateRequest(*CreateReq) (*Approvals, error)
	GetByRequestID(uuid.UUID) ([]Approvals, error)
}

type ApprovalRepository interface {
	Create(*Approvals) (*Approvals, error)
	GetByID(uint) (*Approvals, error)
	GetReceiveRequest(uint, map[string]interface{}) ([]Approvals, error)
	Update(*Approvals) (*Approvals, error)
	GetSendRequest(userId uint, optional map[string]interface{}) ([]Approvals, error)
	DeleteApproval(requestId uint) (*Approvals, error)
	GetAll(map[string]interface{}) ([]Approvals, error)
	GetByUserID(uint, map[string]interface{}) ([]Approvals, error)
	GetByRequestID(uuid.UUID) ([]Approvals, error)
	GetByRequestIDLast(uuid.UUID) (*Approvals, error)
	GetProjectById(id uint) (*Project, error)
}

type ProfileRepositoryDB interface {
	Create(request *UserProfile) error
	Update(req *UserProfile) error
	Delete(Id uint) error
	GetByID(id uint) (*UserProfile, error)
}

type ProducerApproval interface {
	RequestCreated(*ProduceReq, time.Time) error
	ApprovalUpdated(*ProduceReq, time.Time) error
	ApprovalDeleted(uint) error
}

type ConsumerUsecase interface {
	CreateProfile(e events.UserProfile) error
	UpdateProfile(e events.UserProfile) error
	DeleteProfile(e events.UserProfileDeleted) error
	CreateProject(e events.ProjectEvent) error
	UpdateProject(e events.ProjectEvent) error
	DeleteProject(e events.ProjectEventDeleted) error
	//CheckOffsetMessage(topic string, offset int64, partition int32) error
}

type ConsumerRepository interface {
	Get(req *ConsumerOffset) (*ConsumerOffset, error)
	Create(req *ConsumerOffset) error
}

// consumer
type EventHandlerConsume interface {
	// CheckMessage(msg *sarama.ConsumerMessage) error
	Handle(toppic string, eventByte []byte) error
}
type EventProducer interface {
	Produce(events.Event) error
}

type ProjectRepositoryDB interface {
	Create(request *Project) error
	Update(req *Project) error
	Delete(Id uint) error
}
