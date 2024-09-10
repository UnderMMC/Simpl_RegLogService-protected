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
}

type Service struct {
	repo Repository
}

func (s *Service) Registration(user entity.User) error {
	hashedPassword, err := hashPassword(user.Password)
	err = s.repo.UserRegistration(user)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return err
}

func (s *Service) Authorization(user entity.User, session entity.Session) error {
	userID, err := s.repo.GetUserID(user)
	user.ID = userID

	var hashedPassword string
	hashedPassword, err = s.repo.GetUserHashedPass(user)
	if err != nil {
		return err
	}
	if checkPasswordHash(user.Password, hashedPassword) != nil {
		return err
	} else {
		session.UUID, session.Expire, err = s.repo.SessionRegistration(session, user)
	}
	return err
}

func (s *Service) CheckSession(session entity.Session) (int, error) {
	sessionID, err := s.repo.GetSessionID(session)
	if err != nil {
		return 0, err
	} else {
		session.ID = sessionID
	}
	return session.ID, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
