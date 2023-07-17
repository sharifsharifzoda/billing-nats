package repository

import (
	"biling-nats/auth/model"
	"gorm.io/gorm"
)

type Authorization interface {
	IsEmailUsed(email string) bool
	CreateUser(user *model.User) (int, error)
	GetUser(email string) (model.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{NewAuthRepository(db)}
}
