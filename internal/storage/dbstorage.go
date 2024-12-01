package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/storage/storageerror"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type DBStorage struct {
	conn *pgx.Conn
}

func NewDB(ctx context.Context, addr string) (*DBStorage, error) {
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &DBStorage{conn: conn}, nil
}

func (dbs *DBStorage) SaveUser(user models.User) (string, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Passoword), bcrypt.DefaultCost)
	if err != nil {
		return ``, err
	}
	user.Passoword = string(hash)
	uid := uuid.New()
	user.UID = uid
	_, err = dbs.conn.Exec(ctx, "INSERT INTO users (uid, name, email, pass, age) VALUES ($1, $2, $3, $4, $5)",
		user.UID, user.Name, user.Email, user.Passoword, user.Age)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return "", storageerror.ErrUserAlredyExist
			}
		}
		log.Error().Err(err).Msg("failed isert user")
		return "", err
	}
	return uid.String(), nil
}

func (dbs *DBStorage) ValidateUser(user models.UserLogin) (string, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := dbs.conn.QueryRow(ctx, "SELECT uid, email, pass FROM users WHERE email = $1", user.Email)
	var usr models.User
	if err := row.Scan(&usr.UID, &usr.Email, &usr.Passoword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storageerror.ErrUserNoExist
		}
		log.Error().Err(err).Msg("failed scan db data")
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usr.Passoword), []byte(user.Passoword)); err != nil {
		return "", storageerror.ErrInvalidPassword
	}
	return usr.UID.String(), nil
}

func (dbs *DBStorage) GetBooks() ([]models.Book, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := dbs.conn.Query(ctx, "SELECT * FROM books")
	if err != nil {
		log.Error().Err(err).Msg("failed get data from table books")
		return nil, err
	}
	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err = rows.Scan(&book.BID, &book.Lable, &book.Author, &book.Description, &book.WritedAt); err != nil {
			log.Error().Err(err).Msg("failed scan rows data")
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func Migrations(dbDsn string, migratePath string) error {
	log := logger.Get()
	migrPath := fmt.Sprintf("file://%s", migratePath)
	m, err := migrate.New(migrPath, dbDsn)
	if err != nil {
		log.Error().Err(err).Msg("failed to db conntect")
		return err
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Debug().Msg("no migratons apply")
			return nil
		}
		log.Error().Err(err).Msg("run migrations failed")
		return err
	}
	log.Debug().Msg("all migrations apply")
	return nil
}
