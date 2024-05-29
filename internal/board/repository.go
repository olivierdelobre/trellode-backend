package board

import (
	"errors"
	"fmt"
	"image/color"
	"net/http"
	"strconv"
	"strings"
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
	GetBoards(models.Context, bool) ([]*models.Board, int, error)
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
		Preload("Background").
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

	if board.Background != nil {
		//base64String := base64.StdEncoding.EncodeToString(board.Background.Data)
		//board.Background.DataBase64 = base64String

		// calculate colors from dominant color
		menuColorDark, err := darkenColor(board.Background.Color, 0.5)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		board.MenuColorDark = colorToCSS(menuColorDark)
		menuColorLight, err := darkenColor(board.Background.Color, 0.7)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		board.MenuColorLight = colorToCSS(menuColorLight)
		listColor, err := lightenColor(board.Background.Color, 0.5)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		board.ListColor = colorToCSS(listColor)
	}

	return board, http.StatusOK, nil
}

func (repo BoardRepository) GetBoards(context models.Context, archived bool) ([]*models.Board, int, error) {
	boards := []*models.Board{}

	sql := "user_id = ? AND archived_at IS NULL"
	if archived {
		sql = "user_id = ? AND archived_at IS NOT NULL"
	}
	err := repo.db.
		Preload("Background").
		Preload("Lists", repo.db.Where("archived_at IS NULL")).
		Preload("Lists.Cards", repo.db.Where("archived_at IS NULL")).
		Preload("Lists.Cards.Comments").
		Where(sql, context.UserId).
		Find(&boards).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	//for _, board := range boards {
	//	if board.Background != nil {
	//		base64String := base64.StdEncoding.EncodeToString(board.Background.Data)
	//		board.Background.DataBase64 = base64String
	//	}
	//}

	return boards, http.StatusOK, nil
}

func (repo BoardRepository) CreateBoard(context models.Context, board *models.Board) (int, int, error) {
	// override userId
	board.UserID = context.UserId
	board.ArchivedAt = nil

	err := repo.db.Omit("BackgroundID", "Background", "Lists").Create(&board).Error
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return board.ID, http.StatusCreated, nil
}

func (repo BoardRepository) UpdateBoard(context models.Context, board *models.Board) (int, error) {
	board.UpdatedAt = time.Now()
	// if board.ArchivedAt equals epoch 0, nullify archivedAt
	epoch0 := time.Unix(0, 0)
	if board.ArchivedAt != nil && board.ArchivedAt.Format("2006-01-02") == epoch0.Format("2006-01-02") {
		board.ArchivedAt = nil
	}
	err := repo.db.Omit("UserID", "Background", "Lists", "CreatedAt").Save(&board).Error
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

	tx := repo.db.Begin()

	now := time.Now()
	// set archivedAt to current time
	board.ArchivedAt = &now
	err = tx.Omit("UserID", "BackgroundID", "Background", "Lists", "CreatedAt", "UpdatedAt").Save(&board).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	/*
		lists := board.Lists
		for _, list := range lists {
			// set archivedAt to current time
			list.ArchivedAt = &now
			err = tx.Save(&list).Error
			if err != nil {
				tx.Rollback()
				return http.StatusInternalServerError, err
			}

			cards := list.Cards
			for _, card := range cards {
				// set archivedAt to current time
				card.ArchivedAt = &now
				err = tx.Save(&card).Error
				if err != nil {
					tx.Rollback()
					return http.StatusInternalServerError, err
				}
			}
		}
	*/

	tx.Commit()

	return http.StatusAccepted, nil
}

func darkenColor(colorCss string, factor float64) (color.Color, error) {
	c, err := parseHexColor(colorCss)
	if err != nil {
		return nil, err
	}
	if factor < 0 || factor > 1 {
		return nil, errors.New("factor must be between 0 and 1")
	}

	r, g, b, a := c.RGBA()
	var newR, newG, newB uint8
	// Convert to 8-bit values and apply the darkening factor
	newR = uint8(float64(uint8(r>>8)) * factor)
	newG = uint8(float64(uint8(g>>8)) * factor)
	newB = uint8(float64(uint8(b>>8)) * factor)
	// Return the darkened color with the original alpha value
	return color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)}, nil
}

func lightenColor(colorCss string, factor float64) (color.Color, error) {
	c, err := parseHexColor(colorCss)
	if err != nil {
		return nil, err
	}
	if factor < 0 || factor > 1 {
		panic("factor must be between 0 and 1")
	}

	r, g, b, a := c.RGBA()
	// Convert to 8-bit values
	red := uint8(r >> 8)
	green := uint8(g >> 8)
	blue := uint8(b >> 8)

	// Interpolate towards white
	lightenedR := uint8(float64(red) + (255-float64(red))*factor)
	lightenedG := uint8(float64(green) + (255-float64(green))*factor)
	lightenedB := uint8(float64(blue) + (255-float64(blue))*factor)

	// Return the lightened color with the original alpha value
	return color.RGBA{R: lightenedR, G: lightenedG, B: lightenedB, A: uint8(a >> 8)}, nil
}

func parseHexColor(s string) (color.Color, error) {
	// Remove the '#' if it exists
	s = strings.TrimPrefix(s, "#")

	var r, g, b, a uint8
	var err error

	// #RRGGBB
	r, err = parseHexByte(s[0:2])
	if err != nil {
		return nil, err
	}
	g, err = parseHexByte(s[2:4])
	if err != nil {
		return nil, err
	}
	b, err = parseHexByte(s[4:6])
	if err != nil {
		return nil, err
	}
	a = 255 // fully opaque

	return color.RGBA{R: r, G: g, B: b, A: a}, nil
}

func parseHexByte(s string) (uint8, error) {
	v, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0, err
	}
	return uint8(v), nil
}

func colorToCSS(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}
