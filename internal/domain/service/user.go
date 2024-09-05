package service

import (
	"secondTry/internal/domain/entity"
)

type Repository interface {
	UserRegistration(user entity.User) error
	UserLogin(session *entity.Session, user entity.User) (entity.Session, error)
	SessionRegistration(session *entity.Session, user entity.User) (entity.Session, error)
}

type UserRepository struct {
	repo Repository
}

func (r *UserRepository) HashPassword(user entity.User) (string, error) {

}
