package repository

import (
	"biling-nats/billing/model"
	"gorm.io/gorm"
)

type Account interface {
	CreateAccount(acc *model.Account) (string, error)
	GetAccounts(userId int) (model.Accounts, error)
	GetAccountById(id string) (model.Account, error)
	Transaction(tr model.Transaction) error
	CreateTransaction(tr model.Transaction) (int, error)
}

type Repository struct {
	Account
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{NewAccountRepository(db)}
}
