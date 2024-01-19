package usecase

import (
	"approval-service/modules/entities/models"
)

type consumerUsecase struct {
	approvalRepo models.ApprovalRepository
}

func NewConsumerUsecase(approvalRepo models.ApprovalRepository) models.ConsumerUsecase {
	return &consumerUsecase{approvalRepo}
}

// func (u consumerUsecase) RequestCreated(event events.RequestCreatedEvent) error {
// 	_, err := u.approvalRepo.Create(&models.Approvals{
// 		RequestID:    uuid.New(),
// 		To:           event.To,
// 		Status:       "pending",
// 		CreationDate: event.CreationDate,
// 		Project:      event.Project,
// 		RequestUser:  event.RequestUser,
// 		Task:         event.Task,
// 	})
// 	if err != nil {
// 		logs.Error(err)
// 		return errors.New("Can't create request")
// 	}
// 	return nil
// }
