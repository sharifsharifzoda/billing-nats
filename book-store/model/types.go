package model

import "time"

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
