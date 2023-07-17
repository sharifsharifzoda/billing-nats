package handler

import (
	"biling-nats/billing/internal/service"
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
	sub, err := h.Nats.Subscribe("account.create", h.CreateAccount)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())
	}
	defer sub.Unsubscribe()

	//-------------------------------------------------------------------------------

	getSub, err := h.Nats.Subscribe("account.get", h.GetAccounts)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %v", err.Error())
	}
	defer getSub.Unsubscribe()

	//------------------------------------------------------------------------------

	transferSub, err := h.Nats.Subscribe("account.transfer", h.Transaction)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %v", err.Error())
	}
	defer transferSub.Unsubscribe()

	//-----------------------------------------------------------------------------------

	BookSub, err := h.Nats.Subscribe("book.transfer", h.BuyBook)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %v", err.Error())
	}
	defer BookSub.Unsubscribe()

	//-----------------------------------------------------------------------------------
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	msg := <-ch
	fmt.Println("Got signal: ", msg)
	log.Println("Shutting down!")
}
