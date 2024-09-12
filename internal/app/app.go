package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"secondTry/internal/domain/entity"
	"secondTry/internal/domain/repository"
	"secondTry/internal/domain/service"
	"time"
)

type Service interface {
	Registration(user entity.User) error
	Authorization(user entity.User, session entity.Session) (string, time.Time, int, error)
	CheckSession(session entity.Session) (int, string, error)
}

type App struct {
	serv Service
}

var db *sql.DB

func (a *App) registrHandler(w http.ResponseWriter, r *http.Request) {
	var regUser entity.User
	err := json.NewDecoder(r.Body).Decode(&regUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.serv.Registration(regUser)
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
	session.UUID, session.Expire, session.ID, err = a.serv.Authorization(user, session)
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
		var err error
		session.ID, session.UUID, err = a.serv.CheckSession(session)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		// Возвращаем сессию в формате JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode()

		// Добавляем userID в контекст
		ctx := context.WithValue(r.Context(), "userID", session.ID)
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

	// help pls
	repo := repository.NewPostgresUserRepository(db)
	serv := service.NewUserService(repo)
	a.serv = serv

	r := mux.NewRouter()

	r.HandleFunc("/reg", a.registrHandler).Methods("POST")
	r.HandleFunc("/login", a.loginHandler).Methods("POST")

	// Применяем middleware к защищенному маршруту
	r.Handle("/protect", a.sessionMiddleware(http.HandlerFunc(protectedHandler))).Methods("GET")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
