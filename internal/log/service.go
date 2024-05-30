package log

import (
	"trellode-go/internal/models"

	"gorm.io/gorm"
)

type LogServiceInterface interface {
	GetLogs(models.Context, int) ([]*models.Log, int, error)
	CreateLog(models.Context, *gorm.DB, *models.Log) (uint, int, error)
}

type LogService struct {
	repo LogRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewLogService(repo LogRepositoryInterface) LogService {
	return LogService{
		repo: repo,
	}
}

func (s LogService) GetLogs(context models.Context, boardId int) ([]*models.Log, int, error) {
	return s.repo.GetLogs(context, boardId)
}

func (s LogService) CreateLog(context models.Context, tx *gorm.DB, log *models.Log) (int, int, error) {
	return s.repo.CreateLog(context, tx, log)
}
