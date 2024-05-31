package list

import "trellode-go/internal/models"

type ListServiceInterface interface {
	GetList(models.Context, string) (*models.List, int, error)
	CreateList(models.Context, *models.List) (string, int, error)
	UpdateList(models.Context, *models.List) (int, error)
	UpdateCardsOrder(models.Context, string, string) (int, error)
	DeleteList(models.Context, string) (int, error)
}

type ListService struct {
	repo ListRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewListService(repo ListRepositoryInterface) ListService {
	return ListService{
		repo: repo,
	}
}

func (p ListService) GetList(context models.Context, id string) (*models.List, int, error) {
	return p.repo.GetList(context, id)
}

func (p ListService) CreateList(context models.Context, board *models.List) (string, int, error) {
	return p.repo.CreateList(context, board)
}

func (p ListService) UpdateList(context models.Context, board *models.List) (int, error) {
	return p.repo.UpdateList(context, board)
}

func (p ListService) UpdateCardsOrder(context models.Context, listId string, idsOrdered string) (int, error) {
	return p.repo.UpdateCardsOrder(context, listId, idsOrdered)
}

func (p ListService) DeleteList(context models.Context, id string) (int, error) {
	return p.repo.DeleteList(context, id)
}
