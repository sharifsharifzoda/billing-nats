package main

import (
	"biling-nats/billing/config"
	"biling-nats/billing/db"
	"biling-nats/billing/internal/handler"
	"biling-nats/billing/internal/repository"
	"biling-nats/billing/internal/service"
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
		log.Fatal("failed while connection to the nats server. Error is:", err)
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
	viper.AddConfigPath("billing/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}
