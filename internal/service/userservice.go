package service

import (
	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
)

type Storage interface {
	SaveUser(models.User) (string, error)
	ValidateUser(models.UserLogin) (string, error)
}

type UserService struct {
	stor Storage
}

func NewUserService(stor Storage) UserService {
	return UserService{stor: stor}
}

func (us *UserService) LoginUser(user models.UserLogin) (string, error) {
	log := logger.Get()
	uid, err := us.stor.ValidateUser(user)
	if err != nil {
		log.Error().Err(err).Msg("validate user failed")
		return ``, err
	}
	return uid, nil
}

func (us *UserService) RegisterUser(user models.User) (string, error) {
	log := logger.Get()
	uid, err := us.stor.SaveUser(user)
	if err != nil {
		log.Error().Err(err).Msg("save user failed")
		return ``, err
	}
	return uid, nil
}
