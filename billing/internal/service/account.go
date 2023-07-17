package service

import (
	"biling-nats/billing/internal/repository"
	"biling-nats/billing/model"
	"log"
)

type AccountService struct {
	repo repository.Account
}

func NewAuthService(repo repository.Account) *AccountService {
	return &AccountService{repo: repo}
}

func (a *AccountService) CreateAccount(acc *model.Account) (string, error) {
	id, err := a.repo.CreateAccount(acc)
	if err != nil {
		log.Println("could not create a new account due to: ", err.Error())
		return "", err
	}
	return id, nil
}

func (a *AccountService) GetAccounts(userId int) (model.Accounts, error) {
	accounts, err := a.repo.GetAccounts(userId)
	if err != nil {
		log.Println("failed to get list of accounts. error is:", err.Error())
		return nil, err
	}

	return accounts, nil
}

func (a *AccountService) Transaction(tr model.Transaction) (int, error) {
	_, err := a.repo.GetAccountById(tr.SenderAccount)
	if err != nil {
		log.Println("failed to get account by id. error is: ", err.Error())
		return -1, err
	}
	_, err = a.repo.GetAccountById(tr.ReceiverAccount)
	if err != nil {
		log.Println("failed to get account by id. error is: ", err.Error())
		return -1, err
	}

	if err := a.repo.Transaction(tr); err != nil {
		log.Println("unsuccessful operation. error is:", err.Error())
		return -1, err
	}

	id, err := a.repo.CreateTransaction(tr)
	if err != nil {
		return -1, err
	}

	return id, nil
}
