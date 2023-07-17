package service

import (
	"biling-nats/book-store/internal/repository"
	"biling-nats/book-store/model"
	"log"
)

type BookService struct {
	repo repository.Book
}

func NewBookService(repo repository.Book) *BookService {
	return &BookService{repo: repo}
}

func (b *BookService) AddBook(book model.Book) (int, error) {
	bookId, err := b.repo.AddBook(book)
	if err != nil {
		log.Println("failed to create a new book. Error is:", err.Error())
		return -1, err
	}

	return bookId, nil
}

func (b *BookService) GetBook(bookId int) (model.Book, error) {
	book, err := b.repo.GetBook(bookId)
	if err != nil {
		log.Println("failed to get a book by id. Error is:", err.Error())
		return model.Book{}, err
	}

	return book, nil
}
