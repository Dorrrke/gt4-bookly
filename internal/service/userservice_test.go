package service

import (
	"testing"

	"github.com/Dorrrke/gt4-bookly/internal/domain/models"
	"github.com/Dorrrke/gt4-bookly/internal/logger"
	"github.com/Dorrrke/gt4-bookly/internal/service/mocks"
	"github.com/Dorrrke/gt4-bookly/internal/storage/storageerror"
	"github.com/stretchr/testify/assert"
)

func TestLoginUser(t *testing.T) {
	logger.Get(true)
	type want struct {
		uid string
		err error
	}
	type test struct {
		name string
		user models.UserLogin
		want want
	}

	tests := []test{
		{
			name: "successfull user login",
			user: models.UserLogin{
				Email:     "testy@ya.ru",
				Passoword: "qwerty123",
			},
			want: want{
				uid: "test-uid",
				err: nil,
			},
		},
		{
			name: "successfull error abort",
			user: models.UserLogin{
				Email:     "shadowuser@ya.ru",
				Passoword: "qwerty123",
			},
			want: want{
				uid: "",
				err: storageerror.ErrUserNoExist,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storMock := mocks.NewStorage(t)
			storMock.On("ValidateUser", tc.user).Return(tc.want.uid, tc.want.err)
			testService := UserService{
				stor: storMock,
			}
			uid, err := testService.LoginUser(tc.user)
			assert.ErrorIs(t, err, tc.want.err)
			assert.Equal(t, tc.want.uid, uid)
		})
	}
}
