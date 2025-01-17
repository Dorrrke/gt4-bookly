package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Dorrrke/gt4-bookly/internal/config"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/server/utils"
	"github.com/Dorrrke/gt4-bookly/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BooklyAPI struct {
	serve    *http.Server
	valid    *validator.Validate
	uService service.UserService
	bService service.BookService
	delChan  chan struct{}
	ErrChan  chan error
}

func New(cfg config.Config, us service.UserService, bs service.BookService) *BooklyAPI {
	addrStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := http.Server{ //nolint:gosec //todo
		Addr: addrStr,
	}
	vald := validator.New()
	srv := BooklyAPI{
		serve:    &server,
		valid:    vald,
		uService: us,
		bService: bs,
		delChan:  make(chan struct{}, 10),
		ErrChan:  make(chan error, 10),
	}
	return &srv
}

func (s *BooklyAPI) Run(ctx context.Context) error {
	log := logger.Get()
	router := s.configRouting()
	s.serve.Handler = router
	go s.deleter(ctx)
	log.Info().Str("addr", s.serve.Addr).Msg("server start")
	if err := s.serve.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("runing server failed")
		return err
	}
	return nil
}

func (s *BooklyAPI) Shutdown(ctx context.Context) error {
	return s.serve.Shutdown(ctx)
}

func (s *BooklyAPI) JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.Get()
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.String(http.StatusUnauthorized, "token is empty")
			return
		}
		UID, err := utils.ValidToken(token)
		if err != nil {
			log.Error().Err(err).Send()
			if errors.Is(err, utils.ErrInvalidToken) {
				ctx.String(http.StatusUnauthorized, err.Error())
				return
			}
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Set("uid", UID)
		ctx.Next()
	}
}

func (s *BooklyAPI) configRouting() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "Hello, my friend!") })
	users := router.Group("/users")
	{
		users.GET("/info")
		users.POST("/register", s.registerHendler)
		users.POST("/login", s.loginHendler)
	}
	books := router.Group("/books")
	{
		books.GET("/:id", s.getBookHandler)
		books.GET("/", s.getBooksHandler)
		books.POST("/", s.JWTAuthMiddleware(), s.addBookHandler)
		books.DELETE("/:id", s.JWTAuthMiddleware(), s.deleteBookHandler)
	}
	return router
}
