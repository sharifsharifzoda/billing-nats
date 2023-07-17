package service

import (
	"biling-nats/book-store/internal/repository"
	"biling-nats/book-store/model"
)

type Book interface {
	AddBook(book model.Book) (int, error)
	GetBook(bookId int) (model.Book, error)
}

type Service struct {
	Book
}

func NewService(repo *repository.Repository) *Service {
	return &Service{NewBookService(repo.Book)}
}
