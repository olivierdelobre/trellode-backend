package board

import (
	"errors"
	"net/http"
	"time"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BoardRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

type BoardRepositoryInterface interface {
	GetBoard(models.Context, int) (*models.Board, int, error)
	GetBoards(models.Context) ([]*models.Board, int, error)
	CreateBoard(models.Context, *models.Board) (int, int, error)
	UpdateBoard(models.Context, *models.Board) (int, error)
	ArchiveBoard(models.Context, int) (int, error)
}

func NewBoardRepository(db *gorm.DB, log *zap.Logger) BoardRepository {
	return BoardRepository{
		db:  db,
		log: log,
	}
}

func (repo BoardRepository) GetBoard(context models.Context, id int) (*models.Board, int, error) {
	var board *models.Board
	err := repo.db.
		Preload("Lists", repo.db.Where("archived_at IS NULL")).
		Preload("Lists.Cards", repo.db.Where("archived_at IS NULL")).
		Preload("Lists.Cards.Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comments.created_at DESC")
		}).
		Where("id = ?", id).
		First(&board).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return board, http.StatusOK, nil
}

func (repo BoardRepository) GetBoards(context models.Context) ([]*models.Board, int, error) {
	boards := []*models.Board{}
	err := repo.db.
		Preload("Lists", repo.db.Where("archived_at IS NULL")).
		Preload("Lists.Cards", repo.db.Where("archived_at IS NULL")).
		Preload("Lists.Cards.Comments").
		Where("user_id = ?", context.UserId).
		Find(&boards).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return boards, http.StatusOK, nil
}

func (repo BoardRepository) CreateBoard(context models.Context, board *models.Board) (int, int, error) {
	// override userId
	board.UserID = context.UserId

	err := repo.db.Create(&board).Error
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return board.ID, http.StatusCreated, nil
}

func (repo BoardRepository) UpdateBoard(context models.Context, board *models.Board) (int, error) {
	err := repo.db.Save(&board).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

func (repo BoardRepository) ArchiveBoard(context models.Context, id int) (int, error) {
	board, severity, err := repo.GetBoard(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if board.ID == 0 {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "BoardNotFound"))
	}

	// set archivedAt to current time
	now := time.Now()
	board.ArchivedAt = &now
	err = repo.db.Save(&board).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}
