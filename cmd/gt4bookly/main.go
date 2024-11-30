package main

import (
	"github.com/Dorrrke/gt4-bookly/internal/config"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/server"
	"github.com/Dorrrke/gt4-bookly/internal/service"
	"github.com/Dorrrke/gt4-bookly/internal/storage"
)

func main() {
	cfg := config.ReadConfig()
	log := logger.Get(cfg.Debug)

	stor := storage.New()
	userService := service.NewUserService(stor)
	bookService := service.NewBookService(stor)
	serve := server.New(cfg, userService, bookService)
	if err := serve.Run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
