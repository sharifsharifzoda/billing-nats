package repository

import (
	"biling-nats/billing/model"
	"errors"
	"gorm.io/gorm"
)

type AccountRepository struct {
	Db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{Db: db}
}

const (
	CallProcedure = `call transfer($1, $2, $3, null);`
)

func (a *AccountRepository) CreateAccount(acc *model.Account) (string, error) {
	err := a.Db.Model(&model.Account{}).Create(&acc).Error
	if err != nil {
		return "", err
	}

	return acc.Id, nil
}

func (a *AccountRepository) GetAccounts(userId int) (model.Accounts, error) {
	var accounts model.Accounts
	err := a.Db.Model(&model.Account{}).Where("user_id = ?", userId).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (a *AccountRepository) GetAccountById(id string) (model.Account, error) {
	var acc model.Account
	err := a.Db.Model(&model.Account{}).Where("id = ?", id).First(&acc).Error
	if err != nil {
		return model.Account{}, err
	}

	return acc, nil
}

func (a *AccountRepository) Transaction(tr model.Transaction) error {
	var result string
	db, _ := a.Db.DB()
	row := db.QueryRow(CallProcedure, tr.SenderAccount, tr.ReceiverAccount, tr.Amount)
	err := row.Scan(&result)
	if err != nil {
		return err
	}
	if result == "ok" {
		return nil
	}

	err = errors.New(result)
	//fmt.Println(err)
	return err
}

func (a *AccountRepository) CreateTransaction(tr model.Transaction) (int, error) {
	err := a.Db.Model(&model.Transaction{}).Create(&tr).Error
	if err != nil {
		return -1, err
	}
	return tr.Id, nil
}
