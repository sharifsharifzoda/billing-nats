package model

import "time"

type Account struct {
	Id        string    `json:"id" gorm:"primaryKey;unique"`
	Name      string    `json:"name" gorm:"not null"`
	Type      int       `json:"type" gorm:"not null"`
	Balance   float64   `json:"balance" gorm:"not null;default: 0"`
	UserId    int       `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"-" gorm:"index"`
	//AccountType AccountType `json:"-" gorm:"foreignKey:Type"`
}

//type AccountType struct {
//	Id   int    `json:"id" gorm:"serial;primaryKey"`
//	Type string `json:"type" gorm:"not null"`
//}

type Accounts []Account

type Transaction struct {
	Id              int     `json:"id"`
	SenderAccount   string  `json:"sender_account"`
	ReceiverAccount string  `json:"receiver_account"`
	Amount          float64 `json:"amount"`
}
