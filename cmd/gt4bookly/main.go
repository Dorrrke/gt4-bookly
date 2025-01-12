package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dorrrke/gt4-bookly/internal/config"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/server"
	"github.com/Dorrrke/gt4-bookly/internal/service"
	"github.com/Dorrrke/gt4-bookly/internal/storage"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.ReadConfig()
	log := logger.Get(cfg.Debug)
	log.Debug().Any("cfg", cfg).Msg("config")

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c

		log.Info().Msg("gracefully stopping...")
		cancel()
	}()

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

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := serve.Run(gCtx); err != nil {
			log.Error().Err(err).Send()
			return err
		}
		return nil
	})
	group.Go(func() error {
		log.Debug().Msg("start listening error channel")
		defer log.Debug().Msg("stop listening error channel")
		return <-serve.ErrChan
	})
	group.Go(func() error {
		<-gCtx.Done()
		return serve.Shutdown(gCtx)
	})
	group.Go(func() error {
		<-gCtx.Done()
		return stor.Close()
	})

	if err := group.Wait(); err != nil {
		log.Error().Err(err).Send()
	}
	log.Info().Msg("server stoped")
}
