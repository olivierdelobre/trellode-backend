package log

import (
	"trellode-go/internal/models"

	"gorm.io/gorm"
)

type LogServiceInterface interface {
	GetLogs(models.Context, string) ([]*models.Log, int, error)
	CreateLog(models.Context, *gorm.DB, *models.Log) (string, int, error)
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

func (s LogService) GetLogs(context models.Context, boardId string) ([]*models.Log, int, error) {
	return s.repo.GetLogs(context, boardId)
}

func (s LogService) CreateLog(context models.Context, tx *gorm.DB, log *models.Log) (string, int, error) {
	return s.repo.CreateLog(context, tx, log)
}
