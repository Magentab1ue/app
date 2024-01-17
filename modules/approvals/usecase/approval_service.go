package usecase

import "approval-service/modules/entities/models"

type approvalService struct {
	approvalRepo models.ApprovalRepository
}

func NewApprovalService(
	approvalRepo models.ApprovalRepository,
) models.ApprovalUsecase {
	return &approvalService{approvalRepo}
}

func (u approvalService) Update(id int, req *models.UpdateStatusReq) (*models.Approval, error) {

	return nil, nil
}

func (u approvalService) ReceiveRequest(id int, optional map[string]interface{}) ([]models.Approval, error) {

	return nil, nil
}
