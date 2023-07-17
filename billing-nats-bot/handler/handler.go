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
	msgChan := make(chan *nats.Msg, 5)

	sub, err := h.Nats.ChanSubscribe("bot.*", msgChan)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())
	}
	defer sub.Unsubscribe()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := h.Bot.GetUpdatesChan(u)
	if err != nil {
		h.Logger.Fatalln("failed to get updates. Error is:", err.Error())
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := h.handleCommand(update.Message); err != nil {
				h.handleError(update.Message.Chat.ID, err)
			}

			continue
		}

		if err := h.handleMessage(update.Message, msgChan); err != nil {
			h.handleError(update.Message.Chat.ID, err)
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	msg := <-ch
	fmt.Println("Got signal: ", msg)
	log.Println("Shutting down!")
}

func (h *Handler) handleError(id int64, err error) {
	msg := tgbotapi.NewMessage(id, err.Error())
	h.Bot.Send(msg)
}
