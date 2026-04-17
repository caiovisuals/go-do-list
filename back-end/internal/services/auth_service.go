package services

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"go-do-list/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthService) Register(input RegisterInput) (*models.User, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("name, email and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.db.QueryRow(
		`INSERT INTO users (name, email, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, email, created_at`,
		input.Name, input.Email, string(hash),
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, errors.New("email already in use")
	}

	return &user, nil
}

func (s *AuthService) Login(input LoginInput) (string, *models.User, error) {
	var user models.User
	var hash string

	err := s.db.QueryRow(
		`SELECT id, name, email, password_hash, created_at
		 FROM users WHERE email = $1`,
		input.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &hash, &user.CreatedAt)

	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return tokenStr, &user, nil
}

func (s *AuthService) GetByID(id string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(
		`SELECT id, name, email, created_at FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
