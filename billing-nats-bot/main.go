package main

import (
	"biling-nats/billing-nats-bot/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "bot-api", log.LstdFlags)

	if err := godotenv.Load("./billing-nats-bot/.env"); err != nil {
		log.Fatal(err)
	}

	conn, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		log.Fatal("could not connect to nats server. Error is:", err.Error())
	}
	ok := conn.IsConnected()
	if ok != true {
		log.Fatal("not connected to the nats server")
	}

	log.Println("Connected to the nats server")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal("failed to creates a new BotAPI. Error is: ", err.Error())
	}

	newHandler := handler.NewHandler(conn, bot, l)
	newHandler.Init()
}
