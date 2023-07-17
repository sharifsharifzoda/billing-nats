package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nats-io/nats.go"
)

const (
	commandStart = "start"
)

var ChatId []int64

func (h *Handler) SendMessageToBot(msg *nats.Msg) {
	fmt.Println(string(msg.Data))

	//update := tgbotapi.NewUpdate(0)
	//update.Timeout = 60

	//fmt.Println(string(msg.Data))
	//updates, err := h.Bot.GetUpdatesChan(update)
	//if err != nil {
	//	h.Logger.Println("failed to get updates. Error is:", err.Error())
	//	return
	//}

	//for update := range updates {
	//	if strings.ToLower(update.Message.Text) == "/news" {
	//		trId := fmt.Sprintf("Transaction #%v completed successfully", string(msg.Data))
	//		fmt.Println("Received:", update.Message.Text, trId, "->", update.Message.From)
	//
	//		message := tgbotapi.NewMessage(update.Message.Chat.ID, trId)
	//		_, err := h.Bot.Send(message)
	//		if err != nil {
	//			log.Println("failed to send the message to bot. Error is:", err.Error())
	//			continue
	//		}
	//
	//	}
	//
	//	if strings.ToLower(update.Message.Text) == "/start" {
	//		ChatId = append(ChatId, update.Message.Chat.ID)
	//		fmt.Println(update.Message.Chat.ID)
	//	}
	//
	//	for _, id := range ChatId {
	//		log.Println("sending...")
	//		message := tgbotapi.NewMessage(id, string(msg.Data))
	//		_, err := h.Bot.Send(message)
	//		if err != nil {
	//			log.Println("failed to send the message to bot. Error is:", err.Error())
	//			continue
	//		}
	//	}
	//
	//}

	//message := tgbotapi.NewMessage(1238370443, string(msg.Data))
	//_, err := h.Bot.Send(message)
	//if err != nil {
	//	log.Println("failed to send the message to bot. Error is:", err.Error())
	//}
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
	msg := tgbotapi.NewMessage(message.Chat.ID, "Already authorized")
	_, err := h.Bot.Send(msg)
	return err
}

func (h *Handler) handleMessage(message *tgbotapi.Message, ch chan *nats.Msg) error {
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

	for msg := range ch {
		newMessage := tgbotapi.NewMessage(message.Chat.ID, string(msg.Data))
		_, err := h.Bot.Send(newMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command")
	_, err := h.Bot.Send(msg)
	return err
}
