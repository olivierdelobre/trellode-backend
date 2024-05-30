package log

import (
	"errors"
	"net/http"
	"strings"
	"trellode-go/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LogRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

type LogRepositoryInterface interface {
	GetLogs(models.Context, int) ([]*models.Log, int, error)
	CreateLog(models.Context, *gorm.DB, *models.Log) (int, int, error)
}

func NewLogRepository(db *gorm.DB, log *zap.Logger) LogRepository {
	return LogRepository{
		db:  db,
		log: log,
	}
}

func (repo LogRepository) GetLogs(context models.Context, boardId int) ([]*models.Log, int, error) {
	logs := []*models.Log{}

	err := repo.db.
		Preload("User").
		Where("user_id = ? AND board_id = ?", context.UserId, boardId).
		Order("created_at DESC").
		Find(&logs).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// enrich with target object title
	for _, log := range logs {
		if strings.HasSuffix(log.Action, "list") {
			var list *models.List
			err := repo.db.Where("id = ?", log.ActionTargetID).First(&list).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, http.StatusInternalServerError, err
			}
			if list.ID != 0 {
				log.ActionTargetTitle = list.Title
			}
		}
		if strings.HasSuffix(log.Action, "card") {
			var card *models.Card
			err := repo.db.Where("id = ?", log.ActionTargetID).First(&card).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, http.StatusInternalServerError, err
			}
			if card.ID != 0 {
				log.ActionTargetTitle = card.Title
			}
		}
		if strings.HasSuffix(log.Action, "board") {
			var board *models.Board
			err := repo.db.Where("id = ?", log.ActionTargetID).First(&board).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, http.StatusInternalServerError, err
			}
			if board.ID != 0 {
				log.ActionTargetTitle = board.Title
			}
		}
	}

	return logs, http.StatusOK, nil
}

func (repo LogRepository) CreateLog(context models.Context, tx *gorm.DB, log *models.Log) (int, int, error) {
	// override userId
	log.UserID = context.UserId

	err := tx.Create(&log).Error
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return log.ID, http.StatusCreated, nil
}
