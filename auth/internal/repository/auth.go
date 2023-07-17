package repository

import (
	"biling-nats/auth/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	Db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{Db: db}
}

func (a *AuthRepository) IsEmailUsed(email string) bool {
	var user model.User
	tx := a.Db.Where("email = ?", email).Find(&user)
	if tx.Error != nil {
		return false
	}

	if user.Email == "" {
		return false
	}

	return true
}

func (a *AuthRepository) CreateUser(user *model.User) (int, error) {
	tx := a.Db.Create(&user)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return user.ID, nil
}

func (a *AuthRepository) GetUser(email string) (user model.User, err error) {
	tx := a.Db.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return model.User{}, tx.Error
	}

	return user, nil
}
