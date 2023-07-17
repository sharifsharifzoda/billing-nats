package handler

import (
	"biling-nats/auth/internal/service"
	"biling-nats/auth/model"
	"encoding/json"
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

type user struct {
	Id   int
	Role string
}

func NewHandler(nats *nats.Conn, service *service.Service) *Handler {
	return &Handler{nats, service}
}

func (h *Handler) Init() {
	sub, err := h.Nats.Subscribe("user.create", h.CreateUser)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())
	}
	defer sub.Unsubscribe()

	loginSub, err := h.Nats.Subscribe("user.login", h.Login)
	if err != nil {
		log.Fatalf("can't subscribe to the subject. Error is: %s", err.Error())
	}
	defer loginSub.Unsubscribe()

	MiddleSub, err := h.Nats.Subscribe("user.middleware", h.AuthMiddleware)
	if err != nil {
		log.Fatalf("can not subscribe to the subject. Error is: %s", err.Error())
	}
	defer MiddleSub.Unsubscribe()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	msg := <-ch
	fmt.Println("Got signal: ", msg)
	log.Println("Shutting down!")
}

func (h *Handler) CreateUser(msg *nats.Msg) {
	var user model.User
	if err := json.Unmarshal(msg.Data, &user); err != nil {
		msg.Respond([]byte("Invalid data provided"))
		log.Println("invalid data provided")
		return
	}
	log.Println(user)

	if err := h.Service.Auth.ValidateUser(user); err != nil {
		if err := msg.Respond([]byte("validate error")); err != nil {
			log.Println("failed to respond. Error is: ", err.Error())
		}
		log.Println("validate error")
		return
	}

	isUsed := h.Service.Auth.IsEmailUsed(user.Email)
	if isUsed {
		msg.Respond([]byte("email is already used"))
		log.Println("email is already used. please log in.")
		return
	}

	id, err := h.Service.Auth.CreateUser(&user)
	if err != nil {
		msg.Respond([]byte("failed to create user"))
		log.Println("failed to create user. Error is: ", err.Error())
		return
	}

	mapId := make(map[string]int)
	mapId["id"] = id

	bytes, err := json.Marshal(mapId)
	if err != nil {
		msg.Respond([]byte("failed to marshal the data"))
		log.Println("failed to marshal the data. Error is: ", err.Error())
		return
	}

	msg.Respond(bytes)
}

func (h *Handler) Login(msg *nats.Msg) {
	var user model.User
	if err := json.Unmarshal(msg.Data, &user); err != nil {
		msg.Respond([]byte("Invalid data provided"))
		log.Println("invalid data provided")
		return
	}

	if err := h.Service.Auth.ValidateUser(user); err != nil {
		if err := msg.Respond([]byte("validate error")); err != nil {
			log.Println("failed to respond. Error is: ", err.Error())
		}
		log.Println("validate error")
		return
	}

	checkedUser, err := h.Service.Auth.CheckUser(user)
	if err != nil {
		msg.Respond([]byte("invalid email or password"))
		log.Println("Invalid email or password. Error is: ", err.Error())
		return
	}

	token, err := h.Service.Auth.GenerateToken(checkedUser)
	if err != nil {
		msg.Respond([]byte(err.Error()))
		log.Println("could not generate token. Error is: ", err.Error())
		return
	}

	if err := msg.Respond([]byte(token)); err != nil {
		log.Println("could not send the response. Error is: ", err.Error())
		return
	}
}

func (h *Handler) AuthMiddleware(msg *nats.Msg) {
	userId, userRole, err := h.Service.Auth.ParseToken(string(msg.Data))
	if err != nil {
		s := fmt.Sprintf("failed to parse the token. Error is: %v", err.Error())
		msg.Respond([]byte(s))
		log.Println("failed to parse token. Error is: ", err.Error())
		return
	}

	var u = user{
		Id:   userId,
		Role: userRole,
	}

	bytes, err := json.Marshal(u)
	if err != nil {
		msg.Respond([]byte("failed to marshal the data"))
		log.Println("failed to marshal the data. Error is: ", err.Error())
		return
	}

	if err := msg.Respond(bytes); err != nil {
		log.Println("failed to send the request. Error is: ", err.Error())
		return
	}
}
