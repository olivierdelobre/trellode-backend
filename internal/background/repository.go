package background

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"
	"trellode-go/internal/log"
	"trellode-go/internal/models"
	"trellode-go/internal/utils/messages"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BackgroundRepository struct {
	db         *gorm.DB
	log        *zap.Logger
	logService log.LogService
}

type BackgroundRepositoryInterface interface {
	GetBackground(models.Context, int) (*models.Background, int, error)
	GetBackgrounds(models.Context) ([]*models.Background, int, error)
	CreateBackground(models.Context, string) (string, int, error)
	DeleteBackground(models.Context, int) (int, error)
}

func NewBackgroundRepository(db *gorm.DB, log *zap.Logger, logService log.LogService) BackgroundRepository {
	return BackgroundRepository{
		db:         db,
		log:        log,
		logService: logService,
	}
}

func (repo BackgroundRepository) GetBackground(context models.Context, id int) (*models.Background, int, error) {
	var background *models.Background
	err := repo.db.
		Where("id = ?", id).
		First(&background).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	//base64String := base64.StdEncoding.EncodeToString(background.Data)
	//background.DataBase64 = base64String

	return background, http.StatusOK, nil
}

func (repo BackgroundRepository) GetBackgrounds(context models.Context) ([]*models.Background, int, error) {
	backgrounds := []*models.Background{}
	err := repo.db.
		Where("user_id = ?", context.UserId).
		Find(&backgrounds).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	//for _, background := range backgrounds {
	//	base64String := base64.StdEncoding.EncodeToString(background.Data)
	//	background.DataBase64 = base64String
	//}

	return backgrounds, http.StatusOK, nil
}

func (repo BackgroundRepository) CreateBackground(context models.Context, base64Data string) (string, int, error) {
	background := models.Background{}
	// override userId
	background.ID = uuid.NewString()
	background.UserID = context.UserId
	background.Data = base64Data

	// remove header
	// find index of "base64,"
	index := strings.Index(background.Data, "base64,")
	header := background.Data[:index+7]
	dataNoHeader := background.Data[index+7:]

	decoded, err := base64.StdEncoding.DecodeString(dataNoHeader)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// calculate dominant color of an image
	img, _, err := image.Decode(bytes.NewReader(decoded))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// resize
	newWidth := 1920
	newHeight := 0 // Set to 0 to preserve aspect ratio
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	averageColor := averageColor(resizedImg)
	r, g, b, _ := averageColor.RGBA()
	colorCss := fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
	background.Color = colorCss

	// reencode to base64
	var buf bytes.Buffer
	if strings.Contains(header, "jpeg") {
		err = jpeg.Encode(&buf, resizedImg, nil)
	}
	if strings.Contains(header, "png") {
		err = png.Encode(&buf, resizedImg)
	}
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	background.Data = header + base64Str

	tx := repo.db.Begin()

	err = tx.Create(&background).Error
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// log operation
	_, severity, err := repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        "", // not related to a specific board
		Action:         "createbackground",
		ActionTargetID: background.ID,
	})
	if err != nil {
		tx.Rollback()
		return "", severity, err
	}

	tx.Commit()

	return background.ID, http.StatusCreated, nil
}

func (repo BackgroundRepository) DeleteBackground(context models.Context, id int) (int, error) {
	background, severity, err := repo.GetBackground(context, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return severity, err
	}
	if background.ID == "" {
		return http.StatusNotFound, errors.New(messages.GetMessage(context.Lang, "BackgroundNotFound"))
	}

	// check not used in any board
	err = repo.db.Where("background_id = ?", id).First(&models.Board{}).Error
	if err == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusForbidden, errors.New(messages.GetMessage(context.Lang, "BackgroundUsedInBoard"))
	}

	tx := repo.db.Begin()

	// delete background
	err = tx.Where("id = ?", id).Delete(&models.Background{}).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// log operation
	_, severity, err = repo.logService.CreateLog(context, tx, &models.Log{
		UserID:         context.UserId,
		BoardID:        "", // not related to a specific board
		Action:         "deletebackground",
		ActionTargetID: background.ID,
	})
	if err != nil {
		tx.Rollback()
		return severity, err
	}

	tx.Commit()

	return http.StatusAccepted, nil
}

func dominantColor(img image.Image) color.Color {
	colorCount := make(map[color.Color]int)
	bounds := img.Bounds()

	// Count colors in the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			colorCount[c]++
		}
	}

	// Find the most frequent color
	var dominant color.Color
	maxCount := 0
	for c, count := range colorCount {
		if count > maxCount {
			dominant = c
			maxCount = count
		}
	}

	return dominant
}

func averageColor(img image.Image) color.Color {
	bounds := img.Bounds()
	var rTotal, gTotal, bTotal, count uint32

	// Iterate over each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rTotal += r >> 8
			gTotal += g >> 8
			bTotal += b >> 8
			count++
		}
	}

	// Calculate average values
	avgR := uint8(rTotal / count)
	avgG := uint8(gTotal / count)
	avgB := uint8(bTotal / count)

	return color.RGBA{R: avgR, G: avgG, B: avgB, A: 255}
}
