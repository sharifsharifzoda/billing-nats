package db

import (
	"biling-nats/billing/config"
	"biling-nats/billing/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func GetDBConnection(cfg config.DatabaseConnConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Dushanbe",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	conn, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	log.Printf("Connection success host:%s port:%s", cfg.Host, cfg.Port)

	conn.AutoMigrate(&model.Account{}, model.Transaction{})

	return conn, nil
}

func Close(db *gorm.DB) {
	conn, err := db.DB()
	err = conn.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
