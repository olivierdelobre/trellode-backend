package background

import (
	"trellode-go/internal/models"
)

type BackgroundServiceInterface interface {
	GetBackground(models.Context, int) (*models.Background, int, error)
	GetBackgrounds(models.Context) ([]*models.Background, int, error)
	CreateBackground(models.Context, []byte) (uint, int, error)
	UpdateBackground(models.Context, *models.Background) (int, error)
	DeleteBackground(models.Context, int) (int, error)
}

type BackgroundService struct {
	repo BackgroundRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewBackgroundService(repo BackgroundRepositoryInterface) BackgroundService {
	return BackgroundService{
		repo: repo,
	}
}

func (s BackgroundService) GetBackground(context models.Context, id int) (*models.Background, int, error) {
	return s.repo.GetBackground(context, id)
}

func (s BackgroundService) GetBackgrounds(context models.Context) ([]*models.Background, int, error) {
	return s.repo.GetBackgrounds(context)
}

// func (s BackgroundService) CreateBackground(context models.Context, data []byte) (int, int, error) {
func (s BackgroundService) CreateBackground(context models.Context, data string) (int, int, error) {
	return s.repo.CreateBackground(context, data)
}

func (s BackgroundService) DeleteBackground(context models.Context, id int) (int, error) {
	return s.repo.DeleteBackground(context, id)
}
