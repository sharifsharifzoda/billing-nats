package handler

import (
	"biling-nats/book-store/internal/service"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Handler struct {
	Nats    *nats.Conn
	Service *service.Service
}

func NewHandler(nats *nats.Conn, service *service.Service) *Handler {
	return &Handler{nats, service}
}

func (h *Handler) Init() {
	AddSub, err := h.Nats.Subscribe("book.add", h.AddBook)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())
	}
	defer AddSub.Unsubscribe()

	//------------------------------------------------------------------------------

	BuySub, err := h.Nats.Subscribe("book.buy", h.BuyBook)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())

	}
	defer BuySub.Unsubscribe()

	//------------------------------------------------------------------------------

	//------------------------------------------------------------------------------

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	msg := <-ch
	fmt.Println("Got signal: ", msg)
	log.Println("Shutting down!")
}
