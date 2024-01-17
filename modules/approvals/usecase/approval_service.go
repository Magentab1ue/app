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