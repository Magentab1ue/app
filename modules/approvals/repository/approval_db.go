package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"approval-service/logs"
	"approval-service/modules/entities/models"
)

type approvalRepositoryDB struct {
	db *gorm.DB
}

func NewapprovalRepositoryDB(db *gorm.DB) models.ApprovalRepository {
	// Check if the table exists
	if !db.Migrator().HasTable(&models.Approvals{}) {
		err := db.AutoMigrate(models.Approvals{})
		if err != nil {
			panic(err)
		}
	}
	return approvalRepositoryDB{db: db}
}

func (r approvalRepositoryDB) Update(req *models.Approvals) (*models.Approvals, error) {

	//update to database
	if err := r.db.Save(&req).Error; err != nil {
		logs.Error(fmt.Sprintf("Error updating approval with request ID %d: %v", req.ID, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't updating approval with request ID %d", req.ID)
	}
	return req, nil
}

func (r approvalRepositoryDB) GetSendRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optional["sender_id"] = userId
	optionalStr := fmt.Sprintf("%v", optional)

	err := r.db.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approval Send request user with ID %d with optional %s : %v", userId, optionalStr, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't finding approval Send request user with ID %d with optional %s", userId, optionalStr)
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approval Send request user with ID %d with optional %s ", userId, optionalStr)
	}
	return approval, nil
}

func (r approvalRepositoryDB) GetReceiveRequest(userId uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db.Where("? = ANY(\"to\")", userId)

	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approval receive request user with ID %d with optinoal %s %v", userId, optionalStr, err), zap.Error(err))
		return nil, fmt.Errorf("error finding approval receive request user with ID %d with optinoal %s", userId, optionalStr)
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approval receive request user with ID %d with optinoal %s", userId, optionalStr)
	}
	return approval, nil
}

func (r approvalRepositoryDB) DeleteApproval(requestId uint) (*models.Approvals, error) {

	approval := new(models.Approvals)
	if err := r.db.First(&approval, requestId).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't finding approval with request ID %d: %v", requestId, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't finding approval with request ID %d: %v", requestId, err)
	}

	if err := r.db.Delete(&approval, requestId).Error; err != nil {
		logs.Error(fmt.Sprintf("Error cant't delete approval with request ID %d: %v", requestId, err), zap.Error(err))
		return nil, fmt.Errorf("error cant't delete approval with request ID %d: %v", requestId, err)
	}

	return approval, nil
}

func (r approvalRepositoryDB) Create(request *models.Approvals) (*models.Approvals, error) {
	if err := r.db.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return request, nil
}

func (r approvalRepositoryDB) GetByID(id uint) (*models.Approvals, error) {
	var approval models.Approvals
	result := r.db.First(&approval, id)

	if result.Error == nil && result.RowsAffected > 0 {
		return &approval, nil
	} else {
		return nil, errors.New("find position by ID not found")
	}
}

func (r approvalRepositoryDB) GetAll(optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db

	if _, ok := optional["status"]; ok {
		condition = condition.Where("status = ?", optional["status"])

	}
	if _, ok := optional["to"]; ok {
		to := optional["to"]
		delete(optional, "to")
		condition = condition.Where("? = ANY(\"to\")", to)
	}

	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals with optional %s : %v", optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals optional %s ", optionalStr)
	}
	return approval, nil

}

func (r approvalRepositoryDB) GetByUserID(id uint, optional map[string]interface{}) ([]models.Approvals, error) {

	approval := []models.Approvals{}
	optionalStr := fmt.Sprintf("%v", optional)
	condition := r.db
	optional["sender_id"] = id
	optional["to"] = id
	if _, ok := optional["to"]; ok {
		to := optional["to"]
		delete(optional, "to")
		condition = condition.Where("? = ANY(\"to\")", to)
	}

	err := condition.Where(optional).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals by userid %d with optional %s : %v", id, optionalStr, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals by userid %d optional %s ", id, optionalStr)
	}

	return approval, nil
}

func (r approvalRepositoryDB) GetByRequestID(id uuid.UUID) ([]models.Approvals, error) {
	approval := []models.Approvals{}

	err := r.db.Where("request_id = ?", id).Find(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals by request_id : %v  : %v", id, err), zap.Error(err))
		return nil, err
	}
	if len(approval) == 0 {
		return nil, fmt.Errorf("error finding approvals by request_id : %v ", id)
	}

	return approval, nil
}

func (r approvalRepositoryDB) GetByRequestIDLast(id uuid.UUID) (*models.Approvals, error) {
	approval := new(models.Approvals)

	err := r.db.Where("request_id = ?", id).Last(&approval).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding approvals by request_id : %v  : %v to validation", id, err), zap.Error(err))
		return nil, err
	}

	return approval, nil
}

func (r approvalRepositoryDB) GetProjectById(id uint) (*models.Project, error) {
	project := new(models.Project)

	err := r.db.Where("id = ?", id).Last(&project).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding project by id : %v  : %v to validation", id, err), zap.Error(err))
		return nil, fmt.Errorf("error finding project by id : %v  : to validation", id)
	}
	return project, nil
}
func (r approvalRepositoryDB) GetUserById(id uint) (*models.UserProfile, error) {
	userProfile := new(models.UserProfile)

	err := r.db.Where("id = ?", id).Last(&userProfile).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error finding userid by id : %v  : %v to validation", id, err), zap.Error(err))
		return nil, fmt.Errorf("error finding userid by id : %v  : to validation", id)
	}
	return userProfile, nil
}

func (r approvalRepositoryDB) CreateProject(request *models.Project) (*models.Project, error) {
	err := r.db.FirstOrCreate(request).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error Create Project with ID %d: %v", request.ID, err), zap.Error(err))
		return nil, fmt.Errorf("failed to create Project: %v", err)
	}

	return request, nil
}

func (r approvalRepositoryDB) CreateUser(request *models.UserProfile) (*models.UserProfile, error) {
	err := r.db.FirstOrCreate(request).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error Create Project with ID %d: %v", request.ID, err), zap.Error(err))
		return nil, fmt.Errorf("failed to create Project: %v", err)
	}
	return request, nil
}

func (r approvalRepositoryDB) GetListTaskCheck(ids []int64) ([]models.Task, error) {
	tasks := []models.Task{}
	err := r.db.Find(&tasks, ids).Where("status = ?", models.TaskAppproveStatusMap[1]).Or("status = ?", models.TaskAppproveStatusMap[3]).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	return tasks, nil
}

func (r approvalRepositoryDB) GetTask(id int64) (*models.Task, error) {
	task := new(models.Task)
	err := r.db.Find(&task, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tast id %d : %v", id, err)
	}
	return task, nil
}
func (r approvalRepositoryDB) GetTasks(ids []int64) ([]models.Task, error) {
	tasks := []models.Task{}
	err := r.db.Find(&tasks, ids).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	if len(tasks) != len(ids) {
		return nil, fmt.Errorf("no tasks found with the provided IDs")
	}
	return tasks, nil
}
func (r approvalRepositoryDB) UpdateTasksStatus(ids []int64, status string) error {
	err := r.db.Model(&models.Task{}).Where("id IN ?", ids).Updates(map[string]interface{}{"approval_status": status}).Error
	if err != nil {
		return fmt.Errorf("failed to update tasks: %v", err)
	}
	return nil
}

func (r approvalRepositoryDB) CreateTask(request *models.Task) (*models.Task, error) {
	err := r.db.FirstOrCreate(request).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Error Create Task with ID %d: %v", request.ID, err), zap.Error(err))
		return nil, fmt.Errorf("failed to create Task: %v", err)
	}

	return request, nil
}
func (r approvalRepositoryDB) GetAllProject() ([]models.Project, error) {
	project := []models.Project{}
	err := r.db.Find(&project).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	return project, nil
}

func (r approvalRepositoryDB) GetAllTask() ([]models.Task, error) {
	tasks := []models.Task{}
	err := r.db.Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	return tasks, nil
}
