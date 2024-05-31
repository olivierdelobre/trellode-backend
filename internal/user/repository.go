package user

import (
	"errors"
	"net/http"
	"trellode-go/internal/models"

	"trellode-go/internal/utils/tools"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
	user.ID = uuid.NewString()
	user.PasswordHash = tools.HashPassword(user.PasswordHash)

	err := repo.db.Create(&user).Error
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusCreated, nil
}

func (repo UserRepository) Authenticate(context models.Context, user *models.User) (*models.User, int, error) {
	providedPasswordHash := tools.HashPassword(user.Password)
	var dbUser models.User
	err := repo.db.Where("email = ?", user.Email).First(&dbUser).Error
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	if dbUser.PasswordHash != providedPasswordHash {
		return nil, http.StatusForbidden, errors.New("invalid credentials")
	}
	return &dbUser, http.StatusOK, nil
}
