package model

type User struct {
	ID        int    `json:"id" gorm:"serial;primaryKey"`
	Firstname string `json:"firstname" gorm:"not null"`
	Lastname  string `json:"lastname" gorm:"not null"`
	Email     string `json:"email" gorm:"not null;unique"`
	Password  string `json:"password" gorm:"not null"`
	Role      string `json:"role" gorm:"not null;default: buyer"`
}
