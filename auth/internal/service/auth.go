package service

import (
	"biling-nats/auth/internal/repository"
	"biling-nats/auth/model"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
	"time"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID   int    `json:"user_id"`
	UserRole string `json:"user_role"`
}

func (s *AuthService) ValidateUser(user model.User) error {
	if len(user.Email) > 30 || len(user.Email) < 5 {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if len(user.Password) > 20 || len(user.Password) < 6 {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "_") || strings.Contains(user.Password, "-") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "@") || strings.Contains(user.Password, "#") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "$") || strings.Contains(user.Password, "%") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "&") || strings.Contains(user.Password, "*") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "(") || strings.Contains(user.Password, ")") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, ":") || strings.Contains(user.Password, ".") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "/") || strings.Contains(user.Password, `\`) {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, ",") || strings.Contains(user.Password, ";") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "?") || strings.Contains(user.Password, `"`) {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}
	if strings.Contains(user.Password, "!") || strings.Contains(user.Password, "~") {
		log.Println("forbidden")
		return fmt.Errorf("forbidden")
	}

	return nil
}

func (s *AuthService) IsEmailUsed(email string) bool {
	isUsed := s.repo.IsEmailUsed(email)
	if isUsed {
		log.Println("email is already used")
		return true
	}
	return false
}

func (s *AuthService) CreateUser(user *model.User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to generate hash from password due to:", err.Error())
		return 0, err
	}

	user.Password = string(hash)

	id, err := s.repo.CreateUser(user)
	if err != nil {
		log.Println("failed to create user due to:", err.Error())
		return 0, err
	}
	return id, nil
}

func (s *AuthService) CheckUser(user model.User) (model.User, error) {
	u, err := s.repo.GetUser(user.Email)
	if err != nil {
		log.Println("no rows in result set")
		return model.User{}, errors.New("no rows in result set")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		log.Println("invalid email or password")
		return model.User{}, errors.New("invalid email or password")
	}

	return u, nil
}

func (s *AuthService) GenerateToken(user model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
		user.Role,
	})

	signedString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	return signedString, err
}

func (s *AuthService) ParseToken(tokenString string) (int, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return 0, "", errors.New("token expired")
	}

	return claims.UserID, claims.UserRole, nil
}
