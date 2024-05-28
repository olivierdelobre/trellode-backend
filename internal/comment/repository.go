package comment

import (
	"net/http"
	"trellode-go/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

type CommentRepositoryInterface interface {
	GetComment(models.Context, int) (*models.Comment, int, error)
	GetComments(models.Context, int) ([]*models.Comment, int, error)
	CreateComment(models.Context, *models.Comment) (int, int, error)
	UpdateComment(models.Context, *models.Comment) (int, error)
	DeleteComment(models.Context, int) (int, error)
}

func NewCommentRepository(db *gorm.DB, log *zap.Logger) CommentRepository {
	return CommentRepository{
		db:  db,
		log: log,
	}
}

func (repo CommentRepository) GetComment(context models.Context, id int) (*models.Comment, int, error) {
	var comment *models.Comment
	err := repo.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
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
	comment.UserID = context.UserId
	err := repo.db.Create(&comment).Error
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return comment.ID, http.StatusCreated, nil
}

func (repo CommentRepository) UpdateComment(context models.Context, comment *models.Comment) (int, error) {
	err := repo.db.Save(&comment).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}

func (repo CommentRepository) DeleteComment(context models.Context, id int) (int, error) {
	err := repo.db.Where("id = ?", id).Delete(&models.Comment{}).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}
