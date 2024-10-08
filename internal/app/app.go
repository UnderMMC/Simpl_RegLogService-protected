package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"secondTry/internal/domain/entity"
	"secondTry/internal/domain/repository"

	_ "github.com/lib/pq"
)

type Repository interface {
	UserRegistration(user entity.User) error
	UserLogin(session *entity.Session, user entity.User) (entity.Session, error)
	SessionRegistration(session *entity.Session, user entity.User) (entity.Session, error)
}

type App struct {
	repo Repository
}

var db *sql.DB

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (a *App) registrHandler(w http.ResponseWriter, r *http.Request) {
	var regUser entity.User
	err := json.NewDecoder(r.Body).Decode(&regUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	regUser.Password, err = hashPassword(regUser.Password)
	err = a.repo.UserRegistration(regUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	return
}

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var session entity.Session
	session, err = a.repo.UserLogin(&session, user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	// Возвращаем сессию в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// Middleware для проверки сессии
func (a *App) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UUID := r.Header.Get("UUID")
		if UUID == "" {
			http.Error(w, "UUID is missing", http.StatusUnauthorized)
			return
		}

		// Проверка существования сессии в базе данных
		var session entity.Session
		var userID int
		err := db.QueryRow("SELECT uuid FROM sessions WHERE user_id=$1", session.ID).Scan(&userID)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Добавляем userID в контекст
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	// Создаем ответ в формате JSON
	response := map[string]string{
		"message": fmt.Sprintf("Welcome, %s!", userID),
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Устанавливаем статус 200 OK

	// Кодируем ответ в JSON и отправляем его клиенту
	json.NewEncoder(w).Encode(response)
}

//nolint:exhaustruct
func New() *App {
	return &App{}
}

func (a *App) Run() {
	var err error
	connStr := "user=postgres password=pgpwd4habr dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewPostgresUserRepository(db)
	a.repo = userRepo

	r := mux.NewRouter()

	r.HandleFunc("/reg", a.registrHandler).Methods("POST")
	r.HandleFunc("/login", a.loginHandler).Methods("POST")

	// Применяем middleware к защищенному маршруту
	r.Handle("/protected", a.sessionMiddleware(http.HandlerFunc(protectedHandler))).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
