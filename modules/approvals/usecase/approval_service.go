package usecase

import (
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
)

type approvalService struct {
	approvalRepo models.ApprovalRepository
	produce      models.EventProducer
}

func NewApprovalService(
	approvalRepo models.ApprovalRepository,
	produce models.EventProducer,
) models.ApprovalUsecase {
	return &approvalService{approvalRepo, produce}
}

func (u approvalService) UpdateStatus(id uint, req *models.UpdateStatusReq) (*models.Approval, error) {
	approvalRes, err := u.approvalRepo.UpdateStatus(id, req)
	if err != nil {
		return nil, err
	}
	event := events.ApprovalUpdatedEvent{
		RequestId: id,
		Approver:  req.Approver,
	}
	err = u.produce.Produce(event)
	if err != nil {
		return nil, err
	}
	return approvalRes, nil
}

func (u approvalService) ReceiveRequest(id int, optional map[string]interface{}) ([]models.Approval, error) {

	return nil, nil
}
