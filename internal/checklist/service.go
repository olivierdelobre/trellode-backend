package checklist

import "trellode-go/internal/models"

type ChecklistServiceInterface interface {
	GetChecklist(models.Context, string) (*models.Checklist, int, error)
	CreateChecklist(models.Context, *models.Checklist) (string, int, error)
	UpdateChecklist(models.Context, *models.Checklist) (int, error)
	DeleteChecklist(models.Context, string) (int, error)

	GetChecklistItem(models.Context, string) (*models.ChecklistItem, int, error)
	CreateChecklistItem(models.Context, *models.ChecklistItem) (string, int, error)
	UpdateChecklistItem(models.Context, *models.ChecklistItem) (int, error)
	DeleteChecklistItem(models.Context, string) (int, error)
}

type ChecklistService struct {
	repo ChecklistRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewChecklistService(repo ChecklistRepositoryInterface) ChecklistService {
	return ChecklistService{
		repo: repo,
	}
}

func (p ChecklistService) GetChecklist(context models.Context, id string) (*models.Checklist, int, error) {
	return p.repo.GetChecklist(context, id)
}

func (p ChecklistService) CreateChecklist(context models.Context, board *models.Checklist) (string, int, error) {
	return p.repo.CreateChecklist(context, board)
}

func (p ChecklistService) UpdateChecklist(context models.Context, board *models.Checklist) (int, error) {
	return p.repo.UpdateChecklist(context, board)
}

func (p ChecklistService) DeleteChecklist(context models.Context, id string) (int, error) {
	return p.repo.DeleteChecklist(context, id)
}

func (p ChecklistService) GetChecklistItem(context models.Context, id string) (*models.ChecklistItem, int, error) {
	return p.repo.GetChecklistItem(context, id)
}

func (p ChecklistService) CreateChecklistItem(context models.Context, board *models.ChecklistItem) (string, int, error) {
	return p.repo.CreateChecklistItem(context, board)
}

func (p ChecklistService) UpdateChecklistItem(context models.Context, board *models.ChecklistItem) (int, error) {
	return p.repo.UpdateChecklistItem(context, board)
}

func (p ChecklistService) DeleteChecklistItem(context models.Context, id string) (int, error) {
	return p.repo.DeleteChecklistItem(context, id)
}
