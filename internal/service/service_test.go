package service

import (
	"context"
	"errors"
	"log/slog"
	dbmock "merch_shop/internal/db/mocks"
	tokenizermock "merch_shop/pkg/tokenizer/mocks"
	"merch_shop/pkg/xerrors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	datadase := dbmock.NewDB(t)
	tokenizer := tokenizermock.NewTokenizer(t)
	service := New(datadase, slog.Default(), tokenizer)

	type res struct {
		token string
		err   xerrors.Xerror
	}

	t.Run("default invalid parms validation", func(t *testing.T) {
		testCases := []struct {
			name     string
			username string
			password string
			expected res
		}{
			{name: "password lenth < min", username: strings.Repeat("1", maxUsernameLenth-1),
				password: strings.Repeat("1", minPasswordLenth-1), expected: res{"", xerrors.New(errPasswordInvalid, http.StatusBadRequest)}},
			{name: "password lenth > max", username: strings.Repeat("1", maxUsernameLenth-1),
				password: strings.Repeat("1", maxPasswordLenth+1), expected: res{"", xerrors.New(errPasswordInvalid, http.StatusBadRequest)}},
			{name: "username lenth < min", username: strings.Repeat("1", minUsernameLenth-1),
				password: strings.Repeat("1", maxPasswordLenth-1), expected: res{"", xerrors.New(errUsernameInvalid, http.StatusBadRequest)}},
			{name: "username lenth > max", username: strings.Repeat("1", maxUsernameLenth+1),
				password: strings.Repeat("1", maxPasswordLenth-1), expected: res{"", xerrors.New(errUsernameInvalid, http.StatusBadRequest)}},
		}

		for _, tc := range testCases {
			token, err := service.AuthentificateUser(context.Background(), tc.username, tc.password)
			assert.Equal(t, token, tc.expected.token)
			assert.Equal(t, err, tc.expected.err)
		}
	})

	t.Run("get user db error", func(t *testing.T) {
		username := strings.Repeat("1", maxUsernameLenth-1)
		password := strings.Repeat("1", maxPasswordLenth-1)
		datadase.On("GetUser", context.Background(), username).Return(nil, "", errors.New("some error"))

		token, err := service.AuthentificateUser(context.Background(), username, password)
		require.Equal(t, token, "")
		require.Equal(t, err, xerrors.New(errSmthWentWrong, http.StatusInternalServerError))
	})
}
