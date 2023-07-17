package repository

import (
	"biling-nats/book-store/model"
	"database/sql"
)

const (
	InsertBook  = `INSERT INTO books(name, description, price, seller_id) VALUES (($1), ($2), ($3), ($4)) RETURNING id;`
	GetBookById = `SELECT id, name, description, price, seller_id FROM books WHERE id = $1;`
)

type BookRepo struct {
	Db *sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepo {
	return &BookRepo{Db: db}
}

func (b *BookRepo) AddBook(book model.Book) (int, error) {
	var id int
	row := b.Db.QueryRow(InsertBook, book.Name, book.Description, book.Price, book.SellerId)
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (b *BookRepo) GetBook(bookId int) (model.Book, error) {
	var book model.Book
	//fmt.Println(bookId)
	row := b.Db.QueryRow(GetBookById, bookId)
	//fmt.Println(row)
	err := row.Scan(&book.Id, &book.Name, &book.Description, &book.Price, &book.SellerId)
	if err != nil {
		return model.Book{}, err
	}

	return book, nil
}
