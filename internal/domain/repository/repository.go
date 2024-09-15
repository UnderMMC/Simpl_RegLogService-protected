package repository

import (
	"database/sql"
	"github.com/google/uuid"
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

func (r *PostgresUserRepository) GetUserHashedPass(user entity.User) (string, error) {
	var storedPassword string
	err := r.db.QueryRow("SELECT password FROM logdata WHERE login=$1", user.Login).Scan(&storedPassword)

	return storedPassword, err
}

func (r *PostgresUserRepository) SessionRegistration(session entity.Session, user entity.User) (string, time.Time, error) {
	UUID := uuid.New().String()
	Expire := time.Now().Add(2 * time.Minute)

	_, err := r.db.Exec("INSERT INTO sessions (user_id, UUID) VALUES ($1, $2)", user.ID, UUID)
	if err != nil {
		return UUID, Expire, err
	}
	return UUID, Expire, err
}

func (r *PostgresUserRepository) GetUserID(user entity.User) (int, error) {
	var ID int
	err := r.db.QueryRow("SELECT user_id FROM logdata WHERE login=$1", user.Login).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (r *PostgresUserRepository) GetSessionID(session entity.Session) (int, error) {
	var ID int
	err := r.db.QueryRow("SELECT user_id FROM sessions WHERE uuid=$1", session.UUID).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (r *PostgresUserRepository) GetSessionUUID(session entity.Session) (string, error) {
	var UUID string
	err := r.db.QueryRow("SELECT UUID FROM sessions WHERE user_id=$1", session.ID).Scan(&UUID)
	if err != nil {
		return "", err
	}
	return UUID, nil
}
