package service_test

import (
	"secondTry/internal/domain/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"secondTry/internal/domain/entity"
)

// MockRepository - это структура, которая реализует интерфейс Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) UserRegistration(user entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetUserHashedPass(user entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockRepository) SessionRegistration(session entity.Session, user entity.User) (string, time.Time, error) {
	args := m.Called(session, user)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockRepository) GetUserID(user entity.User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetSessionID(session entity.Session) (int, error) {
	args := m.Called(session)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetSessionUUID(session entity.Session) (string, error) {
	args := m.Called(session)
	return args.String(0), args.Error(1)
}

func TestRegistration(t *testing.T) {
	mockRepo := new(MockRepository)
	userService := service.NewUserService(mockRepo)

	user := entity.User{Login: "testuser", Password: "password123"}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	mockRepo.On("UserRegistration", user).Return(nil)

	err := userService.Registration(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthorization_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	userService := service.NewUserService(mockRepo)

	user := entity.User{Login: "testuser", Password: "password123"}
	session := entity.Session{UUID: "uuid123", Expire: time.Now().Add(24 * time.Hour)}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	mockRepo.On("GetUserID", user).Return(1, nil)
	mockRepo.On("GetUserHashedPass", user).Return(string(hashedPassword), nil)
	mockRepo.On("SessionRegistration", session, user).Return("uuid123", session.Expire, nil)

	sessionUUID, expireTime, sessionID, err := userService.Authorization(user, session)

	assert.NoError(t, err)
	assert.Equal(t, "uuid123", sessionUUID)
	assert.Equal(t, session.Expire, expireTime)
	assert.Equal(t, 0, sessionID) // Предполагается, что ID сессии не меняется в этом тесте
	mockRepo.AssertExpectations(t)
}

func TestAuthorization_FailedPasswordCheck(t *testing.T) {
	mockRepo := new(MockRepository)
	userService := service.NewUserService(mockRepo)

	user := entity.User{Login: "testuser", Password: "wrongpassword"}
	session := entity.Session{UUID: "uuid123", Expire: time.Now().Add(24 * time.Hour)}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockRepo.On("GetUserID", user).Return(1, nil)
	mockRepo.On("GetUserHashedPass", user).Return(string(hashedPassword), nil)

	sessionUUID, expireTime, sessionID, err := userService.Authorization(user, session)

	assert.Error(t, err)
	assert.Equal(t, "uuid123", sessionUUID)
	assert.Equal(t, session.Expire, expireTime)
	assert.Equal(t, 0, sessionID) // Предполагается, что ID сессии не меняется в этом тесте
	mockRepo.AssertExpectations(t)
}

func TestCheckSession(t *testing.T) {
	mockRepo := new(MockRepository)
	userService := service.NewUserService(mockRepo)

	session := entity.Session{UUID: "uuid123"}

	mockRepo.On("GetSessionID", session).Return(1, nil)
	mockRepo.On("GetSessionUUID", session).Return("uuid123", nil)

	sessionID, sessionUUID, err := userService.CheckSession(session)

	assert.NoError(t, err)
	assert.Equal(t, 1, sessionID)
	assert.Equal(t, "uuid123", sessionUUID)
	mockRepo.AssertExpectations(t)
}
