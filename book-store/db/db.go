package db

import (
	"biling-nats/book-store/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	CreateBookTable = `CREATE TABLE IF NOT EXISTS books (
    	id SERIAL PRIMARY KEY NOT NULL,
    	name VARCHAR NOT NULL,
    	description VARCHAR NOT NULL,
    	price DECIMAL NOT NULL,
    	seller_id INTEGER NOT NULL,
    	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP,
    	deleted_at TIMESTAMP
	);`
)

func GetDBConnection(cfg config.DatabaseConnConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Dushanbe",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	log.Printf("Connection success host:%s port:%s", cfg.Host, cfg.Port)

	Init(conn)

	return conn, nil
}

func Init(db *sql.DB) {
	DDLs := []string{
		CreateBookTable,
	}

	for i, ddl := range DDLs {
		_, err := db.Exec(ddl)
		if err != nil {
			log.Fatalf("failed to create table #%d due to: %s", i+1, err.Error())
		}
	}
}
