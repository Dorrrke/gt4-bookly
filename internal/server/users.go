package server

import (
	"net/http"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"

	"github.com/gin-gonic/gin"
)

func (s *Server) loginHendler(ctx *gin.Context) {
	log := logger.Get()
	var user models.UserLogin
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		log.Error().Err(err).Msg("unmarshall body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("validate user input data failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	UID, err := s.uService.LoginUser(user)
	if err != nil {
		log.Error().Err(err).Msg("user validate failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid input data", "error": err.Error()})
		return
	}
	ctx.String(http.StatusCreated, "User was logined; user id: %s", UID)
}

func (s *Server) registerHendler(ctx *gin.Context) {
	log := logger.Get()
	var user models.User
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		log.Error().Err(err).Msg("unmarshall body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("validate user input data failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	UID, err := s.uService.RegisterUser(user)
	if err != nil {
		log.Error().Err(err).Msg("user register failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid input data", "error": err.Error()})
		return
	}
	ctx.String(http.StatusCreated, "User was created; user id: %s", UID)
}
