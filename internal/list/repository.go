package list

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"trellode-go/internal/log"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ListRepository struct {
	db         *gorm.DB
	log        *zap.Logger
	logService log.LogService
}

type ListRepositoryInterface interface {
	GetList(models.Context, int) (*models.List, int, error)
	CreateList(models.Context, *models.List) (int, int, error)
	UpdateList(models.Context, *models.List) (int, error)
	DeleteList(models.Context, int) (int, error)
}

func NewListRepository(db *gorm.DB, log *zap.Logger, logService log.LogService) ListRepository {
	return ListRepository{
		db:         db,
		log:        log,
		logService: logService,
	}
}

func (repo ListRepository) GetList(context models.Context, id int) (*models.List, int, error) {
	var list *models.List
	err := repo.db.
		Preload("Cards").
		Preload("Cards.Comments").
		Where("id = ?", id).
		First(&list).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	return list, http.StatusOK, nil
}

func (repo ListRepository) CreateList(context models.Context, list *models.List) (int, int, error) {
	list.ArchivedAt = nil

	tx := repo.db.Begin()

	err := tx.Create(&list).Error
	if err != nil {
		tx.Rollback()
		return 0, http.StatusInternalServerError, err
	}

	// log operation
	_, severity, err := repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        list.BoardID,
		Action:         "createlist",
		ActionTargetID: list.ID,
	})
	if err != nil {
		tx.Rollback()
		return 0, severity, err
	}

	tx.Commit()

	return list.ID, http.StatusCreated, nil
}

func (repo ListRepository) UpdateList(context models.Context, list *models.List) (int, error) {
	// get list from db
	listBefore, severity, err := repo.GetList(context, list.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if listBefore.ID == 0 {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ListdNotFound"))
	}

	// what changed?
	changes, err := whatChanged(listBefore, list)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// marshal changes to JSON string
	changesJson, err := json.Marshal(changes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := repo.db.Begin()

	err = tx.Omit("BoardID", "Cards").Save(&list).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	operation := "updatelist"
	if listBefore.ArchivedAt == nil && list.ArchivedAt != nil {
		operation = "archivelist"
	}
	if listBefore.ArchivedAt != nil && list.ArchivedAt == nil {
		operation = "restorelist"
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        list.BoardID,
		Action:         operation,
		ActionTargetID: list.ID,
		Changes:        string(changesJson),
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	return http.StatusAccepted, nil
}

func (repo ListRepository) DeleteList(context models.Context, id int) (int, error) {
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

	// delete comments
	for _, card := range list.Cards {
		for _, comment := range card.Comments {
			err = tx.Delete(&comment).Error
			if err != nil {
				tx.Rollback()
				return http.StatusInternalServerError, err
			}
		}
	}
	// delete cards
	for _, card := range list.Cards {
		err = tx.Delete(&card).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}
	// delete list
	err = tx.Delete(&list).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        list.BoardID,
		Action:         "deletelist",
		ActionTargetID: list.ID,
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func whatChanged(listBefore *models.List, listAfter *models.List) ([]*models.LogChange, error) {
	changes := []*models.LogChange{}

	if listBefore.Title != listAfter.Title {
		changes = append(changes, &models.LogChange{
			Field:     "title",
			FromValue: listBefore.Title,
			ToValue:   listAfter.Title,
		})
	}
	if listBefore.Position != listAfter.Position {
		changes = append(changes, &models.LogChange{
			Field:     "position",
			FromValue: strconv.Itoa(listBefore.Position),
			ToValue:   strconv.Itoa(listAfter.Position),
		})
	}

	return changes, nil
}
