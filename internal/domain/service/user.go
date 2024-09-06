package service

import (
	"golang.org/x/crypto/bcrypt"
	"secondTry/internal/domain/entity"
)

type Repository interface {
	UserRegistration(user entity.User) error
	UserLogin(session *entity.Session, user entity.User) (entity.Session, error)
	SessionRegistration(session *entity.Session, user entity.User) (entity.Session, error)
	GetUserID(user entity.User) (int, error)
}

type UserRepository struct {
	repo Repository
}

func (s *UserRepository) Create(user entity.User) error {
	hashedPassword, err := hashPassword(user.Password)
	err = s.repo.UserRegistration(user)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
}

func (s *UserRepository) Login(user entity.User, session *entity.Session) (entity.User, error) {
	userID, err := s.repo.GetUserID(user)
	user.ID = userID
	session, err = s.repo.UserLogin(&session, user)

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
