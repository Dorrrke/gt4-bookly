package server

import (
	"net/http"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/server/utils"

	"github.com/gin-gonic/gin"
)

func (s *BooklyAPI) loginHendler(ctx *gin.Context) { //nolint:dupl //todo
	log := logger.Get()
	var user models.UserLogin
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		log.Error().Err(err).Msg("unmarshall login body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("validate login user input data failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := s.uService.LoginUser(user)
	if err != nil {
		log.Error().Err(err).Msg("user login validate failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid input data", "error": err.Error()})
		return
	}
	token, err := utils.CreateJWT(uid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Header("Authorization", token)
	ctx.String(http.StatusCreated, "User was logined; user id: %s", uid)
}

func (s *BooklyAPI) registerHendler(ctx *gin.Context) { //nolint:dupl //todo
	log := logger.Get()
	var user models.User
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		log.Error().Err(err).Msg("unmarshall body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = s.valid.Struct(user); err != nil {
		log.Error().Err(err).Msg("validate user input data failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, err := s.uService.RegisterUser(user)
	if err != nil {
		log.Error().Err(err).Msg("user register failed")
		ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid input data", "error": err.Error()})
		return
	}
	token, err := utils.CreateJWT(uid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Header("Authorization", token)
	ctx.String(http.StatusCreated, "User was created; user id: %s", uid)
}
