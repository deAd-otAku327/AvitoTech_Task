package cryptor

import "golang.org/x/crypto/bcrypt"

// WARNING: for future more parametrisation first args of funcs MUST be context.Context.
// Current configuration is not very heavy for productivity.

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --testonly --name=Cryptor
type Cryptor interface {
	EncryptKeyword(pass string) (string, error)
	CompareHashAndPassword(hash, password string) error
}

type cryptor struct{}

func New() Cryptor {
	return &cryptor{}
}

func (c cryptor) EncryptKeyword(pass string) (string, error) {
	// WARNING: changing cost value may be very heavy for productivity.
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (c cryptor) CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
