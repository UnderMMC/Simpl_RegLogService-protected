package service

import (
	"golang.org/x/crypto/bcrypt"
	"secondTry/internal/domain/entity"
	"time"
)

type Repository interface {
	UserRegistration(user entity.User) error
	GetUserHashedPass(user entity.User) (string, error)
	SessionRegistration(session entity.Session, user entity.User) (string, time.Time, error)
	GetUserID(user entity.User) (int, error)
	GetSessionID(session entity.Session) (int, error)
	GetSessionUUID(session entity.Session) (string, error)
}

type Service struct {
	repo Repository
}

func NewUserService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Registration(user entity.User) error {
	hashedPassword, err := hashPassword(user.Password)
	user.Password = hashedPassword
	err = s.repo.UserRegistration(user)
	if err != nil {
		return err
	}
	return err
}

func (s *Service) Authorization(user entity.User, session entity.Session) (string, time.Time, int, error) {
	userID, err := s.repo.GetUserID(user)
	user.ID = userID

	var hashedPassword string
	hashedPassword, err = s.repo.GetUserHashedPass(user)
	if err != nil {
		return session.UUID, session.Expire, session.ID, err
	}
	if checkPasswordHash(user.Password, hashedPassword) != nil {
		return session.UUID, session.Expire, session.ID, err
	} else {
		session.UUID, session.Expire, err = s.repo.SessionRegistration(session, user)
	}
	return session.UUID, session.Expire, session.ID, err
}

func (s *Service) CheckSession(session entity.Session) (int, string, error) {
	sessionID, err := s.repo.GetSessionID(session)
	if err != nil {
		return 0, "", err
	} else {
		session.ID = sessionID
	}
	sessionUUID, err := s.repo.GetSessionUUID(session)
	if err != nil {
		return 0, "", err
	} else {
		session.UUID = sessionUUID
	}
	return session.ID, session.UUID, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
