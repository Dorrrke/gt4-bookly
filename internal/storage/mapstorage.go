package storage

import (
	"fmt"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/storage/storageerror"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MapStorage struct {
	stor  map[string]models.User
	bStor map[string]models.Book
}

func New() *MapStorage {
	return &MapStorage{stor: make(map[string]models.User)}
}

func (ms *MapStorage) SaveUser(user models.User) (string, error) {
	log := logger.Get()
	for _, usr := range ms.stor {
		if user.Email == usr.Email {
			return ``, fmt.Errorf("user alredy exist")
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Passoword), bcrypt.DefaultCost)
	if err != nil {
		return ``, err
	}
	user.Passoword = string(hash)
	UID := uuid.New()
	user.UID = UID
	ms.stor[user.UID.String()] = user
	log.Debug().Any("storage", ms.stor).Msg("check storage")
	return UID.String(), nil
}

func (ms *MapStorage) ValidateUser(user models.UserLogin) (string, error) {
	for key, usr := range ms.stor {
		if user.Email == usr.Email {
			if err := bcrypt.CompareHashAndPassword([]byte(usr.Passoword), []byte(user.Passoword)); err != nil {
				return ``, fmt.Errorf("invalid user password")
			}
			return key, nil
		}
	}
	return ``, fmt.Errorf("user no exist")
}

func (ms *MapStorage) SaveBook(book models.Book) (string, error) {
	log := logger.Get()
	for _, b := range ms.bStor {
		if book.Lable == b.Lable && book.Author == b.Author {
			return ``, storageerror.ErrBookAlredyExist
		}
	}
	bID := uuid.New()
	book.BID = bID
	ms.bStor[book.BID.String()] = book
	log.Debug().Any("book storage", ms.stor).Msg("check storage")
	return bID.String(), nil
}

func (ms *MapStorage) GetBooks() ([]models.Book, error) {
	if len(ms.bStor) == 0 {
		return nil, storageerror.ErrEmptyStorage
	}
	var books []models.Book
	for _, book := range ms.bStor {
		books = append(books, book)
	}
	return books, nil
}

func (ms *MapStorage) GetBook(bid string) (models.Book, error) {
	book, ok := ms.bStor[bid]
	if !ok {
		return models.Book{}, storageerror.ErrBookNoFound
	}
	return book, nil
}

func (ms *MapStorage) DeleteBook(bid string) error {
	_, ok := ms.bStor[bid]
	if !ok {
		return storageerror.ErrBookNoFound
	}
	delete(ms.bStor, bid)
	return nil
}