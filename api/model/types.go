package model

import "time"

type User struct {
	ID        int    `json:"id" gorm:"serial;primaryKey"`
	Firstname string `json:"firstname" gorm:"not null"`
	Lastname  string `json:"lastname" gorm:"not null"`
	Email     string `json:"email" gorm:"not null;unique"`
	Password  string `json:"password" gorm:"not null"`
	Role      string `json:"role" gorm:"not null;default: buyer"`
}

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
//	Id   int    `json:"id"`
//	Type string `json:"type"`
//}

type Accounts []Account

type Transaction struct {
	Id              int     `json:"id"`
	SenderAccount   string  `json:"sender_account"`
	ReceiverAccount string  `json:"receiver_account"`
	Amount          float64 `json:"amount"`
}

type Book struct {
	Id          int       `json:"id" gorm:"serial;primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	SellerId    int       `json:"seller_id" gorm:"not null"`
	CreatedAt   time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"-" gorm:"autoUpdateTime"`
	DeletedAt   time.Time `json:"-"`
}
