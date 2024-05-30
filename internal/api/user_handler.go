package api

import (
	"trellode-go/internal/background"
	"trellode-go/internal/board"
	"trellode-go/internal/card"
	"trellode-go/internal/comment"
	"trellode-go/internal/list"
	internalLog "trellode-go/internal/log"
	"trellode-go/internal/models"
	"trellode-go/internal/user"
	"trellode-go/internal/utils/config"
	"trellode-go/internal/utils/logging"

	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type server struct {
	db                *gorm.DB
	i18nBundle        *i18n.Bundle
	Router            *gin.Engine
	Log               *zap.Logger
	userService       user.UserService
	boardService      board.BoardService
	listService       list.ListService
	cardService       card.CardService
	commentService    comment.CommentService
	backgroundService background.BackgroundService
	logService        internalLog.LogService
}

func NewServer(db *gorm.DB, router *gin.Engine, log *zap.Logger) *server {
	logService := internalLog.NewLogService(internalLog.NewLogRepository(db, log))
	userService := user.NewUserService(user.NewUserRepository(db, log))
	boardService := board.NewBoardService(board.NewBoardRepository(db, log, logService))
	listService := list.NewListService(list.NewListRepository(db, log, logService))
	cardService := card.NewCardService(card.NewCardRepository(db, log, logService))
	commentService := comment.NewCommentService(comment.NewCommentRepository(db, log, logService))
	backgroundService := background.NewBackgroundService(background.NewBackgroundRepository(db, log, logService))

	// i18n for error messages
	bundle := i18n.NewBundle(language.French)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// load translation files
	_, err := bundle.LoadMessageFile("i18n/fr.toml")
	if err != nil {
		_, err := bundle.LoadMessageFile("../../assets/i18n/fr.toml")
		if err != nil {
			log.Info("Could not load fr.toml file: " + err.Error())
		}
	}
	_, err2 := bundle.LoadMessageFile("i18n/en.toml")
	if err2 != nil {
		_, err2 := bundle.LoadMessageFile("../../assets/i18n/en.toml")
		if err2 != nil {
			log.Info("Could not load en.toml file: " + err.Error())
		}
	}

	return &server{db, bundle, router, log, userService, boardService, listService, cardService, commentService, backgroundService, logService}
}

// RegisterUser 	godoc
// @Summary 		Get list by ID or name
// @Tags    		trellode
// @Param   		id path string true "105179"
// @Accept  		json
// @Produce 		json
// @Success 		200  {object}  models.User
// @Router  		/v1/trellode/{id} [get]
func (s *server) registerUser(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err == nil {
		user, severity, err := s.userService.RegisterUser(context, &user)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (s *server) authenticate(c *gin.Context) {
	context, err := getContext(c)
	if err != nil {
		logging.LogError(s.Log, c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err == nil {
		accessToken, refreshToken, severity, err := s.userService.Authenticate(context, &user)
		if err != nil {
			logging.LogError(s.Log, c, err.Error())
			c.JSON(severity, gin.H{"detail": err.Error()})
			return
		}

		if accessToken != "" {
			c.JSON(http.StatusOK, gin.H{"accesstoken": accessToken, "refreshtoken": refreshToken})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (s *server) options(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func SetupServer() *server {
	c := config.GetConfig()
	// Get logger from config
	log := c.Log
	// Get db from config
	db := c.Db
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	s := NewServer(db, router, log)

	s.Routes()

	return s
}
