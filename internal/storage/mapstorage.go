package storage

import (
	"errors"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/storage/storageerror"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MapUserStorage struct {
	stor map[string]models.User
}

func NewUserStor() *MapUserStorage {
	return &MapUserStorage{
		stor: make(map[string]models.User),
	}
}

func (ms *MapUserStorage) SaveUser(user models.User) (string, error) {
	log := logger.Get()
	for _, usr := range ms.stor {
		if user.Email == usr.Email {
			return ``, errors.New("user alredy exist")
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Passoword), bcrypt.DefaultCost)
	if err != nil {
		return ``, err
	}
	user.Passoword = string(hash)
	uid := uuid.New()
	user.UID = uid
	ms.stor[user.UID.String()] = user
	log.Debug().Any("storage", ms.stor).Msg("check storage")
	return uid.String(), nil
}

func (ms *MapUserStorage) ValidateUser(user models.UserLogin) (string, error) {
	for key, usr := range ms.stor {
		if user.Email == usr.Email {
			if err := bcrypt.CompareHashAndPassword([]byte(usr.Passoword), []byte(user.Passoword)); err != nil {
				return ``, errors.New("invalid user password")
			}
			return key, nil
		}
	}
	return ``, errors.New("user no exist")
}

type MapBookStorage struct {
	bStor map[string]models.Book
}

func NewBookStor() *MapBookStorage {
	return &MapBookStorage{
		bStor: make(map[string]models.Book),
	}
}

func (ms *MapBookStorage) SaveBook(book models.Book) (string, error) {
	log := logger.Get()
	for _, b := range ms.bStor {
		if book.Lable == b.Lable && book.Author == b.Author {
			return ``, storageerror.ErrBookAlredyExist
		}
	}
	bID := uuid.New()
	book.BID = bID
	ms.bStor[book.BID.String()] = book
	log.Debug().Any("book storage", ms.bStor).Msg("check storage")
	return bID.String(), nil
}

func (ms *MapBookStorage) GetBooks() ([]models.Book, error) {
	if len(ms.bStor) == 0 {
		return nil, storageerror.ErrEmptyStorage
	}
	var books []models.Book
	for _, book := range ms.bStor {
		books = append(books, book)
	}
	return books, nil
}

func (ms *MapBookStorage) GetBook(bid string) (models.Book, error) {
	book, ok := ms.bStor[bid]
	if !ok {
		return models.Book{}, storageerror.ErrBookNoFound
	}
	return book, nil
}

func (ms *MapBookStorage) DeleteBook(bid string) error {
	_, ok := ms.bStor[bid]
	if !ok {
		return storageerror.ErrBookNoFound
	}
	delete(ms.bStor, bid)
	return nil
}

func (ms *MapBookStorage) DeleteBooks() error {
	return nil
}

func (ms *MapBookStorage) SetDeleteBookStatus(bid string) error {
	return nil
}
