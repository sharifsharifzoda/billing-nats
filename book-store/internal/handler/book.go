package handler

import (
	"biling-nats/book-store/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
	"time"
)

type BuyBook struct {
	Id      int
	BuyerId int
}

type Transfer struct {
	BuyerId  int
	SellerId int
	Amount   float64
}

func (h *Handler) AddBook(msg *nats.Msg) {
	var book model.Book
	if err := json.Unmarshal(msg.Data, &book); err != nil {
		msg.Respond([]byte("could not unmarshal the data"))
		log.Println("could not unmarshal the data. Error is:", err.Error())
		return
	}

	bookId, err := h.Service.Book.AddBook(book)
	if err != nil {
		msg.Respond([]byte("failed to create a new book"))
		return
	}

	idStr := strconv.Itoa(bookId)

	if err := msg.Respond([]byte(idStr)); err != nil {
		log.Println("failed to send the response. Error is:", err.Error())
	}
}

func (h *Handler) BuyBook(msg *nats.Msg) {
	fmt.Println(string(msg.Data))
	var book BuyBook
	if err := json.Unmarshal(msg.Data, &book); err != nil {
		msg.Respond([]byte("failed to unmarshal the data"))
		log.Println("failed to unmarshal the data. Error is:", err.Error())
		return
	}

	selected, err := h.Service.Book.GetBook(book.Id)
	if err != nil {
		msg.Respond([]byte("failed to get a book by id"))
		return
	}

	var transfer = Transfer{
		BuyerId:  book.BuyerId,
		SellerId: selected.SellerId,
		Amount:   selected.Price,
	}

	bytes, err := json.Marshal(transfer)
	if err != nil {
		msg.Respond([]byte("failed to marshal the data"))
		log.Println("failed to marshal the data. Error is:", err.Error())
		return
	}

	request, err := h.Nats.Request("book.transfer", bytes, 10*time.Second)
	if err != nil {
		msg.Respond([]byte("failed to send the request for transaction"))
		log.Println("failed to send the request. Error is:", err.Error())
		return
	}

	msg.Respond(request.Data)

	if err := h.Nats.Publish("bot.transfer", request.Data); err != nil {
		log.Println("failed to publish the message to telegram bot. Error is:", err.Error())
		return
	}
}
