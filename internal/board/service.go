package board

import (
	"errors"
	"net/http"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"
)

type BoardServiceInterface interface {
	GetBoards(models.Context, bool) ([]*models.Board, int, error)
	CreateBoard(models.Context, *models.Board) (uint, int, error)
	UpdateBoard(models.Context, *models.Board) (int, error)
	DeleteBoard(models.Context, int) (int, error)
}

type BoardService struct {
	repo BoardRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewBoardService(repo BoardRepositoryInterface) BoardService {
	return BoardService{
		repo: repo,
	}
}

func (s BoardService) GetBoard(context models.Context, id int) (*models.Board, int, error) {
	return s.repo.GetBoard(context, id)
}

func (s BoardService) GetBoards(context models.Context, archived bool) ([]*models.Board, int, error) {
	return s.repo.GetBoards(context, archived)
}

func (s BoardService) CreateBoard(context models.Context, board *models.Board) (int, int, error) {
	return s.repo.CreateBoard(context, board)
}

func (s BoardService) UpdateBoard(context models.Context, id int, board *models.Board) (int, error) {
	// check board exists
	existingBoard, severity, err := s.GetBoard(context, id)
	if err != nil {
		return severity, err
	}
	if existingBoard.ID == 0 {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "BoardNotFound"))
	}
	if existingBoard.UserID != context.UserId {
		return http.StatusForbidden, errors.New(messages.GetMessage(context.Lang, "Forbidden"))
	}

	return s.repo.UpdateBoard(context, board)
}

func (s BoardService) DeleteBoard(context models.Context, id int) (int, error) {
	// check board exists
	board, severity, err := s.GetBoard(context, id)
	if err != nil {
		return severity, err
	}
	if board.ID == 0 {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "BoardNotFound"))
	}
	if board.UserID != context.UserId {
		return http.StatusForbidden, errors.New(messages.GetMessage(context.Lang, "Forbidden"))
	}

	return s.repo.DeleteBoard(context, id)
}
