package hasher

import "golang.org/x/crypto/bcrypt"

type BcyptHasher struct{}

// Hash the given value.
func (b *BcyptHasher) Make(value string) string {
	password, _ := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	return string(password)
}

// Check the given plain value against a hash.
func (b *BcyptHasher) Check(value string, hashedValue string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	return err == nil
}

func NewBcyptHasher() *BcyptHasher {
	return &BcyptHasher{}
}
