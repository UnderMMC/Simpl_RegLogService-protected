package repository_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"secondTry/internal/domain/entity"
	"secondTry/mocks"
	"testing"
	"time"
)

func TestUserRegistration(t *testing.T) {
	mockRepo := new(mocks.Repository)
	user := entity.User{Login: "testuser", Password: "testpassword"}

	mockRepo.On("UserRegistration", user).Return(nil)

	err := mockRepo.UserRegistration(user)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetUserHashedPass(t *testing.T) {
	mockRepo := new(mocks.Repository)
	user := entity.User{Login: "testuser"}
	expectedPassword := "testpassword"

	mockRepo.On("GetUserHashedPass", user).Return(expectedPassword, nil)

	password, err := mockRepo.GetUserHashedPass(user)
	assert.NoError(t, err)
	assert.Equal(t, expectedPassword, password)

	mockRepo.AssertExpectations(t)
}

func TestSessionRegistration(t *testing.T) {
	mockRepo := new(mocks.Repository)
	session := entity.Session{}
	user := entity.User{ID: 1} // Assuming user ID is 1
	expectedUUID := uuid.New().String()
	expectedExpire := time.Now().Add(2 * time.Minute)

	mockRepo.On("SessionRegistration", session, user).Return(expectedUUID, expectedExpire, nil)

	uuid, expire, err := mockRepo.SessionRegistration(session, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, uuid)
	assert.True(t, expire.After(time.Now()))

	mockRepo.AssertExpectations(t)
}

func TestGetUserID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	user := entity.User{Login: "testuser"}
	expectedID := 1

	mockRepo.On("GetUserID", user).Return(expectedID, nil)

	id, err := mockRepo.GetUserID(user)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	mockRepo.AssertExpectations(t)
}

func TestGetSessionID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	session := entity.Session{UUID: "some-uuid"}
	expectedID := 1

	mockRepo.On("GetSessionID", session).Return(expectedID, nil)

	id, err := mockRepo.GetSessionID(session)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)

	mockRepo.AssertExpectations(t)
}

func TestGetSessionUUID(t *testing.T) {
	mockRepo := new(mocks.Repository)
	session := entity.Session{ID: 1} // Assuming session ID is 1
	expectedUUID := "some-uuid"

	mockRepo.On("GetSessionUUID", session).Return(expectedUUID, nil)

	uuid, err := mockRepo.GetSessionUUID(session)
	assert.NoError(t, err)
	assert.Equal(t, expectedUUID, uuid)

	mockRepo.AssertExpectations(t)
}
