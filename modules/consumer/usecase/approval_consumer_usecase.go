package usecase

import (
	"approval-service/logs"
	"approval-service/modules/entities/events"
	"approval-service/modules/entities/models"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type consumerUsecase struct {
	//consumerRepo models.ConsumerRepository
	profileRepo models.ProfileRepositoryDB
	projectRepo models.ProjectRepositoryDB
}

func NewConsumerUsecase(profileRepo models.ProfileRepositoryDB, projectRepo models.ProjectRepositoryDB) models.ConsumerUsecase {
	return &consumerUsecase{profileRepo, projectRepo}
}

var validate = validator.New()

func (u consumerUsecase) CreateProfile(e events.UserProfile) error {

	if err := validate.Struct(e); err != nil {
		logs.Warn("Invalid event", zap.Error(err))
		return err
	}
	profile := new(models.UserProfile)
	profile.ID = e.UserId
	profile.Name = e.Name
	profile.ProfileId = e.ProfileId

	err := u.profileRepo.Create(profile)
	if err != nil {
		logs.Error(fmt.Sprintf("Can't create profile with userid %d", e.UserId))
		return err
	}

	return nil
}

func (u consumerUsecase) UpdateProfile(e events.UserProfile) error {
	profile := new(models.UserProfile)
	if err := validate.Struct(e); err != nil {
		logs.Warn("Invalid event", zap.Error(err))
		return err
	}
	profile.ID = e.UserId
	profile.Name = e.Name
	profile.ProfileId = e.ProfileId
	err := u.profileRepo.Update(profile)
	if err != nil {
		logs.Error(fmt.Sprintf("Can't Update profile with userid %d", e.UserId))
		return err
	}

	return nil
}

func (u consumerUsecase) DeleteProfile(e events.UserProfileDeleted) error {
	err := u.profileRepo.Delete(e.UserId)
	if err != nil {
		logs.Error(fmt.Sprintf("Can't Update profile with userid %d", e.UserId))
		return err
	}

	return nil
}

// func (u consumerUsecase) CheckOffsetMessage(topic string, offset int64, partition int32) error {
// 	ConsumeOffset := new(models.ConsumerOffset)
// 	ConsumeOffset.Topic = topic
// 	ConsumeOffset.Offset = offset
// 	ConsumeOffset.Partition = partition

// 	check, err := u.consumerRepo.Get(ConsumeOffset)
// 	if err != nil {
// 		logs.Info(fmt.Sprintf("Save offset for topric %s partition %d offset %d", topic, offset, partition))
// 		u.consumerRepo.Create(ConsumeOffset)
// 		return nil
// 	}
// 	if check != nil || err == nil {
// 		return fmt.Errorf("have offset")
// 	}
// 	return fmt.Errorf("have offset")
// }

func (u consumerUsecase) CreateProject(e events.ProjectEvent) (err error) {

	project := new(models.Project)
	project.ID = e.ID
	project.Project, err = json.Marshal(e)
	if err != nil {
		logs.Error(fmt.Sprintf("Can't create project with userid %d", e.ID))
		return err
	}
	err = u.projectRepo.Create(project)
	if err != nil {
		logs.Error(fmt.Sprintf("Can't create project with userid %d", e.ID))
		return err
	}
	return nil
}

func (u consumerUsecase) UpdateProject(e events.ProjectEvent) (err error) {
	project := new(models.Project)
	project.ID = e.ID
	project.Project, err = json.Marshal(e)
	if err != nil {
		logs.Error(fmt.Sprintf("Can't create project with userid %d", e.ID))
		return err
	}
	err = u.projectRepo.Update(project)
	if err != nil {
		logs.Error(fmt.Sprintf("Can't update project with userid %d", e.ID))
		return err
	}
	return nil
}

func (u consumerUsecase) DeleteProject(e events.ProjectEventDeleted) error {
	err := u.projectRepo.Delete(e.ID)
	if err != nil {
		return err
	}
	return nil
}
