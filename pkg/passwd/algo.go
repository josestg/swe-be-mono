package passwd

import (
	"golang.org/x/crypto/bcrypt"
)

// bcryptImpl is a type that implements the HashComparer interface using the bcrypt algorithm.
type bcryptImpl int

// Hash generates a bcrypt hash from the specified plaintext password using the configured cost.
// It returns the resulting hash as a string and any errors that occur during the hash generation.
func (b bcryptImpl) Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), int(b))
	return string(hash), err
}

// Compare compares the specified plaintext password with the specified bcrypt hash.
// It returns an error if the comparison fails.
func (b bcryptImpl) Compare(hash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

// BcryptDefaultCost is a bcrypt algorithm with default cost.
const BcryptDefaultCost = bcryptImpl(bcrypt.DefaultCost)
