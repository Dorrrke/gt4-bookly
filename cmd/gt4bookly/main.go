package main

import (
	"context"

	"github.com/Dorrrke/gt4-bookly/internal/config"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/server"
	"github.com/Dorrrke/gt4-bookly/internal/service"
	"github.com/Dorrrke/gt4-bookly/internal/storage"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.ReadConfig()
	log := logger.Get(cfg.Debug)

	var userService service.UserService
	var bookService service.BookService

	err := storage.Migrations(cfg.DbDSN, cfg.MigratePath)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	stor, err := storage.NewDB(context.Background(), cfg.DbDSN)
	if err != nil {
		log.Error().Err(err).Send()
		bStor := storage.NewBookStor()
		uStor := storage.NewUserStor()
		userService = service.NewUserService(uStor)
		bookService = service.NewBookService(bStor)
	} else {
		userService = service.NewUserService(stor)
		bookService = service.NewBookService(stor)
	}

	serve := server.New(cfg, userService, bookService)
	if err := serve.Run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
