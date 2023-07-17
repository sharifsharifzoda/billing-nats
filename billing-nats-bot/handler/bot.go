package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nats-io/nats.go"
	"log"
)

func (h *Handler) SendMessageToBot(msg *nats.Msg) {
	fmt.Println(string(msg.Data))

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	fmt.Println(string(msg.Data))
	updates, err := h.Bot.GetUpdatesChan(update)
	if err != nil {
		h.Logger.Println("failed to get updates. Error is:", err.Error())
		return
	}

	for update := range updates {
		//if update.Message.Text == "start"
		if update.Message != nil {
			trId := fmt.Sprintf("Transaction #%v completed successfully", string(msg.Data))
			fmt.Println("Received:", update.Message.Text, trId, "->", update.Message.From)

			message := tgbotapi.NewMessage(update.Message.Chat.ID, trId)
			_, err := h.Bot.Send(message)
			if err != nil {
				log.Println("failed to send the message to bot. Error is:", err.Error())
				continue
			}

		}

	}

}
