package card

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
	"trellode-go/internal/log"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CardRepository struct {
	db         *gorm.DB
	log        *zap.Logger
	logService log.LogService
}

type CardRepositoryInterface interface {
	GetCard(models.Context, string) (*models.Card, int, error)
	CreateCard(models.Context, *models.Card) (string, int, error)
	UpdateCard(models.Context, *models.Card) (int, error)
	DeleteCard(models.Context, string) (int, error)
}

func NewCardRepository(db *gorm.DB, log *zap.Logger, logService log.LogService) CardRepository {
	return CardRepository{
		db:         db,
		log:        log,
		logService: logService,
	}
}

func (repo CardRepository) GetCard(context models.Context, id string) (*models.Card, int, error) {
	var card *models.Card
	err := repo.db.
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Checklists", func(db *gorm.DB) *gorm.DB {
			return db.Where("archived_at IS NULL").Order("title ASC")
		}).
		Preload("Checklists.Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ?", id).
		First(&card).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, err
	}
	if card.ID == "" {
		return nil, http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "CardNotFound"))
	}

	return card, http.StatusOK, nil
}

func (repo CardRepository) CreateCard(context models.Context, card *models.Card) (string, int, error) {
	// get cards of list to determine position of new card
	var list *models.List
	err := repo.db.
		Preload("Cards", repo.db.Where("archived_at IS NULL")).
		Where("id = ?", card.ListID).
		First(&list).Error
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if list.ID == "" {
		return "", http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ListNotFound"))
	}

	// generate UUID
	card.ID = uuid.NewString()
	card.ArchivedAt = nil
	card.Position = len(list.Cards) + 1

	tx := repo.db.Begin()

	err = tx.Create(&card).Error
	if err != nil {
		tx.Rollback()
		return "", http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfCard(card)
	if boardId == "" || err != nil {
		tx.Rollback()
		return "", http.StatusInternalServerError, err
	}
	_, severity, err := repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "createcard",
		ActionTargetID: card.ID,
	})
	if err != nil {
		tx.Rollback()
		return "", severity, err
	}

	tx.Commit()

	return card.ID, http.StatusCreated, nil
}

func (repo CardRepository) UpdateCard(context models.Context, card *models.Card) (int, error) {
	// get card from db
	cardBefore, severity, err := repo.GetCard(context, card.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if cardBefore.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "CardNotFound"))
	}

	card.UpdatedAt = time.Now()
	// if board.ArchivedAt equals epoch 0, nullify archivedAt
	epoch0 := time.Unix(0, 0)
	if card.ArchivedAt != nil && card.ArchivedAt.Format("2006-01-02") == epoch0.Format("2006-01-02") {
		card.ArchivedAt = nil
	}

	// what changed?
	changes, err := whatChanged(cardBefore, card)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// marshal changes to JSON string
	changesJson, err := json.Marshal(changes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := repo.db.Begin()

	err = tx.Omit("Comments", "ListID", "CreatedAt").Save(&card).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	operation := "updatecard"
	if cardBefore.ArchivedAt == nil && card.ArchivedAt != nil {
		operation = "archivecard"
	}
	if cardBefore.ArchivedAt != nil && card.ArchivedAt == nil {
		operation = "restorecard"
	}
	boardId, err := repo.getBoardIdOfCard(card)
	if boardId == "" || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         operation,
		ActionTargetID: card.ID,
		Changes:        string(changesJson),
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo CardRepository) DeleteCard(context models.Context, id string) (int, error) {
	card, severity, err := repo.GetCard(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if card.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ListNotFound"))
	}

	tx := repo.db.Begin()

	// remove comments
	for _, comment := range card.Comments {
		err = tx.Delete(&comment).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}
	// remove card
	err = tx.Delete(&card).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfCard(card)
	if boardId == "" || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "deletecard",
		ActionTargetID: card.ID,
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo CardRepository) getBoardIdOfCard(card *models.Card) (string, error) {
	var list *models.List
	err := repo.db.
		Where("id = ?", card.ListID).
		Find(&list).Error
	if err != nil {
		return "", err
	}

	return list.BoardID, nil
}

// whatChanged calculates the changes between two Card models and returns an array of LogChange models.
//
// Parameters:
// - cardBefore: a pointer to the previous Card model.
// - cardAfter: a pointer to the updated Card model.
//
// Returns:
// - []*models.LogChange: an array of LogChange models representing the changes between the two Card models.
// - error: an error if any occurred during the calculation.
func whatChanged(cardBefore *models.Card, cardAfter *models.Card) ([]*models.LogChange, error) {
	changes := []*models.LogChange{}

	if cardBefore.Title != cardAfter.Title {
		changes = append(changes, &models.LogChange{
			Field:     "title",
			FromValue: cardBefore.Title,
			ToValue:   cardAfter.Title,
		})
	}
	if cardBefore.Description != cardAfter.Description {
		changes = append(changes, &models.LogChange{
			Field:     "description",
			FromValue: cardBefore.Description,
			ToValue:   cardAfter.Description,
		})
	}
	if cardBefore.Position != cardAfter.Position {
		changes = append(changes, &models.LogChange{
			Field:     "position",
			FromValue: strconv.Itoa(cardBefore.Position),
			ToValue:   strconv.Itoa(cardAfter.Position),
		})
	}

	return changes, nil
}
