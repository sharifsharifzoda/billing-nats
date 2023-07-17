package repository

import (
	"biling-nats/book-store/model"
	"database/sql"
)

type Book interface {
	AddBook(book model.Book) (int, error)
	GetBook(bookId int) (model.Book, error)
}

type Repository struct {
	Book
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{NewBookRepo(db)}
}
