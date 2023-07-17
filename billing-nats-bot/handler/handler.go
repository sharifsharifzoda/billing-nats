package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Handler struct {
	Nats   *nats.Conn
	Bot    *tgbotapi.BotAPI
	Logger *log.Logger
}

func NewHandler(nats *nats.Conn, bot *tgbotapi.BotAPI, log *log.Logger) *Handler {
	return &Handler{Nats: nats, Bot: bot, Logger: log}
}

func (h *Handler) Init() {
	sub, err := h.Nats.Subscribe("bot.*", h.SendMessageToBot)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())
	}
	sub, err = h.Nats.Subscribe("bot.*", h.Test)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())
	}
	defer sub.Unsubscribe()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	msg := <-ch
	fmt.Println("Got signal: ", msg)
	log.Println("Shutting down!")
}

func (h *Handler) Test(msg *nats.Msg) {
	fmt.Println("printing from test ----", string(msg.Data))
}
