package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Dorrrke/gt4-bookly/internal/config"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	serve    *http.Server
	valid    *validator.Validate
	uService service.UserService
	bService service.BookService
}

func New(cfg config.Config, us service.UserService, bs service.BookService) *Server {
	addrStr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	server := http.Server{
		Addr: addrStr,
	}
	vald := validator.New()
	return &Server{
		serve:    &server,
		valid:    vald,
		uService: us,
		bService: bs,
	}
}

func (s *Server) Run() error {
	log := logger.Get()
	router := s.configRouting()
	s.serve.Handler = router
	log.Info().Str("addr", s.serve.Addr).Msg("server start")
	if err := s.serve.ListenAndServe(); err != nil {
		log.Error().Err(err).Msg("runing server failed")
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.serve.Shutdown(ctx)
}

func (s *Server) configRouting() *gin.Engine {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) { ctx.String(200, "Hello, my friend!") })
	users := router.Group("/users")
	{
		users.GET("/info")
		users.POST("/register", s.registerHendler)
		users.POST("/login", s.loginHendler)
	}
	books := router.Group("/books")
	{
		books.GET("/:id")
		books.DELETE("/:id")
		books.GET("/", s.getBooksHandler)
		books.POST("/", s.addBookHandler)
	}

	return router
}
