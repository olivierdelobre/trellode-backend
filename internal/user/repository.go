package user

import (
	"errors"
	"net/http"
	"regexp"
	"trellode-go/internal/models"

	"trellode-go/internal/utils/messages"
	"trellode-go/internal/utils/tools"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

type UserRepositoryInterface interface {
	RegisterUser(models.Context, *models.User) (*models.User, int, error)
	Authenticate(models.Context, *models.User) (*models.User, int, error)
}

func NewUserRepository(db *gorm.DB, log *zap.Logger) UserRepository {
	return UserRepository{
		db:  db,
		log: log,
	}
}

func (repo UserRepository) RegisterUser(context models.Context, user *models.User) (*models.User, int, error) {
	// check password strength
	if len(user.Password) < 8 {
		return nil, http.StatusBadRequest, errors.New(messages.GetMessage(context.Lang, "PasswordTooShort"))
	}
	// check password contains uppercase and lowercase letters
	matchLowerCase := regexp.MustCompile(`[a-z]`)
	matchUpperCase := regexp.MustCompile(`[A-Z]`)
	matchFigure := regexp.MustCompile(`[0-9]`)
	matchSpecial := regexp.MustCompile(`[!@#$%^&*\(\)\+\-=\[\]{};':"\\|,\.<>\/\?]`)
	if !matchLowerCase.MatchString(user.Password) || !matchUpperCase.MatchString(user.Password) || !matchFigure.MatchString(user.Password) || !matchSpecial.MatchString(user.Password) {
		return nil, http.StatusBadRequest, errors.New(messages.GetMessage(context.Lang, "PasswordNotMatchingPolicies"))
	}

	// check if email already exists
	var dbUser models.User
	err := repo.db.Where("email = ?", user.Email).First(&dbUser).Error
	if err == nil {
		return nil, http.StatusBadRequest, errors.New(messages.GetMessage(context.Lang, "EmailAlreadyRegistered"))
	}

	user.ID = uuid.NewString()
	user.PasswordHash = tools.HashPassword(user.Password)

	err = repo.db.Create(&user).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusCreated, nil
}

func (repo UserRepository) Authenticate(context models.Context, user *models.User) (*models.User, int, error) {
	var dbUser models.User
	err := repo.db.Where("email = ?", user.Email).First(&dbUser).Error
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	errMatch := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.Password))
	if errMatch != nil {
		return nil, http.StatusForbidden, errors.New(messages.GetMessage(context.Lang, "InvalidCredentials"))
	}
	return &dbUser, http.StatusOK, nil
}
