package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nats-io/nats.go"
	"log"
)

const (
	commandStart = "start"
)

var ChatId []int64

func (h *Handler) SendMessageToBot(msg *nats.Msg) {
	fmt.Println(string(msg.Data))

	for _, id := range ChatId {
		message := tgbotapi.NewMessage(id, string(msg.Data))
		_, err := h.Bot.Send(message)
		if err != nil {
			log.Println("failed to send the message to Telegram. Error is: ", err.Error())
			continue
		}
	}

}

func (h *Handler) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return h.handleStartCommand(message)
	default:
		return h.handleUnknownCommand(message)
	}
}

func (h *Handler) handleStartCommand(message *tgbotapi.Message) error {
	ChatId = append(ChatId, message.Chat.ID)
	fmt.Println(ChatId)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Already authorized")
	_, err := h.Bot.Send(msg)
	return err
}

func (h *Handler) handleMessage(message *tgbotapi.Message) error {
	//for sub.IsValid() {
	//	m, err := sub.NextMsg(10 * time.Second)
	//	if err != nil {
	//		log.Println("Error is: ", err.Error())
	//		return err
	//	}
	//
	//	msg := tgbotapi.NewMessage(message.Chat.ID, string(m.Data))
	//	_, err = h.Bot.Send(msg)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//return nil

	//for msg := range ch {
	//	newMessage := tgbotapi.NewMessage(message.Chat.ID, string(msg.Data))
	//	_, err := h.Bot.Send(newMessage)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

func (h *Handler) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command")
	_, err := h.Bot.Send(msg)
	return err
}
