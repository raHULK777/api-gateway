package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-gateway/internal/repository"
	"api-gateway/internal/service"
)

// Mock implementations for testing
type mockKeyStore struct {
	mock.Mock
}

func (m *mockKeyStore) CreateUser(user *repository.User) error {
	args := m.Called(user)
	return args.Get(0).(error)
}

func (m *mockKeyStore) GetUserByEmail(email string) (*repository.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*repository.User), args.Get(1).(error)
}

func (m *mockKeyStore) GetUserByID(id uint) (*repository.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*repository.User), args.Get(1).(error)
}

func (m *mockKeyStore) CreateAPIKey(key *repository.APIKey) error {
	args := m.Called(key)
	return args.Get(0).(error)
}

func (m *mockKeyStore) ValidateAPIKey(key string) (*repository.APIKey, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*repository.APIKey), args.Get(1).(error)
}

func (m *mockKeyStore) GetUserItems(userID uint) ([]repository.Item, error) {
	args := m.Called(userID)
	return args.Get(0).([]repository.Item), args.Get(1).(error)
}

func (m *mockKeyStore) CreateItem(item *repository.Item) error {
	args := m.Called(item)
	return args.Get(0).(error)
}

func (m *mockKeyStore) GetItem(id uint) (*repository.Item, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*repository.Item), args.Get(1).(error)
}

func (m *mockKeyStore) UpdateItem(item *repository.Item) error {
	args := m.Called(item)
	return args.Get(0).(error)
}

func (m *mockKeyStore) DeleteItem(id uint) error {
	args := m.Called(id)
	return args.Get(0).(error)
}

type mockTokenService struct {
	mock.Mock
}

func (m *mockTokenService) GenerateToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Get(1).(error)
}

func (m *mockTokenService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*jwt.MapClaims), args.Get(1).(error)
}

func (m *mockTokenService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Get(1).(error)
}

func (m *mockTokenService) CheckPassword(hash, password string) error {
	args := m.Called(hash, password)
	return args.Get(0).(error)
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKS := new(mockKeyStore)
	mockTS := new(mockTokenService)

	keyStore := &repository.KeyStore(mockKS)
	tokenService := &service.TokenService(mockTS)

	router := gin.New()
	router.POST("/register", Register(keyStore, tokenService))

	// Mock responses
	mockKS.On("GetUserByEmail", "test@example.com").Return((*repository.User)(nil), repository.ErrRecordNotFound)
	mockTS.On("HashPassword", "password123").Return("hashedpassword", nil)
	mockKS.On("CreateUser", mock.MatchedBy(func(u *repository.User) bool {
		return u.Email == "test@example.com"
	})).Return(nil)
	mockTS.On("GenerateToken", uint(1), "test@example.com").Return("test-token", nil)

	w := httptest.NewRecorder()
	body := `{"email":"test@example.com","password":"password123","name":"Test User"}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKS := new(mockKeyStore)
	mockTS := new(mockTokenService)

	keyStore := &repository.KeyStore(mockKS)
	tokenService := &service.TokenService(mockTS)

	router := gin.New()
	router.POST("/login", Login(keyStore, tokenService))

	// Mock responses
	user := &repository.User{ID: 1, Email: "test@example.com", Password: "hashedpassword"}
	mockKS.On("GetUserByEmail", "test@example.com").Return(user, nil)
	mockTS.On("CheckPassword", "hashedpassword", "password123").Return(nil)
	mockTS.On("GenerateToken", uint(1), "test@example.com").Return("test-token", nil)

	w := httptest.NewRecorder()
	body := `{"email":"test@example.com","password":"password123"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}
