package card

import "trellode-go/internal/models"

type CardServiceInterface interface {
	GetCard(models.Context, int) (*models.Card, int, error)
	CreateCard(models.Context, *models.Card) (int, int, error)
	UpdateCard(models.Context, *models.Card) (int, error)
	DeleteCard(models.Context, int) (int, error)
}

type CardService struct {
	repo CardRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewCardService(repo CardRepositoryInterface) CardService {
	return CardService{
		repo: repo,
	}
}

func (p CardService) GetCard(context models.Context, id int) (*models.Card, int, error) {
	return p.repo.GetCard(context, id)
}

func (p CardService) CreateCard(context models.Context, board *models.Card) (int, int, error) {
	return p.repo.CreateCard(context, board)
}

func (p CardService) UpdateCard(context models.Context, board *models.Card) (int, error) {
	return p.repo.UpdateCard(context, board)
}

func (p CardService) DeleteCard(context models.Context, id int) (int, error) {
	return p.repo.DeleteCard(context, id)
}
