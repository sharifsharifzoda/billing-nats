package main

import (
	"biling-nats/api/handler"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	l := log.New(os.Stdout, "task-api", log.LstdFlags)

	conn, err := nats.Connect(nats.DefaultURL, nats.Name("sharif"))
	if err != nil {
		log.Fatal("error while connecting to the NATS server. Error is: ", err.Error())
	}
	fmt.Println("connected to the nats server")

	defer conn.Close()

	h := handler.NewHandlerApi(l, conn)

	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.Login)
	}

	account := r.Group("/account", h.AuthMiddleware)
	{
		account.POST("/", h.CreateAccount)
		account.GET("/", h.GetAccounts)
	}

	tr := r.Group("/transaction")
	{
		tr.POST("/", h.Transaction)
	}

	book := r.Group("/book", h.AuthMiddleware)
	{
		book.POST("/", h.AddBook)
		book.POST("/buy/:id", h.BuyBook)
	}

	s := http.Server{
		Addr:         ":9090",
		Handler:      r,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s.Shutdown(ctx)
}
