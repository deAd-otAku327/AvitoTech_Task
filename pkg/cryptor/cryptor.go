package cryptor

import "golang.org/x/crypto/bcrypt"

// WARNING: for future more parametrisation first args of funcs MUST be context.Context.
// Current configuration is not very heavy for productivity.

func EncryptKeyword(pass string) (string, error) {
	// WARNING: changing cost value may be very heavy for productivity.
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
