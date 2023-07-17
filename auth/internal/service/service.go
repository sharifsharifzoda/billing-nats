package service

import (
	"biling-nats/auth/internal/repository"
	"biling-nats/auth/model"
)

type Auth interface {
	ValidateUser(user model.User) error
	IsEmailUsed(email string) bool
	CreateUser(user *model.User) (int, error)
	CheckUser(user model.User) (model.User, error)
	GenerateToken(user model.User) (string, error)
	ParseToken(token string) (int, string, error)
}

type Service struct {
	Auth
}

func NewService(repo *repository.Repository) *Service {
	return &Service{NewAuthService(repo.Authorization)}
}
