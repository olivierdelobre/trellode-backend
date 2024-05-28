package card

import (
	"net/http"
	"trellode-go/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CardRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

type CardRepositoryInterface interface {
	GetCard(models.Context, int) (*models.Card, int, error)
	CreateCard(models.Context, *models.Card) (int, int, error)
	UpdateCard(models.Context, *models.Card) (int, error)
	ArchiveCard(models.Context, int) (int, error)
}

func NewCardRepository(db *gorm.DB, log *zap.Logger) CardRepository {
	return CardRepository{
		db:  db,
		log: log,
	}
}

func (repo CardRepository) GetCard(context models.Context, id int) (*models.Card, int, error) {
	var card *models.Card
	err := repo.db.
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comments.created_at DESC")
		}).
		Where("id = ? AND archived_at IS NULL", id).
		Find(&card).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return card, http.StatusOK, nil
}

func (repo CardRepository) CreateCard(context models.Context, card *models.Card) (int, int, error) {
	err := repo.db.Create(&card).Error
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return card.ID, http.StatusCreated, nil
}

func (repo CardRepository) UpdateCard(context models.Context, card *models.Card) (int, error) {
	err := repo.db.Save(&card).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

func (repo CardRepository) ArchiveCard(context models.Context, id int) (int, error) {
	err := repo.db.Where("id = ?", id).Delete(&models.Card{}).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}