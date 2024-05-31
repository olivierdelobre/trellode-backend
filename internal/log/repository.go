package log

import (
	"errors"
	"net/http"
	"strings"
	"trellode-go/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LogRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

type LogRepositoryInterface interface {
	GetLogs(models.Context, string) ([]*models.Log, int, error)
	CreateLog(models.Context, *gorm.DB, *models.Log) (string, int, error)
}

func NewLogRepository(db *gorm.DB, log *zap.Logger) LogRepository {
	return LogRepository{
		db:  db,
		log: log,
	}
}

func (repo LogRepository) GetLogs(context models.Context, boardId string) ([]*models.Log, int, error) {
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
			if list.ID != "" {
				log.ActionTargetTitle = list.Title
			}
		}
		if strings.HasSuffix(log.Action, "card") {
			var card *models.Card
			err := repo.db.Where("id = ?", log.ActionTargetID).First(&card).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, http.StatusInternalServerError, err
			}
			if card.ID != "" {
				log.ActionTargetTitle = card.Title
			}
		}
		if strings.HasSuffix(log.Action, "board") {
			var board *models.Board
			err := repo.db.Where("id = ?", log.ActionTargetID).First(&board).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, http.StatusInternalServerError, err
			}
			if board.ID != "" {
				log.ActionTargetTitle = board.Title
			}
		}
	}

	return logs, http.StatusOK, nil
}

func (repo LogRepository) CreateLog(context models.Context, tx *gorm.DB, log *models.Log) (string, int, error) {
	log.ID = uuid.NewString()
	// override userId
	log.UserID = context.UserId

	err := tx.Create(&log).Error
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	return log.ID, http.StatusCreated, nil
}
