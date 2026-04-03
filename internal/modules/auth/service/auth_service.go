package service

import (
	"errors"

	"github.com/thecemakin/hr-project/internal/modules/auth/model"
	"github.com/thecemakin/hr-project/internal/modules/auth/repository"
	"github.com/thecemakin/hr-project/internal/platform/auth"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService interface {
	Login(email, password string) (string, *model.User, error)
	Register(email, password, role string, employeeID *uint) (*model.User, error)
}

type authService struct {
	repo          repository.UserRepository
	tokenProvider *auth.TokenProvider
}

func NewAuthService(repo repository.UserRepository, tp *auth.TokenProvider) AuthService {
	return &authService{
		repo:          repo,
		tokenProvider: tp,
	}
}

func (s *authService) Login(email, password string) (string, *model.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.tokenProvider.GenerateToken(user.ID, user.Email, user.Role, user.EmployeeID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) Register(email, password, role string, employeeID *uint) (*model.User, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
		EmployeeID:   employeeID,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
