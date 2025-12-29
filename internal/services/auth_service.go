package services

import (
	"todo-go-backend/internal/errors"
	"todo-go-backend/internal/models"
	"todo-go-backend/internal/repositories"
	"todo-go-backend/pkg/utils"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(username, email, password string) (*models.User, string, error)
	Login(identifier, password string) (*models.User, string, error) // identifier can be username or email
}

type authService struct {
	userRepo repositories.UserRepository
	jwtSecret string
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(username, email, password string) (*models.User, string, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByUsernameOrEmail(username, email)
	if err != nil {
		return nil, "", errors.NewInternalServerError(err)
	}
	if exists {
		return nil, "", errors.NewUserAlreadyExistsError()
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", errors.NewInternalServerError(err)
	}

	// Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", errors.NewInternalServerError(err)
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret)
	if err != nil {
		return nil, "", errors.NewInternalServerError(err)
	}

	return user, token, nil
}

func (s *authService) Login(identifier, password string) (*models.User, string, error) {
	// Find user by username or email
	user, err := s.userRepo.FindByUsernameOrEmailValue(identifier)
	if err != nil {
		return nil, "", errors.NewInvalidCredentialsError()
	}

	// Verify password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.NewInvalidCredentialsError()
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, s.jwtSecret)
	if err != nil {
		return nil, "", errors.NewInternalServerError(err)
	}

	return user, token, nil
}

