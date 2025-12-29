package services

import (
	"testing"
	"todo-go-backend/internal/errors"

	"github.com/stretchr/testify/assert"
)

// MockUserRepository é um mock do UserRepository para testes
type MockUserRepository struct {
	users        map[uint]*MockUser
	usersByUser  map[string]*MockUser
	usersByEmail map[string]*MockUser
	nextID       uint
}

type MockUser struct {
	ID       uint
	Username string
	Email    string
	Password string
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:        make(map[uint]*MockUser),
		usersByUser:  make(map[string]*MockUser),
		usersByEmail: make(map[string]*MockUser),
		nextID:       1,
	}
}

func (m *MockUserRepository) Create(user interface{}) error {
	u := user.(*MockUser)
	u.ID = m.nextID
	m.nextID++
	m.users[u.ID] = u
	m.usersByUser[u.Username] = u
	m.usersByEmail[u.Email] = u
	return nil
}

func (m *MockUserRepository) FindByID(id uint) (interface{}, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByUsername(username string) (interface{}, error) {
	user, ok := m.usersByUser[username]
	if !ok {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(email string) (interface{}, error) {
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByUsernameOrEmail(username, email string) (interface{}, error) {
	if user, ok := m.usersByUser[username]; ok {
		return user, nil
	}
	if user, ok := m.usersByEmail[email]; ok {
		return user, nil
	}
	return nil, errors.ErrUserNotFound
}

func (m *MockUserRepository) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	_, userExists := m.usersByUser[username]
	_, emailExists := m.usersByEmail[email]
	return userExists || emailExists, nil
}

// Implementação real do UserRepository para o mock
type mockUserRepo struct {
	mock *MockUserRepository
}

func (m *mockUserRepo) Create(user interface{}) error {
	return m.mock.Create(user)
}

func (m *mockUserRepo) FindByID(id uint) (interface{}, error) {
	return m.mock.FindByID(id)
}

func (m *mockUserRepo) FindByUsername(username string) (interface{}, error) {
	return m.mock.FindByUsername(username)
}

func (m *mockUserRepo) FindByEmail(email string) (interface{}, error) {
	return m.mock.FindByEmail(email)
}

func (m *mockUserRepo) FindByUsernameOrEmail(username, email string) (interface{}, error) {
	return m.mock.FindByUsernameOrEmail(username, email)
}

func (m *mockUserRepo) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	return m.mock.ExistsByUsernameOrEmail(username, email)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret")

	t.Run("Register new user successfully", func(t *testing.T) {
		user, token, err := service.Register("testuser", "test@example.com", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("Register duplicate username", func(t *testing.T) {
		_, _, err := service.Register("testuser", "test2@example.com", "password123")

		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
		appErr := err.(*errors.AppError)
		assert.Equal(t, errors.ErrUserAlreadyExists, appErr.Err)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret")

	// Create a user first
	_, _, _ = service.Register("testuser", "test@example.com", "password123")

	t.Run("Login with valid credentials", func(t *testing.T) {
		user, token, err := service.Login("testuser", "password123")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
	})

	t.Run("Login with invalid password", func(t *testing.T) {
		_, _, err := service.Login("testuser", "wrongpassword")

		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
		appErr := err.(*errors.AppError)
		assert.Equal(t, errors.ErrInvalidCredentials, appErr.Err)
	})

	t.Run("Login with non-existent user", func(t *testing.T) {
		_, _, err := service.Login("nonexistent", "password123")

		assert.Error(t, err)
		assert.IsType(t, &errors.AppError{}, err)
		appErr := err.(*errors.AppError)
		assert.Equal(t, errors.ErrInvalidCredentials, appErr.Err)
	})
}
