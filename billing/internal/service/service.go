package service

import (
	"biling-nats/billing/internal/repository"
	"biling-nats/billing/model"
)

type Account interface {
	CreateAccount(acc *model.Account) (string, error)
	GetAccounts(userId int) (model.Accounts, error)
	Transaction(tr model.Transaction) (int, error)
}

type Service struct {
	Account
}

func NewService(repo *repository.Repository) *Service {
	return &Service{NewAuthService(repo.Account)}
}
