package main

import (
	"biling-nats/book-store/config"
	"biling-nats/book-store/db"
	"biling-nats/book-store/internal/handler"
	"biling-nats/book-store/internal/repository"
	"biling-nats/book-store/internal/service"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load("./auth/.env"); err != nil {
		log.Fatal(err)
	}
	//reading from yaml
	if err := InitConfigs(); err != nil {
		log.Fatalf("error while reading config file. error is %v", err.Error())
	}

	conn, err := nats.Connect(nats.DefaultURL, nats.Name("sharif"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the nats server")

	defer conn.Close()

	var cfg config.DatabaseConnConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Couldn't unmarshal the config into struct. error is %v", err.Error())
	}
	cfg.Password = os.Getenv("DB_PASSWORD")
	//log.Println(cfg)

	database, err := db.GetDBConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	newRepository := repository.NewRepository(database)
	newService := service.NewService(newRepository)
	newHandler := handler.NewHandler(conn, newService)
	newHandler.Init()
}

func InitConfigs() error {
	viper.AddConfigPath("book-store/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
