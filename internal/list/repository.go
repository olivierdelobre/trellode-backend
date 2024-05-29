package list

import (
	"errors"
	"net/http"
	"time"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ListRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

type ListRepositoryInterface interface {
	GetList(models.Context, int) (*models.List, int, error)
	CreateList(models.Context, *models.List) (int, int, error)
	UpdateList(models.Context, *models.List) (int, error)
	ArchiveList(models.Context, int) (int, error)
}

func NewListRepository(db *gorm.DB, log *zap.Logger) ListRepository {
	return ListRepository{
		db:  db,
		log: log,
	}
}

func (repo ListRepository) GetList(context models.Context, id int) (*models.List, int, error) {
	var list *models.List
	err := repo.db.
		Preload("Cards").
		Preload("Cards.Comments").
		Where("id = ?", id).
		First(&list).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return list, http.StatusOK, nil
}

func (repo ListRepository) CreateList(context models.Context, list *models.List) (int, int, error) {
	err := repo.db.Create(&list).Error
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return list.ID, http.StatusCreated, nil
}

func (repo ListRepository) UpdateList(context models.Context, list *models.List) (int, error) {
	err := repo.db.Save(&list).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

func (repo ListRepository) ArchiveList(context models.Context, id int) (int, error) {
	list, severity, err := repo.GetList(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if list.ID == 0 {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ListNotFound"))
	}

	tx := repo.db.Begin()

	// set archivedAt to current time
	now := time.Now()
	list.ArchivedAt = &now
	err = tx.Save(&list).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	/*
		// archive all cards
		err = tx.Model(&models.Card{}).Where("list_id = ?", id).Update("archived_at", now).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	*/

	tx.Commit()

	return http.StatusAccepted, nil
}
