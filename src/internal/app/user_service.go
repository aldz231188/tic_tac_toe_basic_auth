package app

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"strings"
	"t03/internal/api/dto"
	"t03/internal/domain"
)

type UserServiceImpl struct {
	repo domain.GameRepository
}

func NewUserService(repo domain.GameRepository) domain.UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Register(request dto.SignUpRequest) (string, error) {
	_, err := s.repo.GetUser(request.Login)
	if err == nil {
		return "", errors.New("user already exists")
	}
	id := uuid.New()
	user := &domain.User{
		ID:       id,
		Login:    request.Login,
		Password: request.Password,
	}
	return id.String(), s.repo.SaveUser(user)
}

func (s *UserServiceImpl) AuthenticateBasic(encoded string) (string, error) {
	if !strings.HasPrefix(encoded, "Basic ") {
		return "", errors.New("invalid format")
	}

	trimed := strings.TrimPrefix(encoded, "Basic ")
	data, err := base64.StdEncoding.DecodeString(trimed)
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid format")
	}
	user, err := s.repo.GetUser(parts[0])
	if err != nil || user.Password != parts[1] {
		return "", errors.New("invalid login or password")
	}
	return user.ID.String(), nil
}
