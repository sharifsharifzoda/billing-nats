package handler

import (
	"biling-nats/billing/model"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"strconv"
)

type Transfer struct {
	BuyerId  int
	SellerId int
	Amount   float64
}

func (h *Handler) CreateAccount(msg *nats.Msg) {
	var account *model.Account
	if err := json.Unmarshal(msg.Data, &account); err != nil {
		msg.Respond([]byte("error while unmarshalling"))
		log.Println("error while unmarshalling. error is: ", err.Error())
		return
	}

	idStr, err := h.Service.Account.CreateAccount(account)
	if err != nil {
		msg.Respond([]byte("failed to create a new account"))
		log.Println("failed to create a new account. Error is: ", err.Error())
		return
	}

	if err := msg.Respond([]byte(idStr)); err != nil {
		log.Println("could not send the request. Error is:", err.Error())
		return
	}
}

func (h *Handler) GetAccounts(msg *nats.Msg) {
	s := string(msg.Data)
	userId, err := strconv.Atoi(s)
	if err != nil {
		msg.Respond([]byte("invalid type provided for converting to integer"))
		log.Println("invalid type provided for converting to integer. error is:", err.Error())
		return
	}

	accounts, err := h.Service.Account.GetAccounts(userId)
	if err != nil {
		msg.Respond([]byte("failed to get the list of accounts"))
		log.Println("failed to get the list of accounts. error is:", err.Error())
		return
	}

	bytes, err := json.Marshal(accounts)
	if err != nil {
		msg.Respond([]byte("error while marshalling"))
		log.Println("error while marshalling. error is: ", err.Error())
		return
	}

	if err := msg.Respond(bytes); err != nil {
		log.Println("could not send the response. Error is:", err.Error())
		return
	}
}

func (h *Handler) Transaction(msg *nats.Msg) {
	var transaction model.Transaction
	if err := json.Unmarshal(msg.Data, &transaction); err != nil {
		msg.Respond([]byte("error while unmarshalling"))
		log.Println("error while unmarshalling. error is: ", err.Error())
		return
	}

	id, err := h.Service.Account.Transaction(transaction)
	if err != nil {
		msg.Respond([]byte("failed to complete transaction"))
		log.Println("failed to complete transaction. Error is:", err.Error())
		return
	}

	idStr := strconv.Itoa(id)

	if err := msg.Respond([]byte(idStr)); err != nil {
		log.Println("could not send the response. Error is:", err.Error())
		return
	}
}

func (h *Handler) BuyBook(msg *nats.Msg) {
	var transfer Transfer
	if err := json.Unmarshal(msg.Data, &transfer); err != nil {
		msg.Respond([]byte("error while unmarshalling"))
		log.Println("error while unmarshalling. error is: ", err.Error())
		return
	}

	buyerAccounts, err := h.Service.Account.GetAccounts(transfer.BuyerId)
	sellerAccounts, err := h.Service.Account.GetAccounts(transfer.SellerId)

	if err != nil {
		msg.Respond([]byte("error while getting the list of accounts"))
		log.Println("error while getting the list of accounts. error is: ", err.Error())
		return
	}

	var transaction = model.Transaction{
		SenderAccount:   buyerAccounts[0].Id,
		ReceiverAccount: sellerAccounts[0].Id,
		Amount:          transfer.Amount,
	}

	tranId, err := h.Service.Account.Transaction(transaction)
	if err != nil {
		msg.Respond([]byte("error while transaction"))
		log.Println("error while transaction. error is: ", err.Error())
		return
	}

	idStr := strconv.Itoa(tranId)

	if err := msg.Respond([]byte(idStr)); err != nil {
		log.Println("failed to response. Error is:", err.Error())
	}
}
