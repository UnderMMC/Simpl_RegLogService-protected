package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"secondTry/internal/domain/entity"
	"time"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) UserRegistration(user entity.User) error {
	_, err := r.db.Exec("INSERT INTO logdata (login, password) VALUES ($1, $2)", user.Login, user.Password)
	if err != nil {
		return err
	}
	return err
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (r *PostgresUserRepository) SessionRegistration(session *entity.Session, user entity.User) (entity.Session, error) {

	SessionId := uuid.New().String()
	Expire := time.Now().Add(2 * time.Minute)

	var ID int
	err := r.db.QueryRow("SELECT user_id FROM logdata WHERE login=$1", user.Login).Scan(&ID)
	if err != nil {
		return entity.Session{}, err
	}

	session = &entity.Session{
		UUID:   SessionId,
		Expire: Expire,
		ID:     ID,
	}

	_, err = r.db.Exec("INSERT INTO sessions (user_id, UUID) VALUES ($1, $2)", session.ID, session.UUID)
	if err != nil {
		return entity.Session{}, err
	}
	return *session, err
}

func (r *PostgresUserRepository) UserLogin(session *entity.Session, user entity.User) (entity.Session, error) {
	var result entity.Session
	var storedPassword string
	err := r.db.QueryRow("SELECT password FROM logdata WHERE login=$1", user.Login).Scan(&storedPassword)
	if err != nil || checkPasswordHash(user.Password, storedPassword) != nil {
		return entity.Session{}, err
	} else {
		result, err = r.SessionRegistration(session, user)
	}
	return result, err
}
