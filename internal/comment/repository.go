package comment

import (
	"encoding/json"
	"errors"
	"net/http"
	"trellode-go/internal/log"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db         *gorm.DB
	log        *zap.Logger
	logService log.LogService
}

type CommentRepositoryInterface interface {
	GetComment(models.Context, int) (*models.Comment, int, error)
	GetComments(models.Context, int) ([]*models.Comment, int, error)
	CreateComment(models.Context, *models.Comment) (int, int, error)
	UpdateComment(models.Context, *models.Comment) (int, error)
	DeleteComment(models.Context, int) (int, error)
}

func NewCommentRepository(db *gorm.DB, log *zap.Logger, logService log.LogService) CommentRepository {
	return CommentRepository{
		db:         db,
		log:        log,
		logService: logService,
	}
}

func (repo CommentRepository) GetComment(context models.Context, id int) (*models.Comment, int, error) {
	var comment *models.Comment
	err := repo.db.Where("id = ?", id).First(&comment).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, err
	}

	return comment, http.StatusOK, nil
}

func (repo CommentRepository) GetComments(context models.Context, boardId int) ([]*models.Comment, int, error) {
	comments := []*models.Comment{}
	err := repo.db.Where("card_id = ?", boardId).Find(&comments).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return comments, http.StatusOK, nil
}

func (repo CommentRepository) CreateComment(context models.Context, comment *models.Comment) (int, int, error) {
	tx := repo.db.Begin()

	comment.UserID = context.UserId
	err := tx.Create(&comment).Error
	if err != nil {
		tx.Rollback()
		return 0, http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfComment(comment)
	if boardId == 0 || err != nil {
		tx.Rollback()
		return 0, http.StatusInternalServerError, err
	}
	_, severity, err := repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "createcomment",
		ActionTargetID: comment.ID,
	})
	if err != nil {
		tx.Rollback()
		return 0, severity, err
	}

	tx.Commit()

	return comment.ID, http.StatusCreated, nil
}

func (repo CommentRepository) UpdateComment(context models.Context, comment *models.Comment) (int, error) {
	// get comment from db
	commentBefore, severity, err := repo.GetComment(context, comment.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if commentBefore.ID == 0 {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "CommentNotFound"))
	}

	// what changed?
	changes, err := whatChanged(commentBefore, comment)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// marshal changes to JSON string
	changesJson, err := json.Marshal(changes)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	tx := repo.db.Begin()

	err = tx.Save(&comment).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfComment(comment)
	if boardId == 0 || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "updatecomment",
		ActionTargetID: comment.ID,
		Changes:        string(changesJson),
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo CommentRepository) DeleteComment(context models.Context, id int) (int, error) {
	// get comment from db
	commentBefore, severity, err := repo.GetComment(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if commentBefore.ID == 0 {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "CommentNotFound"))
	}

	tx := repo.db.Begin()

	err = tx.Where("id = ?", id).Delete(&models.Comment{}).Error
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	// log operation
	boardId, err := repo.getBoardIdOfComment(commentBefore)
	if boardId == 0 || err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        boardId,
		Action:         "deletecomment",
		ActionTargetID: id,
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func (repo CommentRepository) getBoardIdOfComment(comment *models.Comment) (int, error) {
	var card *models.Card
	err := repo.db.
		Where("id = ?", comment.CardID).
		First(&card).Error
	if err != nil {
		return 0, err
	}
	// get list
	var list *models.List
	err = repo.db.
		Where("id = ?", card.ListID).
		First(&list).Error
	if err != nil {
		return 0, err
	}

	return list.BoardID, nil
}

func whatChanged(commentBefore *models.Comment, commentAfter *models.Comment) ([]*models.LogChange, error) {
	changes := []*models.LogChange{}

	if commentBefore.Content != commentAfter.Content {
		changes = append(changes, &models.LogChange{
			Field:     "content",
			FromValue: commentBefore.Content,
			ToValue:   commentAfter.Content,
		})
	}

	return changes, nil
}
