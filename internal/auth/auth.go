package auth

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/internal/models"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

// User представляет модель пользователя
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // хранится как хеш
}

// AuthHandler обрабатывает запросы аутентификации
type AuthHandler struct {
	users     map[string]User
	mu        sync.RWMutex
	TokenAuth *jwtauth.JWTAuth
}

// NewAuthHandler создает новый обработчик аутентификации
func NewAuthHandler(tokenSecret string) *AuthHandler {
	return &AuthHandler{
		users:     make(map[string]User),
		TokenAuth: jwtauth.New("HS256", []byte(tokenSecret), nil),
	}
}

// RegisterRequest представляет запрос на регистрацию
// swagger:model RegisterRequest
type RegisterRequest struct {
	// Имя пользователя
	// required: true
	// example: admin
	Username string `json:"username"`

	// Пароль
	// required: true
	// example: password123
	Password string `json:"password"`
}

// LoginRequest представляет запрос на вход
// swagger:model LoginRequest
type LoginRequest struct {
	// Имя пользователя
	// required: true
	// example: admin
	Username string `json:"username"`

	// Пароль
	// required: true
	// example: password123
	Password string `json:"password"`
}

// LoginResponse представляет ответ на успешный вход
// swagger:model LoginResponse
type LoginResponse struct {
	// JWT токен
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
}

// ErrorResponse представляет ответ с ошибкой
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Описание ошибки
	// example: user not found
	Error string `json:"error"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body RegisterRequest true "Данные пользователя"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		models.RenderError(w, r, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		models.RenderError(w, r, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		models.RenderError(w, r, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Проверяем, существует ли пользователь
	if _, exists := h.users[req.Username]; exists {
		models.RenderError(w, r, ErrUserExists.Error(), http.StatusBadRequest)
		return
	}

	// Сохраняем пользователя
	h.users[req.Username] = User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	w.WriteHeader(http.StatusOK)
}

// Login godoc
// @Summary Вход пользователя
// @Description Аутентифицирует пользователя и возвращает JWT токен
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body LoginRequest true "Учетные данные"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		models.RenderError(w, r, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		models.RenderError(w, r, "Username and password are required", http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	user, exists := h.users[req.Username]
	h.mu.RUnlock()

	if !exists {
		models.RenderError(w, r, ErrUserNotFound.Error(), http.StatusForbidden)
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		models.RenderError(w, r, ErrInvalidPassword.Error(), http.StatusForbidden)
		return
	}

	// Генерируем JWT токен
	_, tokenString, err := h.TokenAuth.Encode(map[string]interface{}{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	})

	if err != nil {
		models.RenderError(w, r, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, LoginResponse{Token: tokenString})
}

// Authenticator middleware проверяет JWT токен
func (h *AuthHandler) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil || token == nil {
			models.RenderError(w, r, "Unauthorized", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
