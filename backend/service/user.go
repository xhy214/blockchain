package service

import (
	"errors"
	"time"

	"blockchain/backend/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserService struct{}

var JWTSecret string
var JWTExpireHours int

func (s *UserService) Register(username, password, realName string) (*model.User, error) {
	existing, _ := model.FindUserByUsername(username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	user := &model.User{
		ID:       uuid.New().String(),
		Username: username,
		RealName: realName,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	if err := model.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(username, password string) (string, *model.User, error) {
	user, err := model.FindUserByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}
	if !user.CheckPassword(password) {
		return "", nil, errors.New("invalid username or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Duration(JWTExpireHours) * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", nil, err
	}
	return tokenStr, user, nil
}
