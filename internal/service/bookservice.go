package service

import "github.com/Dorrrke/gt4-bookly/internal/domain/models"

type BookStorage interface {
	SaveBook(models.Book) (string, error)
	GetBooks() ([]models.Book, error)
	GetBook(string) (models.Book, error)
	DeleteBook(string) error
}

type BookService struct {
	stor BookStorage
}

func NewBookService(stor BookStorage) BookService {
	return BookService{stor: stor}
}

func (bs *BookService) AddBook(book models.Book) (string, error) {
	return bs.stor.SaveBook(book)
}

func (bs *BookService) GetBooks() ([]models.Book, error) {
	return bs.stor.GetBooks()
}

func (bs *BookService) GetBook(bid string) (models.Book, error) {
	return bs.stor.GetBook(bid)
}

func (bs *BookService) DeleteBook(bid string) error {
	return bs.stor.DeleteBook(bid)
}
