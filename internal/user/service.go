package user

import (
	"net/http"
	"os"
	"time"
	"trellode-go/internal/models"

	"github.com/golang-jwt/jwt"
)

type UserServiceInterface interface {
	RegisterUser(models.Context, *models.User) (*models.User, int, error)
	Authenticate(models.Context, string, string) (string, string, int, error)
}

type UserClaims struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Profile   string `json:"profile"`
	jwt.StandardClaims
}

type UserService struct {
	repo UserRepositoryInterface
}

// NewPersonService returns a service to manipulate unit
func NewUserService(repo UserRepositoryInterface) UserService {
	return UserService{
		repo: repo,
	}
}

func (s UserService) RegisterUser(context models.Context, user *models.User) (*models.User, int, error) {
	return s.repo.RegisterUser(context, user)
}

func (s UserService) Authenticate(context models.Context, user *models.User) (string, string, int, error) {
	user, severity, err := s.repo.Authenticate(context, user)
	if err != nil {
		return "", "", severity, err
	}
	userClaims := UserClaims{
		Id:        user.ID,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Profile:   "user",
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		},
	}

	signedAccessToken, err := s.CreateAccessToken(userClaims)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	refreshClaims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	signedRefreshToken, err := s.GetRefreshToken(refreshClaims)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	return signedAccessToken, signedRefreshToken, http.StatusOK, nil
}

func (s UserService) CreateAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func (s UserService) GetRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func (s UserService) ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAccessToken.Claims.(*UserClaims)
}

func (s UserService) ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedRefreshToken.Claims.(*jwt.StandardClaims)
}
