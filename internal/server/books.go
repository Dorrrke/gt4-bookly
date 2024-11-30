package server

import (
	"errors"
	"net/http"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/storage/storageerror"

	"github.com/gin-gonic/gin"
)

func (s *Server) addBookHandler(ctx *gin.Context) {
	log := logger.Get()
	var book models.Book
	err := ctx.ShouldBindBodyWithJSON(&book)
	if err != nil {
		log.Error().Err(err).Msg("unmarshall body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bid, err := s.bService.AddBook(book)
	if err != nil {
		log.Error().Err(err).Msg("save book failed")
		if errors.Is(err, storageerror.ErrBookAlredyExist) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusCreated, "Book %s was saved", bid)
}

func (s *Server) getBooksHandler(ctx *gin.Context) {
	log := logger.Get()
	books, err := s.bService.GetBooks()
	if err != nil {
		log.Error().Err(err).Msg("get all books form storage failed")
		if errors.Is(err, storageerror.ErrEmptyStorage) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}
