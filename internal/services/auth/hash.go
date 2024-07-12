package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func ValidateHash(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
