package checklist

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"trellode-go/internal/log"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChecklistRepository struct {
	db         *gorm.DB
	log        *zap.Logger
	logService log.LogService
}

type ChecklistRepositoryInterface interface {
	GetChecklist(models.Context, string) (*models.Checklist, int, error)
	CreateChecklist(models.Context, *models.Checklist) (string, int, error)
	UpdateChecklist(models.Context, *models.Checklist) (int, error)
	DeleteChecklist(models.Context, string) (int, error)

	GetChecklistItem(models.Context, string) (*models.ChecklistItem, int, error)
	CreateChecklistItem(models.Context, *models.ChecklistItem) (string, int, error)
	UpdateChecklistItem(models.Context, *models.ChecklistItem) (int, error)
	DeleteChecklistItem(models.Context, string) (int, error)
}

func NewChecklistRepository(db *gorm.DB, log *zap.Logger, logService log.LogService) ChecklistRepository {
	return ChecklistRepository{
		db:         db,
		log:        log,
		logService: logService,
	}
}

func (repo ChecklistRepository) GetChecklist(context models.Context, id string) (*models.Checklist, int, error) {
	var checklist *models.Checklist
	err := repo.db.
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("id = ?", id).
		First(&checklist).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, err
	}
	if checklist.ID == "" {
		return nil, http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ChecklistNotFound"))
	}

	return checklist, http.StatusOK, nil
}

func (repo ChecklistRepository) CreateChecklist(context models.Context, checklist *models.Checklist) (string, int, error) {
	checklist.ID = uuid.NewString()
	checklist.ArchivedAt = nil

	tx := repo.db.Begin()

	err := tx.Create(&checklist).Error
	if err != nil {
		tx.Rollback()
		return "", http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfChecklist(checklist)
	if boardId == "" || err != nil {
		tx.Rollback()
		return "", http.StatusInternalServerError, err
	}
	_, severity, err := repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "createchecklist",
		ActionTargetID: checklist.ID,
	})
	if err != nil {
		tx.Rollback()
		return "", severity, err
	}

	tx.Commit()

	return checklist.ID, http.StatusCreated, nil
}

func (repo ChecklistRepository) UpdateChecklist(context models.Context, checklist *models.Checklist) (int, error) {
	// get checklist from db
	checklistBefore, severity, err := repo.GetChecklist(context, checklist.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if checklistBefore.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ChecklistNotFound"))
	}

	checklist.UpdatedAt = time.Now()
	// if board.ArchivedAt equals epoch 0, nullify archivedAt
	epoch0 := time.Unix(0, 0)
	if checklist.ArchivedAt != nil && checklist.ArchivedAt.Format("2006-01-02") == epoch0.Format("2006-01-02") {
		checklist.ArchivedAt = nil
	}

	// what changed?
	changes, err := whatChanged(checklistBefore, checklist)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// marshal changes to JSON string
	changesJson, err := json.Marshal(changes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := repo.db.Begin()

	err = tx.Omit("Comments", "ListID", "CreatedAt").Save(&checklist).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	operation := "updatechecklist"
	if checklistBefore.ArchivedAt == nil && checklist.ArchivedAt != nil {
		operation = "archivechecklist"
	}
	if checklistBefore.ArchivedAt != nil && checklist.ArchivedAt == nil {
		operation = "restorechecklist"
	}
	boardId, err := repo.getBoardIdOfChecklist(checklist)
	if boardId == "" || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         operation,
		ActionTargetID: checklist.ID,
		Changes:        string(changesJson),
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo ChecklistRepository) DeleteChecklist(context models.Context, id string) (int, error) {
	checklist, severity, err := repo.GetChecklist(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if checklist.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ChecklistNotFound"))
	}

	tx := repo.db.Begin()

	// remove items
	for _, item := range checklist.Items {
		err = tx.Delete(&item).Error
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}
	// remove checklist
	err = tx.Delete(&checklist).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfChecklist(checklist)
	if boardId == "" || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "deletechecklist",
		ActionTargetID: checklist.ID,
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo ChecklistRepository) GetChecklistItem(context models.Context, id string) (*models.ChecklistItem, int, error) {
	var checklistItem *models.ChecklistItem
	err := repo.db.
		Where("id = ?", id).
		First(&checklistItem).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	return checklistItem, http.StatusOK, nil
}

func (repo ChecklistRepository) CreateChecklistItem(context models.Context, checklistItem *models.ChecklistItem) (string, int, error) {
	checklistItem.ID = uuid.NewString()

	tx := repo.db.Begin()

	err := tx.Create(&checklistItem).Error
	if err != nil {
		tx.Rollback()
		return "", http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfChecklistItem(checklistItem)
	if boardId == "" || err != nil {
		tx.Rollback()
		return "", http.StatusInternalServerError, err
	}
	_, severity, err := repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "createchecklistitem",
		ActionTargetID: checklistItem.ID,
	})
	if err != nil {
		tx.Rollback()
		return "", severity, err
	}

	tx.Commit()

	return checklistItem.ID, http.StatusCreated, nil
}

func (repo ChecklistRepository) UpdateChecklistItem(context models.Context, checklistItem *models.ChecklistItem) (int, error) {
	// get checklistItem from db
	checklistItemBefore, severity, err := repo.GetChecklistItem(context, checklistItem.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if checklistItemBefore.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ChecklistItemNotFound"))
	}

	checklistItem.UpdatedAt = time.Now()

	// what changed?
	changes, err := whatChangedItem(checklistItemBefore, checklistItem)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// marshal changes to JSON string
	changesJson, err := json.Marshal(changes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := repo.db.Begin()

	err = tx.Omit("CreatedAt").Save(&checklistItem).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfChecklistItem(checklistItem)
	if boardId == "" || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "updatechecklistitem",
		ActionTargetID: checklistItem.ID,
		Changes:        string(changesJson),
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo ChecklistRepository) DeleteChecklistItem(context models.Context, id string) (int, error) {
	checklistItem, severity, err := repo.GetChecklistItem(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if checklistItem.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "ChecklistItemNotFound"))
	}

	tx := repo.db.Begin()

	err = tx.Delete(&checklistItem).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfChecklistItem(checklistItem)
	if boardId == "" || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "deletechecklist",
		ActionTargetID: checklistItem.ID,
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo ChecklistRepository) getBoardIdOfChecklist(checklist *models.Checklist) (string, error) {
	var card *models.Card
	err := repo.db.
		Where("id = ?", checklist.CardID).
		First(&card).Error
	if err != nil {
		return "", err
	}
	// get list
	var list *models.List
	err = repo.db.
		Where("id = ?", card.ListID).
		First(&list).Error
	if err != nil {
		return "", err
	}

	return list.BoardID, nil
}

func (repo ChecklistRepository) getBoardIdOfChecklistItem(checklistItem *models.ChecklistItem) (string, error) {
	var checklist *models.Checklist
	err := repo.db.
		Where("id = ?", checklistItem.ChecklistID).
		First(&checklist).Error
	if err != nil {
		return "", err
	}
	// get card
	var card *models.Card
	err = repo.db.
		Where("id = ?", checklist.CardID).
		First(&card).Error
	if err != nil {
		return "", err
	}
	// get list
	var list *models.List
	err = repo.db.
		Where("id = ?", card.ListID).
		First(&list).Error
	if err != nil {
		return "", err
	}

	return list.BoardID, nil
}

// whatChanged calculates the changes between two Checklist models and returns an array of LogChange models.
//
// Parameters:
// - before: a pointer to the previous Checklist model.
// - after: a pointer to the updated Checklist model.
//
// Returns:
// - []*models.LogChange: an array of LogChange models representing the changes between the two Checklist models.
// - error: an error if any occurred during the calculation.
func whatChanged(before *models.Checklist, after *models.Checklist) ([]*models.LogChange, error) {
	changes := []*models.LogChange{}

	if before.Title != after.Title {
		changes = append(changes, &models.LogChange{
			Field:     "title",
			FromValue: before.Title,
			ToValue:   after.Title,
		})
	}

	return changes, nil
}

func whatChangedItem(before *models.ChecklistItem, after *models.ChecklistItem) ([]*models.LogChange, error) {
	changes := []*models.LogChange{}

	if before.Title != after.Title {
		changes = append(changes, &models.LogChange{
			Field:     "title",
			FromValue: before.Title,
			ToValue:   after.Title,
		})
	}

	return changes, nil
}
