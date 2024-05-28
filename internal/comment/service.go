package comment

import "trellode-go/internal/models"

type CommentServiceInterface interface {
	GetComment(models.Context, int) (*models.Comment, int, error)
	GetComments(models.Context, int) ([]*models.Comment, int, error)
	CreateComment(models.Context, *models.Comment) (int, int, error)
	UpdateComment(models.Context, *models.Comment) (int, error)
	DeleteComment(models.Context, int) (int, error)
}

type CommentService struct {
	repo CommentRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewCommentService(repo CommentRepositoryInterface) CommentService {
	return CommentService{
		repo: repo,
	}
}

func (p CommentService) GetComment(context models.Context, id int) (*models.Comment, int, error) {
	return p.repo.GetComment(context, id)
}

func (p CommentService) GetComments(context models.Context, boardId int) ([]*models.Comment, int, error) {
	return p.repo.GetComments(context, boardId)
}

func (p CommentService) CreateComment(context models.Context, board *models.Comment) (int, int, error) {
	return p.repo.CreateComment(context, board)
}

func (p CommentService) UpdateComment(context models.Context, board *models.Comment) (int, error) {
	return p.repo.UpdateComment(context, board)
}

func (p CommentService) DeleteComment(context models.Context, id int) (int, error) {
	return p.repo.DeleteComment(context, id)
}
