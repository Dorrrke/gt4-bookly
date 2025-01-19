package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/storage/storageerror"

	"github.com/gin-gonic/gin"
)

func (s *BooklyAPI) addBookHandler(ctx *gin.Context) {
	log := logger.Get()
	_, exist := ctx.Get("uid")
	if !exist {
		log.Error().Msg("user ID not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}
	// UID := fmt.Sprintf("%v", uid)
	var bookReq models.BookRequest
	err := ctx.ShouldBindBodyWithJSON(&bookReq)
	if err != nil {
		log.Error().Err(err).Msg("unmarshall body failed")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	writed_at, err := time.Parse("2006-01", bookReq.WritedAt)
	if err != nil {
		log.Error().Err(err).Msg("failed parsing writed time")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	book := models.Book{
		BID:         bookReq.BID,
		Lable:       bookReq.Lable,
		Author:      bookReq.Author,
		Description: bookReq.Description,
		WritedAt:    writed_at,
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

func (s *BooklyAPI) getBooksHandler(ctx *gin.Context) {
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

func (s *BooklyAPI) getBookHandler(ctx *gin.Context) {
	log := logger.Get()
	bid := ctx.Param("id")
	log.Debug().Str("bid", bid).Msg("chek bid from param")
	book, err := s.bService.GetBook(bid)
	if err != nil {
		log.Error().Err(err).Msg("get all books form storage failed")
		if errors.Is(err, storageerror.ErrBookNoFound) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reqBook := models.BookRequest{
		BID:         book.BID,
		Lable:       book.Lable,
		Author:      book.Author,
		Description: book.Description,
		WritedAt:    book.WritedAt.Format("2006-01"),
	}
	ctx.JSON(http.StatusOK, reqBook)
}

func (s *BooklyAPI) deleteBookHandler(ctx *gin.Context) {
	log := logger.Get()
	bid := ctx.Param("id")
	err := s.bService.SetDeleteStatus(bid)
	if err != nil {
		log.Error().Err(err).Msg("delete book failed")
		if errors.Is(err, storageerror.ErrBookNoFound) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.delChan <- struct{}{}
	ctx.String(http.StatusOK, "Book %s was deleted", bid)
}

func (s *BooklyAPI) deleter(ctx context.Context) {
	log := logger.Get()
	defer log.Debug().Msg("deleter stoped")
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
			log.Debug().Int("len", len(s.delChan)).Int("cap", cap(s.delChan)).Msg("deleter chan check")
			if len(s.delChan) == cap(s.delChan) {
				log.Debug().Int("len", len(s.delChan)).Int("cap", cap(s.delChan)).Msg("deleter chan is full")
				for i := 0; i < cap(s.delChan); i++ {
					<-s.delChan
				}
				if err := s.bService.DeleteBooks(); err != nil {
					log.Error().Err(err).Msg("delete all books failed")
					s.ErrChan <- err
					return
				}
			}
		}
	}
}
