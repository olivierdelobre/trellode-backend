package list

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"trellode-go/internal/log"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ListRepository struct {
	db         *gorm.DB
	log        *zap.Logger
	logService log.LogService
}

type ListRepositoryInterface interface {
	GetList(models.Context, string) (*models.List, int, error)
	CreateList(models.Context, *models.List) (string, int, error)
	UpdateList(models.Context, *models.List) (int, error)
	UpdateCardsOrder(models.Context, string, string) (int, error)
	DeleteList(models.Context, string) (int, error)
}

func NewListRepository(db *gorm.DB, log *zap.Logger, logService log.LogService) ListRepository {
	return ListRepository{
		db:         db,
		log:        log,
		logService: logService,
	}
}

func (repo ListRepository) GetList(context models.Context, id string) (*models.List, int, error) {
	var list *models.List
	err := repo.db.
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("archived_at IS NULL").Order("position ASC")
		}).
		Preload("Cards.Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("id = ?", id).
		First(&list).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, err
	}
	if list.ID == "" {
		return nil, http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ListNotFound"))
	}

	return list, http.StatusOK, nil
}

func (repo ListRepository) CreateList(context models.Context, list *models.List) (string, int, error) {
	// get lists of board to determine position of new list
	var board *models.Board
	err := repo.db.
		Preload("Lists", repo.db.Where("archived_at IS NULL")).
		Where("id = ?", list.BoardID).
		First(&board).Error
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if board.ID == "" {
		return "", http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "BoardNotFound"))
	}

	list.ID = uuid.NewString()
	list.ArchivedAt = nil
	list.Position = len(board.Lists) + 1

	tx := repo.db.Begin()

	err = tx.Create(&list).Error
	if err != nil {
		tx.Rollback()
		return "", http.StatusInternalServerError, err
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
		return "", severity, err
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
	if listBefore.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ListNotFound"))
	}

	// get lists of board to reassign positions (if something was archived or restored)
	var board *models.Board
	err = repo.db.
		Preload("Lists", repo.db.Where("archived_at IS NULL")).
		Where("id = ?", list.BoardID).
		First(&board).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if board.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "BoardNotFound"))
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

	// reassign positions
	newPosition := 0
	for _, loopList := range board.Lists {
		if loopList.ID == list.ID {
			loopList = *list
		}
		if loopList.ArchivedAt != nil {
			continue
		}
		newPosition++
		list.Position = newPosition
		err := tx.Model(&models.List{}).Where("id = ?", loopList.ID).Update("position", newPosition).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
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

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo ListRepository) UpdateCardsOrder(context models.Context, listId string, idsOrdered string) (int, error) {
	// get list from db
	list, severity, err := repo.GetList(context, listId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if list.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ListNotFound"))
	}

	tx := repo.db.Begin()

	idsOrderedSplit := strings.Split(idsOrdered, ",")

	for i, id := range idsOrderedSplit {
		var card *models.Card
		err := tx.Where("id = ?", id).First(&card).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
		if card.ID == "" {
			tx.Rollback()
			return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "CardNotFound"))
		}
		card.Position = i + 1
		err = tx.Save(&card).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}

	// log operation
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        list.BoardID,
		Action:         "reordercards",
		ActionTargetID: list.ID,
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo ListRepository) DeleteList(context models.Context, id string) (int, error) {
	list, severity, err := repo.GetList(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if list.ID == "" {
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

	// get board from db to update lists' positions
	boardId := list.BoardID
	var board models.Board
	err = tx.Where("id = ?", boardId).First(&board).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	if board.ID == "" {
		tx.Rollback()
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "BoardNotFound"))
	}

	// update positions of lists
	newPosition := 0
	for _, loopList := range board.Lists {
		// don't process removed record
		if loopList.ID == id {
			continue
		}
		newPosition++
		err := tx.Model(&models.List{}).Where("id = ?", loopList.ID).Update("position", newPosition).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
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
