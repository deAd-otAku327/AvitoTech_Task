package tokenizer

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	errGenerateToken = errors.New("token generation failed")
	errParseToken    = errors.New("token parsing failed")
	errVerifyToken   = errors.New("token verefication failed")
)

type Tokenizer struct {
	secretKey []byte
}

func New(key string) *Tokenizer {
	return &Tokenizer{
		secretKey: []byte(key),
	}
}

func (t *Tokenizer) GenerateToken(userID int) (*string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iss": "avito-shop",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString(t.secretKey)
	if err != nil {
		return nil, errGenerateToken
	}

	return &token, nil
}

func (t *Tokenizer) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return t.secretKey, nil
	})
	if err != nil {
		return nil, errParseToken
	}

	if !token.Valid {
		return nil, errVerifyToken
	}

	return token, nil
}
